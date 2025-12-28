package router

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/blog4/internal"
	"github.com/tokuhirom/blog4/internal/public"
	"github.com/tokuhirom/blog4/internal/sobs"

	"github.com/tokuhirom/blog4/internal/middleware"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/internal/admin"
)

func BuildRouter(cfg internal.Config, sqlDB *sql.DB, sobsClient *sobs.SobsClient) (*gin.Engine, error) {
	// Validate admin config
	if cfg.AdminUser == "" {
		slog.Warn("AdminUser is not set")
	}
	if cfg.AdminPassword == "" {
		slog.Warn("AdminPassword is not set")
		return nil, fmt.Errorf("AdminPassword is not set")
	}

	// Set gin mode
	if !cfg.LocalDev {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create main router
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Add WebAccel guard middleware if configured
	if cfg.WebAccelGuard != "" {
		slog.Info("CheckWebAccelGuard validation enabled")
		r.Use(middleware.CheckWebAccelGuard(cfg))
	}

	// Health check endpoints
	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.GET("/git_hash", func(c *gin.Context) {
		gitHash := os.Getenv("GIT_HASH")
		if gitHash == "" {
			c.String(http.StatusInternalServerError, "GIT_HASH not set")
			return
		}
		c.String(http.StatusOK, gitHash)
	})

	// Setup admin routes
	adminQueries := admindb.New(sqlDB)
	adminGroup := r.Group("/admin")
	admin.SetupAdminRoutes(adminGroup, adminQueries, sobsClient, cfg)

	// Setup public routes
	publicQueries := publicdb.New(sqlDB)
	public.SetupPublicRoutes(r, publicQueries)

	return r, nil
}
