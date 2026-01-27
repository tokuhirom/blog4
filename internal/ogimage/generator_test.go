package ogimage

import (
	"context"
	"image/png"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockS3Uploader is a mock implementation of S3Uploader for testing
type mockS3Uploader struct {
	uploadedKey         string
	uploadedContentType string
	uploadedMetadata    map[string]string
	uploadError         error
}

func (m *mockS3Uploader) PutObject(ctx context.Context, key string, body io.Reader, contentType string, metadata map[string]string) error {
	if m.uploadError != nil {
		return m.uploadError
	}
	m.uploadedKey = key
	m.uploadedContentType = contentType
	m.uploadedMetadata = metadata
	// Read body to simulate upload
	_, _ = io.ReadAll(body)
	return nil
}

func TestRenderImage(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{"short_english", "Hello World"},
		{"short_japanese", "こんにちは"},
		{"long_japanese", "これは非常に長い日本語のタイトルです。自動的に折り返されて複数行になることを期待します。最大3行まで表示されます。"},
		{"mixed", "Goで実装するOGP画像生成"},
		{"emoji", "🎉 新機能リリース 🚀"},
		{"long_english", "This is a very long title that should wrap across multiple lines automatically. We expect it to display up to three lines maximum."},
	}

	generator := NewGenerator(nil, "https://example.com", "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf, err := generator.renderImage(tt.title, time.Now())
			require.NoError(t, err)
			require.NotNil(t, buf)

			// Verify it's a valid PNG
			img, err := png.Decode(buf)
			require.NoError(t, err)
			assert.NotNil(t, img)

			// Verify dimensions
			bounds := img.Bounds()
			assert.Equal(t, imageWidth, bounds.Dx(), "width should be 1200")
			assert.Equal(t, imageHeight, bounds.Dy(), "height should be 630")
		})
	}
}

func TestGenerateOGImage(t *testing.T) {
	mockS3 := &mockS3Uploader{}
	generator := NewGenerator(mockS3, "https://example.com", "")

	entry := EntryInfo{
		Path:        "2026/01/test-entry",
		Title:       "テストエントリー",
		PublishedAt: time.Date(2026, 1, 27, 12, 0, 0, 0, time.UTC),
	}

	ctx := context.Background()
	url, err := generator.GenerateOGImage(ctx, entry)
	require.NoError(t, err)
	assert.NotEmpty(t, url)

	// Verify S3 upload
	assert.Contains(t, mockS3.uploadedKey, "og-images/")
	assert.Equal(t, "image/png", mockS3.uploadedContentType)
	assert.Equal(t, "blog4-og-generator", mockS3.uploadedMetadata["generated-by"])
	assert.Equal(t, "2026/01/test-entry", mockS3.uploadedMetadata["entry-path"])
	assert.Equal(t, "UTF-8", mockS3.uploadedMetadata["charset"])

	// Verify URL format
	assert.Contains(t, url, "https://example.com/og-images/")
}

func TestSanitizeTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal", "Hello World", "Hello World"},
		{"with_newline", "Hello\nWorld", "Hello World"},
		{"with_tab", "Hello\tWorld", "Hello World"},
		{"with_control_chars", "Hello\x00\x01World", "HelloWorld"},
		{"extra_spaces", "  Hello   World  ", "Hello   World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeTitle(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test very long title separately
	t.Run("very_long", func(t *testing.T) {
		// Create a string with 250 'a' characters
		longString := strings.Repeat("a", 250)
		result := sanitizeTitle(longString)
		assert.Equal(t, 200, len([]rune(result)))
	})
}

func TestWrapText(t *testing.T) {
	generator := NewGenerator(nil, "https://example.com", "")

	// Load font for testing
	fnt, err := generator.loadFont()
	require.NoError(t, err)

	// Create a context for measuring
	dc := newTestContext(fnt)

	tests := []struct {
		name           string
		text           string
		maxLines       int
		expectMaxLines int
		expectEllipsis bool
	}{
		{"short", "Short", 3, 1, false},
		{"medium", "This is a medium length title that should fit", 3, 3, false},
		{"long", "This is a very long title that should wrap across multiple lines automatically and demonstrate the wrapping functionality", 3, 3, true},
		{"japanese", "これは日本語のタイトルです", 3, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := generator.wrapText(dc, tt.text, titleMaxWidth, tt.maxLines)
			assert.LessOrEqual(t, len(lines), tt.expectMaxLines)
			if tt.expectEllipsis {
				assert.Contains(t, lines[len(lines)-1], "...")
			}
		})
	}
}

// newTestContext creates a gg context for testing text wrapping
func newTestContext(fnt *truetype.Font) *gg.Context {
	dc := gg.NewContext(imageWidth, imageHeight)
	titleFace := truetype.NewFace(fnt, &truetype.Options{Size: titleFontSize})
	dc.SetFontFace(titleFace)
	return dc
}
