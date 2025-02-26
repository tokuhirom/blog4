package admin

import (
	"bytes"
	"database/sql"
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/admin/openapi"
	"github.com/tokuhirom/blog4/server/sobs"
	"log"
	"net/http"
	"time"
)

//go:embed frontend/dist/index.html
var frontendFS embed.FS

/*
//go:embed frontend/*
var frontendFS embed.FS

func handleAssets(writer http.ResponseWriter, request *http.Request) {
	path := chi.URLParam(request, "*")
	file, err := frontendFS.ReadFile("frontend/dist/assets/" + path)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return // 404
	}
	http.ServeContent(writer, request, path, time.Time{}, bytes.NewReader(file))
}
*/

// TODO auth
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
		log.Printf("LocalDevelopment mode disabled. CORS is not allowed")
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := frontendFS.ReadFile("frontend/dist/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return // 404
		}
		http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(file))
	})
	//r.Get("/assets/*", handleAssets)

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
	r.Mount("/api/", adminApiHandler)

	return r
}
