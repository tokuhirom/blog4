package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestGinEntryRouteMatching tests if gin's wildcard route /entry/*filepath works correctly
func TestGinEntryRouteMatching(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a simple gin router with the same route structure as the public router
	r := gin.New()
	r.Use(gin.Recovery())

	// Track which routes were called
	rootCalled := false
	entryCalled := false
	entryPath := ""

	r.GET("/", func(c *gin.Context) {
		rootCalled = true
		c.String(http.StatusOK, "root")
	})

	r.GET("/entry/*filepath", func(c *gin.Context) {
		entryCalled = true
		entryPath = c.Param("filepath")
		// Strip leading slash (same as production code)
		strippedPath := strings.TrimPrefix(entryPath, "/")
		c.String(http.StatusOK, "entry: "+strippedPath)
	})

	tests := []struct {
		name        string
		path        string
		expectRoute string
		expectParam string
	}{
		{
			name:        "Root path",
			path:        "/",
			expectRoute: "root",
		},
		{
			name:        "Entry with simple path",
			path:        "/entry/getting-started",
			expectRoute: "entry",
			expectParam: "/getting-started",
		},
		{
			name:        "Entry with date-based path",
			path:        "/entry/2024/01/01/120000",
			expectRoute: "entry",
			expectParam: "/2024/01/01/120000",
		},
		{
			name:        "Entry with complex path",
			path:        "/entry/path/with/multiple/slashes",
			expectRoute: "entry",
			expectParam: "/path/with/multiple/slashes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			rootCalled = false
			entryCalled = false
			entryPath = ""

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code == http.StatusNotFound {
				t.Errorf("Path %s returned 404 - route not matched", tt.path)
				t.Logf("Response body: %s", w.Body.String())
				return
			}

			if w.Code != http.StatusOK {
				t.Errorf("Path %s returned status %d, expected 200", tt.path, w.Code)
				return
			}

			if tt.expectRoute == "root" && !rootCalled {
				t.Errorf("Path %s: expected root route to be called", tt.path)
			}

			if tt.expectRoute == "entry" {
				if !entryCalled {
					t.Errorf("Path %s: expected entry route to be called", tt.path)
				}
				if entryPath != tt.expectParam {
					t.Errorf("Path %s: expected param %q, got %q", tt.path, tt.expectParam, entryPath)
				}
			}

			t.Logf("✓ Path %s correctly routed to %s", tt.path, tt.expectRoute)
		})
	}
}
