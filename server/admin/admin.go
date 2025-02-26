package admin

import (
	"bytes"
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/admin/openapi"
	"github.com/tokuhirom/blog4/server/sobs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Router(cfg server.Config, db *sql.DB, sobsClient *sobs.SobsClient) *chi.Mux {
	if cfg.AdminUser == "" {
		println("AdminUser is not set")
	}
	if cfg.AdminPassword == "" {
		println("AdminPassword is not set")
	}

	r := chi.NewRouter()
	//r.Use(middleware.BasicAuth("admin", map[string]string{cfg.AdminUser: cfg.AdminPassword}))
	if cfg.LocalDev {
		log.Print("LocalDevelopment mode enabled. CORS is allowed for http://localhost:5173")
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:5173"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	} else {
		log.Printf("LocalDevelopment mode disabled. CORS is not allowed... And enable BasicAuth for %s",
			cfg.AdminUser)
		r.Use(middleware.BasicAuth(
			"admin",
			map[string]string{
				cfg.AdminUser: cfg.AdminPassword,
			},
		))
	}
	dir, _ := os.Getwd()
	indexHtmlHandler := func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile(filepath.Join("server/admin/frontend/dist/index.html"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return // 404
		}
		http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(file))
	}
	r.Get("/", indexHtmlHandler)
	r.HandleFunc("/entry/*", indexHtmlHandler)
	filesDir := http.Dir(filepath.Join(dir, "server/admin/frontend/dist"))
	r.Handle("/assets/*", http.StripPrefix("/admin/", http.FileServer(filesDir)))

	queries := admindb.New(db)
	apiService := adminApiService{
		queries:     queries,
		db:          db,
		hubUrls:     cfg.GetHubUrls(),
		paapiClient: NewPAAPIClient(cfg.AmazonPaapi5AccessKey, cfg.AmazonPaapi5SecretKey),
		S3Client:    sobsClient,
	}
	adminApiHandler, err := openapi.NewServer(&apiService, openapi.WithPathPrefix("/admin/api"))
	if err != nil {
		return nil
	}
	if cfg.AdminPassword == "" {
		log.Fatalf("Missing AdminPassword")
	}
	r.Mount("/api/", adminApiHandler)

	return r
}
