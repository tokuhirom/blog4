package admin

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/internal"
	"github.com/tokuhirom/blog4/internal/ogimage"
	"github.com/tokuhirom/blog4/internal/sobs"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

// NoCacheMiddleware sets Cache-Control headers to prevent caching of dynamic content
func NoCacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Don't add no-cache headers for static files (they should be cached)
		if strings.HasPrefix(path, "/admin/static/") ||
			path == "/admin/manifest.webmanifest" ||
			path == "/admin/sw.js" ||
			strings.HasPrefix(path, "/admin/icons/") {
			c.Next()
			return
		}

		// For all other admin routes, prevent caching
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Next()
	}
}

// GinSessionMiddleware validates session and redirects to login if needed
func GinSessionMiddleware(queries *admindb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for login routes, static files, PWA files
		path := c.Request.URL.Path
		if path == "/login" ||
			strings.HasPrefix(path, "/admin/static/") ||
			path == "/admin/manifest.webmanifest" ||
			path == "/admin/sw.js" ||
			strings.HasPrefix(path, "/admin/icons/") {
			c.Next()
			return
		}

		// Get session ID from cookie
		sessionID := getSessionID(c.Request)
		if sessionID == "" {
			slog.Info("No session found, redirecting to login",
				slog.String("path", c.Request.URL.Path))
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}

		// Validate session
		session, err := queries.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				slog.Info("Invalid session, redirecting to login", slog.String("sessionID", sessionID))
				c.Redirect(http.StatusFound, "/admin/login")
				c.Abort()
				return
			}
			slog.Error("Failed to get session", slog.String("error", err.Error()))
			c.String(http.StatusInternalServerError, "Internal Server Error")
			c.Abort()
			return
		}

		// Update last accessed time in background
		go func() {
			ctx := context.Background()
			if err := queries.UpdateSessionLastAccessed(ctx, sessionID); err != nil {
				slog.Error("Failed to update session last accessed", slog.String("error", err.Error()))
			}
		}()

		// Add username to gin context
		c.Set("username", session.Username)
		c.Next()
	}
}

// SetupAdminRoutes configures admin routes on the given router group
func SetupAdminRoutes(adminGroup *gin.RouterGroup, queries *admindb.Queries, sobsClient *sobs.SobsClient, cfg internal.Config) {
	// Initialize OG image service
	var ogImageService *ogimage.Service
	if cfg.OGImageEnabled {
		s3Adapter := ogimage.NewSobsAdapter(sobsClient)
		ogGenerator := ogimage.NewGenerator(s3Adapter, cfg.S3AttachmentsBaseUrl, cfg.OGImageFontPath)
		ogImageService = ogimage.NewService(ogGenerator, queries)
	}

	// Create handler
	handler := NewHtmxHandler(queries, sobsClient, cfg.AdminUser, cfg.AdminPassword, !cfg.LocalDev, cfg.S3AttachmentsBaseUrl, ogImageService)

	// Login page (no session middleware needed)
	adminGroup.GET("/login", handler.RenderLoginPage)
	adminGroup.POST("/login", handler.HandleLogin)
	adminGroup.POST("/api/login", handler.APILogin)

	// PWA files (served before session middleware for accessibility)
	adminGroup.StaticFile("/manifest.webmanifest", "admin/manifest.webmanifest")
	adminGroup.StaticFile("/sw.js", "admin/static/sw.js")
	adminGroup.Static("/icons", "admin/static/icons/")

	// Add middlewares for authenticated routes
	adminGroup.Use(NoCacheMiddleware())
	adminGroup.Use(GinSessionMiddleware(queries))

	// Web Share Target endpoint (requires authentication)
	adminGroup.POST("/share-target", handler.HandleShareTarget)

	// Root redirect to entries
	adminGroup.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/entries/search")
	})

	// Entry list routes
	adminGroup.GET("/entries", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/entries/search")
	})
	adminGroup.GET("/entries/search", handler.RenderEntriesPage)
	adminGroup.POST("/entries/create", handler.CreateEntry)

	// Entry routes with query parameter (supports both slug and date-based paths)
	// Examples: /entries/edit?path=getting-started, /entries/edit?path=2024/01/01/120000
	adminGroup.GET("/entries/edit", handler.RenderEntryEditPage)
	adminGroup.POST("/entries/title", handler.UpdateEntryTitle)
	adminGroup.POST("/entries/body", handler.UpdateEntryBody)
	adminGroup.POST("/entries/visibility", handler.UpdateEntryVisibility)
	adminGroup.POST("/entries/image/regenerate", handler.RegenerateEntryImage)
	adminGroup.POST("/entries/upload", handler.UploadEntryImage)
	adminGroup.DELETE("/entries/delete", handler.DeleteEntry)

	// JSON API routes (used by Preact apps)
	adminGroup.GET("/api/entries", handler.APIListEntries)
	adminGroup.POST("/api/entries/create", handler.APICreateEntry)
	adminGroup.PUT("/api/entries/title", handler.APIUpdateTitle)
	adminGroup.PUT("/api/entries/body", handler.APIUpdateBody)
	adminGroup.PUT("/api/entries/visibility", handler.APIUpdateVisibility)
	adminGroup.DELETE("/api/entries/delete", handler.APIDeleteEntry)
	adminGroup.POST("/api/entries/image/regenerate", handler.APIRegenerateEntryImage)
	adminGroup.POST("/api/entries/preview", handler.APIPreviewMarkdown)

	// Static files
	adminGroup.Static("/static", "admin/static/")
}
