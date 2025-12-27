package admin

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

type contextKey string

const usernameKey contextKey = "username"

// SessionMiddleware creates a middleware that checks for valid sessions
func SessionMiddleware(queries *admindb.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip auth check for auth endpoints
			if strings.HasPrefix(r.URL.Path, "/admin/api/auth/") {
				next.ServeHTTP(w, r)
				return
			}

			// Get session ID from cookie
			sessionID := getSessionID(r)
			if sessionID == "" {
				slog.Info("No session found", slog.String("path", r.URL.Path))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Validate session
			session, err := queries.GetSession(r.Context(), sessionID)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					slog.Info("Invalid session", slog.String("sessionID", sessionID))
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				slog.Error("Failed to get session", slog.String("error", err.Error()))
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Update last accessed time in background
			go func() {
				ctx := context.Background()
				if err := queries.UpdateSessionLastAccessed(ctx, sessionID); err != nil {
					slog.Error("Failed to update session last accessed", slog.String("error", err.Error()))
				}
			}()

			// Add username to context
			ctx := context.WithValue(r.Context(), usernameKey, session.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
