package ogimage

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	imageWidth    = 1200
	imageHeight   = 630
	padding       = 60.0
	logoSize      = 120.0
	brandFontSize = 36.0
	titleFontSize = 48.0
	dateFontSize  = 24.0
	titleMaxWidth = 900.0
	titleMaxLines = 3
	lineHeight    = 1.3
)

// S3Uploader defines the interface for uploading files to S3
type S3Uploader interface {
	PutObject(ctx context.Context, key string, body io.Reader, contentType string, metadata map[string]string) error
}

// EntryInfo contains the information needed to generate an OG image
type EntryInfo struct {
	Path        string
	Title       string
	PublishedAt time.Time
}

// Generator generates OG images
type Generator struct {
	s3Client  S3Uploader
	s3BaseURL string
	fontPath  string
	fontCache *truetype.Font
	fontOnce  sync.Once
}

// NewGenerator creates a new Generator
func NewGenerator(s3Client S3Uploader, s3BaseURL string, fontPath string) *Generator {
	return &Generator{
		s3Client:  s3Client,
		s3BaseURL: s3BaseURL,
		fontPath:  fontPath,
	}
}

// GenerateOGImage generates an OG image and uploads it to S3, returning the URL
func (g *Generator) GenerateOGImage(ctx context.Context, entry EntryInfo) (string, error) {
	// Render the image
	buf, err := g.renderImage(entry.Title, entry.PublishedAt)
	if err != nil {
		return "", fmt.Errorf("failed to render image: %w", err)
	}

	// Generate S3 key with timestamp
	now := time.Now()
	key := fmt.Sprintf("og-images/%04d/%02d/%02d/%s.png",
		now.Year(), now.Month(), now.Day(), now.Format("150405"))

	// Upload to S3
	metadata := map[string]string{
		"generated-by": "blog4-og-generator",
		"entry-path":   entry.Path,
		"charset":      "UTF-8",
	}

	err = g.s3Client.PutObject(ctx, key, bytes.NewReader(buf.Bytes()), "image/png", metadata)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct URL
	url := fmt.Sprintf("%s/%s", g.s3BaseURL, key)
	return url, nil
}

// renderImage creates a 1200x630px PNG image
func (g *Generator) renderImage(title string, publishedAt time.Time) (*bytes.Buffer, error) {
	// Create canvas
	dc := gg.NewContext(imageWidth, imageHeight)

	// Draw gradient background (#1976d2 -> #1565c0)
	gradient := gg.NewLinearGradient(0, 0, 0, imageHeight)
	gradient.AddColorStop(0, color.RGBA{R: 0x19, G: 0x76, B: 0xd2, A: 0xff})
	gradient.AddColorStop(1, color.RGBA{R: 0x15, G: 0x65, B: 0xc0, A: 0xff})
	dc.SetFillStyle(gradient)
	dc.DrawRectangle(0, 0, imageWidth, imageHeight)
	dc.Fill()

	// Load font
	fnt, err := g.loadFont()
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %w", err)
	}

	// Draw logo "B4"
	logoFace := truetype.NewFace(fnt, &truetype.Options{Size: logoSize})
	dc.SetFontFace(logoFace)
	dc.SetColor(color.White)
	dc.DrawString("B4", padding, padding+logoSize)

	// Draw brand name "Blog4"
	brandFace := truetype.NewFace(fnt, &truetype.Options{Size: brandFontSize})
	dc.SetFontFace(brandFace)
	dc.DrawString("Blog4", padding+logoSize+20, padding+logoSize-20)

	// Sanitize and wrap title
	sanitizedTitle := sanitizeTitle(title)
	titleFace := truetype.NewFace(fnt, &truetype.Options{Size: titleFontSize, Hinting: font.HintingFull})
	dc.SetFontFace(titleFace)
	lines := g.wrapText(dc, sanitizedTitle, titleMaxWidth, titleMaxLines)

	// Calculate title Y position (center vertically)
	titleHeight := float64(len(lines)) * titleFontSize * lineHeight
	titleY := (imageHeight - titleHeight) / 2

	// Draw title lines (centered horizontally)
	for i, line := range lines {
		lineY := titleY + float64(i)*titleFontSize*lineHeight
		w, _ := dc.MeasureString(line)
		x := (imageWidth - w) / 2
		dc.DrawString(line, x, lineY)
	}

	// Draw date (bottom right)
	dateFace := truetype.NewFace(fnt, &truetype.Options{Size: dateFontSize})
	dc.SetFontFace(dateFace)
	dateStr := fmt.Sprintf("Published: %s", publishedAt.Format("2006-01-02"))
	dateW, _ := dc.MeasureString(dateStr)
	dc.DrawString(dateStr, imageWidth-dateW-padding, imageHeight-padding)

	// Encode to PNG
	var buf bytes.Buffer
	err = dc.EncodePNG(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}

	return &buf, nil
}

