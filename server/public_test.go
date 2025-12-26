package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
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
		c.String(http.StatusOK, "entry: "+entryPath)
	})

	tests := []struct {
		name           string
		path           string
		expectRoute    string
		expectParam    string
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

// TestPublicRouterWithActualChiMount tests using the actual chi router
func TestPublicRouterWithActualChiMount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Import chi to test actual chi.Mount behavior
	chiRouter := chi.NewRouter()

	// Create the public router (mock version)
	publicRouter := gin.New()
	publicRouter.Use(gin.Recovery())

	called := make(map[string]bool)

	publicRouter.GET("/", func(c *gin.Context) {
		called["/"] = true
		c.String(http.StatusOK, "top page")
	})

	publicRouter.GET("/entry/*filepath", func(c *gin.Context) {
		called["/entry/*"] = true
		c.String(http.StatusOK, "entry page: "+c.Param("filepath"))
	})

	// Use actual chi.Mount (this is what's used in production)
	chiRouter.Mount("/", publicRouter)

	tests := []struct {
		name        string
		path        string
		shouldWork  bool
	}{
		{
			name:       "Root should work",
			path:       "/",
			shouldWork: true,
		},
		{
			name:       "Entry path should work",
			path:       "/entry/test",
			shouldWork: true,
		},
		{
			name:       "Entry path with date should work",
			path:       "/entry/2024/01/01/120000",
			shouldWork: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called = make(map[string]bool)

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			chiRouter.ServeHTTP(w, req)

			if tt.shouldWork && w.Code == http.StatusNotFound {
				t.Errorf("Path %s returned 404 when it should work", tt.path)
				t.Logf("Response: %s", w.Body.String())
				t.Logf("Called: %v", called)
			} else {
				t.Logf("✓ Path %s: status=%d, called=%v", tt.path, w.Code, called)
			}
		})
	}
}

// TestActualProductionRouterSetup simulates the exact production router setup
func TestActualProductionRouterSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create main chi router (like in router/router.go)
	r := chi.NewRouter()

	// Create admin router (gin-based, mounted at /admin)
	adminRouter := gin.New()
	adminRouter.Use(gin.Recovery())

	// Admin routes (simplified)
	adminRouter.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/entries")
	})
	adminRouter.GET("/entries", func(c *gin.Context) {
		c.String(http.StatusOK, "admin entries")
	})

	// Wrap admin router to strip /admin prefix (like admin does now)
	adminHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Strip /admin prefix for gin router
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/admin")
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
		adminRouter.ServeHTTP(w, req)
	})

	// Mount admin at /admin
	r.Mount("/admin", adminHandler)

	// Create public router (gin-based, mounted at /)
	publicRouter := gin.New()
	publicRouter.Use(gin.Recovery())

	publicCalled := make(map[string]bool)

	publicRouter.GET("/", func(c *gin.Context) {
		publicCalled["/"] = true
		c.String(http.StatusOK, "top page")
	})

	publicRouter.GET("/entry/*filepath", func(c *gin.Context) {
		publicCalled["/entry/*"] = true
		c.String(http.StatusOK, "entry page: "+c.Param("filepath"))
	})

	// Mount public at / (THIS IS THE KEY - after admin is mounted)
	r.Mount("/", publicRouter)

	// Add other routes
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	tests := []struct {
		name        string
		path        string
		expectCode  int
		description string
	}{
		{
			name:        "Admin root",
			path:        "/admin",
			expectCode:  http.StatusFound,
			description: "Should redirect to /admin/entries",
		},
		{
			name:        "Admin entries",
			path:        "/admin/entries",
			expectCode:  http.StatusOK,
			description: "Should show admin entries",
		},
		{
			name:        "Healthz",
			path:        "/healthz",
			expectCode:  http.StatusOK,
			description: "Should return ok",
		},
		{
			name:        "Public root",
			path:        "/",
			expectCode:  http.StatusOK,
			description: "Should show top page",
		},
		{
			name:        "Public entry page",
			path:        "/entry/getting-started",
			expectCode:  http.StatusOK,
			description: "Should show entry page - THIS IS THE KEY TEST",
		},
		{
			name:        "Public entry with date path",
			path:        "/entry/2024/01/01/120000",
			expectCode:  http.StatusOK,
			description: "Should show entry page with date path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			publicCalled = make(map[string]bool)

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("%s: expected status %d, got %d", tt.description, tt.expectCode, w.Code)
				t.Logf("Response body: %s", w.Body.String())
				t.Logf("Public called: %v", publicCalled)
			} else {
				t.Logf("✓ %s (status=%d, public_called=%v)", tt.description, w.Code, publicCalled)
			}
		})
	}
}
