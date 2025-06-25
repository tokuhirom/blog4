package admin

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/admin/openapi"
	"github.com/tokuhirom/blog4/server/sobs"
)

func Router(cfg server.Config, db *sql.DB, sobsClient *sobs.SobsClient) (*chi.Mux, error) {
	if cfg.AdminUser == "" {
		slog.Warn("AdminUser is not set")
	}
	if cfg.AdminPassword == "" {
		slog.Warn("AdminPassword is not set")
		return nil, fmt.Errorf("AdminPassword is not set")
	}

	r := chi.NewRouter()
	//r.Use(middleware.BasicAuth("admin", map[string]string{cfg.AdminUser: cfg.AdminPassword}))
	if len(cfg.AllowedOrigins) > 0 {
		slog.Info("CORS is allowed", slog.Any("origins", cfg.AllowedOrigins))
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   cfg.AllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}
	dir, _ := os.Getwd()
	indexHtmlHandler := func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile(filepath.Join("web/admin/dist/index.html"))
		if err != nil {
			slog.Error("failed to read index.html", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return // 404
		}
		http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(file))
	}
	r.Get("/", indexHtmlHandler)
	r.Get("/login", indexHtmlHandler)
	r.HandleFunc("/entry/*", indexHtmlHandler)
	filesDir := http.Dir(filepath.Join(dir, "web/admin/dist"))
	r.Handle("/assets/*", http.StripPrefix("/admin/", http.FileServer(filesDir)))

	queries := admindb.New(db)
	apiService := adminApiService{
		queries:       queries,
		db:            db,
		hubUrls:       cfg.GetHubUrls(),
		paapiClient:   NewPAAPIClient(cfg.AmazonPaapi5AccessKey, cfg.AmazonPaapi5SecretKey),
		S3Client:      sobsClient,
		adminUser:     cfg.AdminUser,
		adminPassword: cfg.AdminPassword,
		isSecure:      !cfg.LocalDev,
	}
	adminApiHandler, err := openapi.NewServer(&apiService,
		openapi.WithPathPrefix("/admin/api"),
		openapi.WithErrorHandler(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) {
			slog.Error("API error",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Any("error", err))
			// Use ogen's default error handler to properly format the response
			ogenerrors.DefaultErrorHandler(ctx, w, r, err)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin API handler: %w", err)
	}
	// Create a subrouter for API with session middleware
	apiRouter := chi.NewRouter()

	// Add HTTP context middleware
	apiRouter.Use(HTTPContextMiddleware)

	// Apply session middleware only to non-auth endpoints
	apiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth check for auth endpoints
			if strings.HasPrefix(r.URL.Path, "/admin/api/auth/") {
				next.ServeHTTP(w, r)
				return
			}
			// Apply session middleware for other endpoints
			SessionMiddleware(queries)(next).ServeHTTP(w, r)
		})
	})

	apiRouter.Mount("/", adminApiHandler)
	r.Mount("/api/", apiRouter)

	return r, nil
}
