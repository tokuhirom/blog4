package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tokuhirom/blog4/db/admin/admindb"
)

// GinSessionMiddleware converts the session middleware to gin middleware
func GinSessionMiddleware(queries *admindb.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the underlying session middleware
		sessionMiddleware := SessionMiddleware(queries)

		// Wrap the gin handler to work with the session middleware
		sessionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Update the gin context with the modified request
			c.Request = r
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}

// SetupHtmxRouter creates and configures the htmx router using gin
func SetupHtmxRouter(queries *admindb.Queries) http.Handler {
	// Create gin router in release mode for production
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add recovery middleware
	router.Use(gin.Recovery())

	// Add session middleware
	router.Use(GinSessionMiddleware(queries))

	// Create handler
	handler := NewHtmxHandler(queries)

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
