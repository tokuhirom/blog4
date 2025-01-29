package admin

import (
	"bytes"
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tokuhirom/blog3/server"
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
func Router(cfg server.Config) *chi.Mux {
	if cfg.AdminUser == "" {
		println("AdminUser is not set")
	}
	if cfg.AdminPassword == "" {
		println("AdminPassword is not set")
	}

	r := chi.NewRouter()
	r.Use(middleware.BasicAuth("admin", map[string]string{cfg.AdminUser: cfg.AdminPassword}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := frontendFS.ReadFile("frontend/dist/index.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return // 404
		}
		http.ServeContent(w, r, "index.html", time.Time{}, bytes.NewReader(file))
	})
	r.Get("/assets/*", handleAssets)
	return r
}
