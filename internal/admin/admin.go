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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ogen-go/ogen/ogenerrors"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/internal/admin/openapi"
	"github.com/tokuhirom/blog4/server"
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
	if cfg.LocalDev {
		slog.Info("LocalDevelopment mode enabled. CORS is allowed", slog.String("origin", "http://localhost:6173"))
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:6173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	} else {
		slog.Info("LocalDevelopment mode disabled. CORS is not allowed. BasicAuth enabled", slog.String("admin_user", cfg.AdminUser))
		r.Use(middleware.BasicAuth(
			"admin",
			map[string]string{
				cfg.AdminUser: cfg.AdminPassword,
			},
		))
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
	r.HandleFunc("/entry/*", indexHtmlHandler)
	filesDir := http.Dir(filepath.Join(dir, "web/admin/dist"))
	r.Handle("/assets/*", http.StripPrefix("/admin/", http.FileServer(filesDir)))

	queries := admindb.New(db)
	apiService := adminApiService{
		queries:     queries,
		db:          db,
		hubUrls:     cfg.GetHubUrls(),
		paapiClient: NewPAAPIClient(cfg.AmazonPaapi5AccessKey, cfg.AmazonPaapi5SecretKey),
		S3Client:    sobsClient,
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
	r.Mount("/api/", adminApiHandler)

	return r, nil
}
