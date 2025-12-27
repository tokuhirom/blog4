package admin

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
)

func Router(cfg server.Config, db *sql.DB) (*chi.Mux, error) {
	if cfg.AdminUser == "" {
		slog.Warn("AdminUser is not set")
	}
	if cfg.AdminPassword == "" {
		slog.Warn("AdminPassword is not set")
		return nil, fmt.Errorf("AdminPassword is not set")
	}

	r := chi.NewRouter()
	queries := admindb.New(db)

	// htmx routes with gin and session middleware (mounted at root)
	htmxRouter := SetupHtmxRouter(queries, cfg)
	r.Mount("/", htmxRouter)

	return r, nil
}
