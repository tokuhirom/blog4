package admin

import (
	"testing"
	"time"
)

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
