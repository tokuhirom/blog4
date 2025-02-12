package admin

import (
	"bytes"
	"database/sql"
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
	"log"
	"net/http"
	"time"
)

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

// TODO auth
func Router(cfg server.Config, queries *admindb.Queries) *chi.Mux {
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
	r.Get("/api/entries", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Query().Get("")
		_, err := queries.GetLatestEntries(r.Context(), admindb.GetLatestEntriesParams{
			Column1: nil,
			LastEditedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			Limit: 100,
		})
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("{\"hello\": \"world\"}"))
		if err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})
	r.Get("/assets/*", handleAssets)
	return r
}
