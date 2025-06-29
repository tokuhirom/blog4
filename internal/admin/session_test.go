package admin

import (
	"testing"
)

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
