package router

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/internal/admin"
	middleware2 "github.com/tokuhirom/blog4/internal/middleware"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/sobs"
)

func BuildRouter(cfg server.Config, sqlDB *sql.DB, sobsClient *sobs.SobsClient) (*chi.Mux, error) {
	publicQueries := publicdb.New(sqlDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	if cfg.WebAccelGuard != "" {
		slog.Info("CheckWebAccelHeader validation enabled")
		r.Use(middleware2.CheckWebAccelHeader(cfg.WebAccelGuard))
	}

	adminRouter, err := admin.Router(cfg, sqlDB, sobsClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create admin router: %w", err)
	}
	r.Mount("/admin", adminRouter)
	r.Mount("/", server.Router(publicQueries))
	r.Get("/healthz", HealthzHandler)
	r.Get("/git_hash", GitHashHandler)
	return r, nil
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("ok"))
	if err != nil {
		slog.Error("failed to write healthz response", slog.Any("error", err))
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
		slog.Error("failed to write git hash response", slog.Any("error", err))
	}
}