// wrapText wraps text to fit within maxWidth, returning at most maxLines lines
func (g *Generator) wrapText(dc *gg.Context, text string, maxWidth float64, maxLines int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " " + word
		} else {
			testLine = word
		}

		w, _ := dc.MeasureString(testLine)
		if w > maxWidth {
			// Current line is too long
			if currentLine == "" {
				// Single word is too long, need to break it
				lines = append(lines, g.breakWord(dc, word, maxWidth))
				if len(lines) >= maxLines {
					break
				}
			} else {
				// Save current line and start new one with current word
				lines = append(lines, currentLine)
				if len(lines) >= maxLines {
					break
				}
				currentLine = word
			}
		} else {
			currentLine = testLine
		}
	}

	// Add remaining line
	if currentLine != "" && len(lines) < maxLines {
		lines = append(lines, currentLine)
	}

	// Truncate last line with ellipsis if we hit maxLines
	if len(lines) == maxLines {
		lastLine := lines[maxLines-1]
		w, _ := dc.MeasureString(lastLine + "...")
		if w > maxWidth {
			// Need to shorten the last line
			for utf8.RuneCountInString(lastLine) > 0 {
				runes := []rune(lastLine)
				lastLine = string(runes[:len(runes)-1])
				w, _ := dc.MeasureString(lastLine + "...")
				if w <= maxWidth {
					break
				}
			}
		}
		lines[maxLines-1] = lastLine + "..."
	}

	return lines
}

// breakWord breaks a single word that's too long to fit on one line
func (g *Generator) breakWord(dc *gg.Context, word string, maxWidth float64) string {
	runes := []rune(word)
	for i := len(runes); i > 0; i-- {
		substr := string(runes[:i])
		w, _ := dc.MeasureString(substr + "...")
		if w <= maxWidth {
			return substr + "..."
		}
	}
	return "..."
}

// loadFont loads the Japanese font with fallback options
func (g *Generator) loadFont() (*truetype.Font, error) {
	var loadErr error
	g.fontOnce.Do(func() {
		// Try primary font path (Noto Sans CJK)
		if g.fontPath != "" {
			fnt, err := loadFontFromFile(g.fontPath)
			if err == nil {
				g.fontCache = fnt
				slog.Info("Loaded font", slog.String("path", g.fontPath))
				return
			}
			slog.Warn("Failed to load primary font", slog.String("path", g.fontPath), slog.Any("error", err))
		}

		// Fallback fonts: prioritize .ttf files over .ttc files
		fallbackPaths := []string{
			"/usr/share/fonts/opentype/ipafont-gothic/ipagp.ttf",  // IPA P Gothic (TTF)
			"/usr/share/fonts/opentype/ipafont-gothic/ipag.ttf",   // IPA Gothic (TTF)
			"/usr/share/fonts/truetype/ipafont/ipagp.ttf",         // IPA P Gothic (alt path)
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Bold.ttc", // Noto CJK (TTC)
			"/usr/share/fonts/truetype/noto/NotoSansCJK-Bold.ttc", // Noto CJK (alt path)
			"/System/Library/Fonts/Hiragino Sans GB.ttc",          // macOS
		}

		for _, path := range fallbackPaths {
			fnt, err := loadFontFromFile(path)
			if err == nil {
				g.fontCache = fnt
				slog.Info("Loaded fallback font", slog.String("path", path))
				return
			}
		}

		// Final fallback: embedded Go font (no Japanese support, but better than nothing)
		fnt, err := truetype.Parse(goregular.TTF)
		if err == nil {
			g.fontCache = fnt
			slog.Warn("Using embedded font (no Japanese support)")
			return
		}

		loadErr = fmt.Errorf("failed to load any font")
	})

	if loadErr != nil {
		return nil, loadErr
	}

	if g.fontCache == nil {
		return nil, fmt.Errorf("font cache is nil")
	}

	return g.fontCache, nil
}

// loadFontFromFile loads a TrueType font from a file
func loadFontFromFile(path string) (*truetype.Font, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	fnt, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	return fnt, nil
}

// sanitizeTitle removes control characters and limits length
func sanitizeTitle(title string) string {
	// Remove control characters except newline
	var cleaned strings.Builder
	for _, r := range title {
		if r == '\n' || r == '\r' || r == '\t' {
			cleaned.WriteRune(' ')
		} else if r >= 32 || r == '\t' {
			cleaned.WriteRune(r)
		}
	}

	result := cleaned.String()

	// Limit to 200 characters
	runes := []rune(result)
	if len(runes) > 200 {
		result = string(runes[:200])
	}

	return strings.TrimSpace(result)
}
