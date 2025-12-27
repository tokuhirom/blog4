package router

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/internal/admin"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/sobs"
)

func BuildRouter(cfg server.Config, sqlDB *sql.DB, sobsClient *sobs.SobsClient) (*gin.Engine, error) {
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
		slog.Info("CheckWebAccelHeader validation enabled")
		r.Use(func(c *gin.Context) {
			if c.Request.URL.Path != "/healthz" {
				gotToken := c.GetHeader("X-WebAccel-Guard")
				if gotToken != cfg.WebAccelGuard {
					slog.Warn("invalid X-WebAccel-Guard header", slog.String("got_token", gotToken), slog.String("path", c.Request.URL.Path))
					c.String(http.StatusBadRequest, "Invalid X-WebAccel-Guard header")
					c.Abort()
					return
				}
			}
			c.Next()
		})
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
	setupPublicRoutes(r, publicQueries)

	return r, nil
}

func setupPublicRoutes(r *gin.Engine, queries *publicdb.Queries) {
	r.GET("/", func(c *gin.Context) {
		server.RenderTopPage(c, queries)
	})
	r.GET("/feed", func(c *gin.Context) {
		server.RenderFeed(c, queries)
	})
	r.GET("/entry/*filepath", func(c *gin.Context) {
		server.RenderEntryPage(c, queries)
	})
	r.GET("/static/main.css", func(c *gin.Context) {
		server.RenderStaticMainCss(c)
	})
}
