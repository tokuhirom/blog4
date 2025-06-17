package markdown

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/tokuhirom/blog4/db/public/publicdb"
)

type Markdown struct {
	md goldmark.Markdown
}

type WikiLinkResolver struct {
}

func NewMarkdown(ctx context.Context, queries *publicdb.Queries) *Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,     // Enable GitHub Flavored Markdown
			extension.Linkify, // Enable auto-linking
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
			),
			&AsinLink{
				Context: ctx,
				Queries: queries,
			},
			&WikiLink{
				Context: ctx,
			},
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),  // Render as XHTML
			html.WithUnsafe(), // Allow unsafe HTML (e.g., raw HTML tags)
		),
	)

	return &Markdown{
		md,
	}
}

func (m *Markdown) Render(input string) (template.HTML, error) {
	var buf bytes.Buffer
	if err := m.md.Convert([]byte(input), &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}
	return template.HTML(buf.String()), nil
}
