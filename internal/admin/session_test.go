package admin

import (
	"testing"
	"time"
)

func TestSessionTimeouts(t *testing.T) {
	tests := []struct {
		name            string
		rememberMe      bool
		expectedTimeout time.Duration
	}{
		{
			name:            "regular session",
			rememberMe:      false,
			expectedTimeout: defaultSessionTimeout,
		},
		{
			name:            "remember me session",
			rememberMe:      true,
			expectedTimeout: extendedSessionTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test verifies the timeout values
			if tt.rememberMe && tt.expectedTimeout != 30*24*time.Hour {
				t.Errorf("Remember me session should have 30 day timeout, got %v", tt.expectedTimeout)
			}
			if !tt.rememberMe && tt.expectedTimeout != 24*time.Hour {
				t.Errorf("Regular session should have 24 hour timeout, got %v", tt.expectedTimeout)
			}
		})
	}
}

func TestGenerateSessionID(t *testing.T) {
	// Test that session IDs are generated correctly
	id1, err := generateSessionID()
	if err != nil {
		t.Fatalf("Failed to generate session ID: %v", err)
	}

	id2, err := generateSessionID()
	if err != nil {
		t.Fatalf("Failed to generate session ID: %v", err)
	}

	// Verify IDs are not empty
	if id1 == "" || id2 == "" {
		t.Error("Generated session IDs should not be empty")
	}

	// Verify IDs are unique
	if id1 == id2 {
		t.Error("Generated session IDs should be unique")
	}

	// Verify ID length (base64 encoded 32 bytes should be 44 chars)
	expectedLen := 44 // base64 encoding of 32 bytes
	if len(id1) != expectedLen {
		t.Errorf("Session ID length should be %d, got %d", expectedLen, len(id1))
	}
}
