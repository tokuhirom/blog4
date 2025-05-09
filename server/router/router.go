package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/admin"
	middleware2 "github.com/tokuhirom/blog4/server/middleware"
	"github.com/tokuhirom/blog4/server/sobs"
	"log"
	"net/http"
	"os"
)

func BuildRouter(cfg server.Config, sqlDB *sql.DB, sobsClient *sobs.SobsClient) *chi.Mux {
	publicQueries := publicdb.New(sqlDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	if cfg.WebAccelGuard != "" {
		log.Printf("CheckWebAccelHeader validation enabled")
		r.Use(middleware2.CheckWebAccelHeader(cfg.WebAccelGuard))
	}

	r.Mount("/admin", admin.Router(cfg, sqlDB, sobsClient))
	r.Mount("/", server.Router(publicQueries))
	r.Get("/healthz", HealthzHandler)
	r.Get("/git_hash", GitHashHandler)
	return r
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func GitHashHandler(w http.ResponseWriter, r *http.Request) {
	gitHash := os.Getenv("GIT_HASH")
	if gitHash == "" {
		http.Error(w, "GIT_HASH not set", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(gitHash))
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}
