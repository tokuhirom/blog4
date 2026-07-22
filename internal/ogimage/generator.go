package ogimage

import (
	"bytes"
	"context"
	"fmt"
	"image/color"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/fogleman/gg"
	xfont "golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

const (
	imageWidth  = 1200
	imageHeight = 630

	marginX = 80.0 // left/right content margin
	barW    = 12.0 // left accent bar width

	siteFontSize  = 30.0
	titleFontSize = 60.0
	metaFontSize  = 26.0

	titleLineHeight = 82.0 // baseline-to-baseline for the title
	titleMaxLines   = 3

	siteY      = 100.0               // baseline of the site name
	headerZone = 130.0               // bottom of the header area
	footerZone = imageHeight - 100.0 // top of the footer area
	footerY    = imageHeight - 60.0  // baseline of the footer text
	dpi        = 72.0
)

// Palette (Tailwind-ish slate/sky).
var (
	bgTop     = color.RGBA{0x0f, 0x17, 0x2a, 0xff} // slate-900
	bgBottom  = color.RGBA{0x1e, 0x29, 0x3b, 0xff} // slate-800
	accentBar = color.RGBA{0x38, 0xbd, 0xf8, 0xff} // sky-400
	siteColor = color.RGBA{0x94, 0xa3, 0xb8, 0xff} // slate-400
	metaColor = color.RGBA{0x64, 0x74, 0x8b, 0xff} // slate-500
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

// Config configures a Generator.
type Config struct {
	S3Client  S3Uploader
	S3BaseURL string
	// FontPath is an optional explicit font override (used as a fallback if the
	// bundled Noto Sans CJK fonts are unavailable).
	FontPath string
	// SiteName is shown as the header label (e.g. "tokuhirom's blog").
	SiteName string
	// SiteURL is used to derive the domain label shown in the footer.
	SiteURL string
}

// Generator generates OG images
type Generator struct {
	s3Client  S3Uploader
	s3BaseURL string
	fontPath  string
	siteName  string
	siteHost  string

	fontOnce    sync.Once
	fontErr     error
	boldFont    *sfnt.Font
	mediumFont  *sfnt.Font
	regularFont *sfnt.Font
}

// NewGenerator creates a new Generator
func NewGenerator(cfg Config) *Generator {
	siteName := cfg.SiteName
	if siteName == "" {
		siteName = "blog"
	}
	host := cfg.SiteURL
	if u, err := url.Parse(cfg.SiteURL); err == nil && u.Host != "" {
		host = u.Host
	}
	return &Generator{
		s3Client:  cfg.S3Client,
		s3BaseURL: cfg.S3BaseURL,
		fontPath:  cfg.FontPath,
		siteName:  siteName,
		siteHost:  host,
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

// renderImage creates a 1200x630px PNG image (design "A": navy, left-aligned).
func (g *Generator) renderImage(title string, publishedAt time.Time) (*bytes.Buffer, error) {
	if err := g.loadFonts(); err != nil {
		return nil, fmt.Errorf("failed to load fonts: %w", err)
	}

	dc := gg.NewContext(imageWidth, imageHeight)

	// Background: subtle vertical gradient.
	grad := gg.NewLinearGradient(0, 0, 0, imageHeight)
	grad.AddColorStop(0, bgTop)
	grad.AddColorStop(1, bgBottom)
	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, imageWidth, imageHeight)
	dc.Fill()

	// Left accent bar.
	dc.SetColor(accentBar)
	dc.DrawRectangle(0, 0, barW, imageHeight)
	dc.Fill()

	// Header: site name.
	dc.SetFontFace(g.face(g.mediumFont, siteFontSize))
	dc.SetColor(siteColor)
	dc.DrawString(g.siteName, marginX, siteY)

	// Title: bold, left-aligned, vertically centered in the content zone.
	dc.SetFontFace(g.face(g.boldFont, titleFontSize))
	lines := wrapText(dc, sanitizeTitle(title), imageWidth-marginX*2, titleMaxLines)
	blockH := float64(len(lines)) * titleLineHeight
	// Center the block between header and footer, then offset to the first baseline.
	startY := headerZone + (footerZone-headerZone-blockH)/2 + titleFontSize
	dc.SetColor(color.White)
	for i, ln := range lines {
		dc.DrawString(ln, marginX, startY+float64(i)*titleLineHeight)
	}

	// Footer: domain (left) and date (right).
	dc.SetFontFace(g.face(g.regularFont, metaFontSize))
	dc.SetColor(metaColor)
	dc.DrawString(g.siteHost, marginX, footerY)
	dateStr := publishedAt.Format("2006-01-02")
	dw, _ := dc.MeasureString(dateStr)
	dc.DrawString(dateStr, imageWidth-marginX-dw, footerY)

	var buf bytes.Buffer
	if err := dc.EncodePNG(&buf); err != nil {
		return nil, fmt.Errorf("failed to encode PNG: %w", err)
	}
	return &buf, nil
}

// face builds a font.Face at the given size from a parsed sfnt.Font.
func (g *Generator) face(f *sfnt.Font, size float64) xfont.Face {
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		// Should not happen for an already-parsed font; fall back to a basic face.
		slog.Warn("failed to build font face", slog.Any("error", err))
	}
	return face
}

// wrapText wraps text to fit within maxWidth, returning at most maxLines lines.
// It tokenizes into CJK single characters and runs of non-CJK (word-like)
// characters so both Japanese and English wrap reasonably.
func wrapText(dc *gg.Context, text string, maxWidth float64, maxLines int) []string {
	tokens := tokenize(text)
	if len(tokens) == 0 {
		return []string{""}
	}

	var lines []string
	var cur string
	truncated := false
	for _, tok := range tokens {
		candidate := cur + tok
		w, _ := dc.MeasureString(candidate)
		if w > maxWidth && cur != "" {
			lines = append(lines, cur)
			if len(lines) == maxLines {
				// Ran out of lines but tokens remain: mark for ellipsis.
				truncated = true
				break
			}
			cur = strings.TrimLeft(tok, " ")
		} else {
			cur = candidate
		}
	}
	if !truncated && cur != "" && len(lines) < maxLines {
		lines = append(lines, cur)
	}

	// If content remains (we hit maxLines), ellipsize the last line.
	if truncated && len(lines) == maxLines {
		last := []rune(strings.TrimRight(lines[maxLines-1], " "))
		for len(last) > 0 {
			w, _ := dc.MeasureString(string(last) + "…")
			if w <= maxWidth {
				break
			}
			last = last[:len(last)-1]
		}
		lines[maxLines-1] = string(last) + "…"
	}

	return lines
}

// tokenize splits text into wrap units: each CJK/space-adjacent rune becomes its
// own token, while consecutive non-CJK non-space runes (Latin words) stay together.
func tokenize(text string) []string {
	var tokens []string
	var word strings.Builder
	flush := func() {
		if word.Len() > 0 {
			tokens = append(tokens, word.String())
			word.Reset()
		}
	}
	for _, r := range text {
		if isWordRune(r) {
			word.WriteRune(r)
			continue
		}
		flush()
		tokens = append(tokens, string(r))
	}
	flush()
	return tokens
}

// isWordRune reports whether r should stick to adjacent runes (kept as one word).
func isWordRune(r rune) bool {
	if unicode.IsSpace(r) {
		return false
	}
	if r > 0x2E7F && (unicode.Is(unicode.Han, r) ||
		unicode.Is(unicode.Hiragana, r) ||
		unicode.Is(unicode.Katakana, r) ||
		unicode.Is(unicode.Hangul, r)) {
		return false
	}
	return true
}

// loadFonts loads the bold/medium/regular weights, caching the result.
func (g *Generator) loadFonts() error {
	g.fontOnce.Do(func() {
		notoDirs := []string{
			"/usr/share/fonts/opentype/noto",
			"/usr/share/fonts/truetype/noto",
		}
		g.boldFont = loadWeight("Bold", notoDirs, g.fontPath)
		g.mediumFont = loadWeight("Medium", notoDirs, g.fontPath)
		g.regularFont = loadWeight("Regular", notoDirs, g.fontPath)

		if g.boldFont == nil || g.mediumFont == nil || g.regularFont == nil {
			g.fontErr = fmt.Errorf("failed to load required fonts")
		}
	})
	return g.fontErr
}

// loadWeight loads a Noto Sans CJK weight, falling back to the explicit override
// path and finally to the embedded Go font.
func loadWeight(weight string, notoDirs []string, overridePath string) *sfnt.Font {
	for _, dir := range notoDirs {
		path := fmt.Sprintf("%s/NotoSansCJK-%s.ttc", dir, weight)
		if f := loadFontFile(path); f != nil {
			slog.Info("Loaded OG font", slog.String("weight", weight), slog.String("path", path))
			return f
		}
	}
	if overridePath != "" {
		if f := loadFontFile(overridePath); f != nil {
			slog.Info("Loaded OG font (override)", slog.String("weight", weight), slog.String("path", overridePath))
			return f
		}
	}
	// Last resort: embedded Go font (no Japanese glyphs).
	if f, err := sfnt.Parse(goregular.TTF); err == nil {
		slog.Warn("Using embedded OG font (no Japanese support)", slog.String("weight", weight))
		return f
	}
	return nil
}

// loadFontFile parses a .ttf/.otf/.ttc file, preferring the JP variant for
// collections, and returns nil on any failure.
func loadFontFile(path string) *sfnt.Font {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	coll, err := opentype.ParseCollection(b)
	if err != nil {
		return nil
	}
	idx := 0
	for i := 0; i < coll.NumFonts(); i++ {
		f, err := coll.Font(i)
		if err != nil {
			continue
		}
		name, err := f.Name(nil, sfnt.NameIDFamily)
		if err == nil && strings.Contains(name, "JP") {
			idx = i
			break
		}
	}
	f, err := coll.Font(idx)
	if err != nil {
		return nil
	}
	return f
}

// sanitizeTitle removes control characters and limits length
func sanitizeTitle(title string) string {
	// Remove control characters except newline
	var cleaned strings.Builder
	for _, r := range title {
		if r == '\n' || r == '\r' || r == '\t' {
			cleaned.WriteRune(' ')
		} else if r >= 32 {
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
