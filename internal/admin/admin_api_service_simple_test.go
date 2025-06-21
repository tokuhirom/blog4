package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test the utility functions that don't require mocks

func TestGetDefaultTitle(t *testing.T) {
	title := getDefaultTitle()
	// Should be in format YYYYMMDDHHMMSS
	assert.Regexp(t, `^\d{14}$`, title)
}

func TestExtractLinks(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected []string
	}{
		{
			name:     "Single link",
			markdown: "This is a [[Test Link]] in the text.",
			expected: []string{"Test Link"},
		},
		{
			name:     "Multiple links",
			markdown: "Here are [[Link 1]] and [[Link 2]] in the text.",
			expected: []string{"Link 1", "Link 2"},
		},
		{
			name:     "Duplicate links - removes duplicates",
			markdown: "[[Test]] and [[Test]] and [[Test]]",
			expected: []string{"Test"},
		},
		{
			name:     "Links with spaces",
			markdown: "[[  Spaced Link  ]] should be trimmed",
			expected: []string{"Spaced Link"},
		},
		{
			name:     "No links",
			markdown: "No links in this text",
			expected: nil,
		},
		{
			name:     "Empty link",
			markdown: "[[]] should not be included",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractLinks(tt.markdown)
			assert.Equal(t, tt.expected, result)
		})
	}
}