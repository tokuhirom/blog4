package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/server"
)

// CheckWebAccelGuard is a Gin middleware that checks for the X-WebAccel-Guard header
func CheckWebAccelGuard(cfg server.Config) func(c *gin.Context) {
	return func(c *gin.Context) {
		if c.Request.URL.Path != "/healthz" {
			gotToken := c.GetHeader("X-WebAccel-Guard")
			if gotToken != cfg.WebAccelGuard {
				slog.Warn("invalid X-WebAccel-Guard header",
					slog.String("got_token", gotToken),
					slog.String("path", c.Request.URL.Path))
				c.String(http.StatusBadRequest, "Invalid X-WebAccel-Guard header")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
