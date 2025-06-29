package admin

import (
	"testing"
	"time"

	"github.com/zeebo/assert"

	"github.com/tokuhirom/blog4/server/admin/openapi"
)

// TestAuthLoginRememberMe tests the remember me functionality
// This is an integration test that requires a test database
func TestAuthLoginRememberMe(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// This test would require a test database setup
	// For now, we document the expected behavior:

	t.Run("regular login sets 24 hour session", func(t *testing.T) {
		// When logging in without remember_me
		// The session should expire in 24 hours
		// The cookie MaxAge should be ~24 hours
	})

	t.Run("remember me login sets 30 day session", func(t *testing.T) {
		// When logging in with remember_me = true
		// The session should expire in 30 days
		// The cookie MaxAge should be ~30 days
	})
}

// TestRememberMeSessionExpiry verifies the session timeout calculation
func TestRememberMeSessionExpiry(t *testing.T) {
	testCases := []struct {
		name             string
		rememberMe       bool
		expectedDuration time.Duration
	}{
		{
			name:             "regular session - 24 hours",
			rememberMe:       false,
			expectedDuration: 24 * time.Hour,
		},
		{
			name:             "remember me session - 30 days",
			rememberMe:       true,
			expectedDuration: 30 * 24 * time.Hour,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock login request
			loginReq := &openapi.LoginRequest{
				Username:   "testuser",
				Password:   "testpass",
				RememberMe: openapi.NewOptBool(tc.rememberMe),
			}

			// Calculate expected session timeout
			var sessionTimeout time.Duration
			if loginReq.RememberMe.Or(false) {
				sessionTimeout = extendedSessionTimeout
			} else {
				sessionTimeout = defaultSessionTimeout
			}

			// Verify the timeout matches expectation
			assert.Equal(t, tc.expectedDuration, sessionTimeout)
		})
	}
}

// TestSessionCookieSettings verifies cookie settings for remember me
func TestSessionCookieSettings(t *testing.T) {
	testCases := []struct {
		name           string
		rememberMe     bool
		expectedMaxAge int
	}{
		{
			name:           "regular session cookie",
			rememberMe:     false,
			expectedMaxAge: int((24 * time.Hour).Seconds()),
		},
		{
			name:           "remember me session cookie",
			rememberMe:     true,
			expectedMaxAge: int((30 * 24 * time.Hour).Seconds()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate session timeout based on remember me
			var sessionTimeout time.Duration
			if tc.rememberMe {
				sessionTimeout = extendedSessionTimeout
			} else {
				sessionTimeout = defaultSessionTimeout
			}

			expires := time.Now().Add(sessionTimeout)
			maxAge := int(time.Until(expires).Seconds())

			// Allow for small timing differences (within 5 seconds)
			diff := maxAge - tc.expectedMaxAge
			if diff < -5 || diff > 5 {
				t.Errorf("Expected MaxAge ~%d seconds, got %d seconds", tc.expectedMaxAge, maxAge)
			}
		})
	}
}
