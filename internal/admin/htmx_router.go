package admin

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
)

// GinSessionMiddleware validates session and redirects to login if needed
func GinSessionMiddleware(queries *admindb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for login routes
		if c.Request.URL.Path == "/login" {
			c.Next()
			return
		}

		// Get session ID from cookie
		sessionID := getSessionID(c.Request)
		if sessionID == "" {
			slog.Info("No session found, redirecting to login", slog.String("path", c.Request.URL.Path))
			c.Redirect(http.StatusFound, "/admin/htmx/login")
			c.Abort()
			return
		}

		// Validate session
		session, err := queries.GetSession(c.Request.Context(), sessionID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				slog.Info("Invalid session, redirecting to login", slog.String("sessionID", sessionID))
				c.Redirect(http.StatusFound, "/admin/htmx/login")
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

// SetupHtmxRouter creates and configures the htmx router using gin
func SetupHtmxRouter(queries *admindb.Queries, cfg server.Config) http.Handler {
	// Create gin router in release mode for production
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Create handler
	handler := NewHtmxHandler(queries, cfg.AdminUser, cfg.AdminPassword, !cfg.LocalDev)

	// Login page (no session middleware needed)
	router.GET("/login", handler.RenderLoginPage)
	router.POST("/login", handler.HandleLogin)

	// Add session middleware for authenticated routes
	router.Use(GinSessionMiddleware(queries))

	// Entry list routes
	router.GET("/entries", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/htmx/entries/search")
	})
	router.GET("/entries/search", handler.RenderEntriesPage)
	router.POST("/entries/create", handler.CreateEntry)

	// Entry routes with path parameters (yyyy/mm/dd/hhmmss format)
	entries := router.Group("/entries/:year/:month/:day/:time")
	{
		entries.GET("/edit", handler.RenderEntryEditPage)
		entries.POST("/title", handler.UpdateEntryTitle)
		entries.POST("/body", handler.UpdateEntryBody)
		entries.POST("/visibility", handler.UpdateEntryVisibility)
		entries.POST("/image/regenerate", handler.RegenerateEntryImage)
		entries.DELETE("", handler.DeleteEntry)
	}

	// Static files
	router.Static("/static", "web/static/admin")

	return router
}
