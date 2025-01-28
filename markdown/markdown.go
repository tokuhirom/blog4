package markdown

import (
	"bytes"
	"context"
	"github.com/tokuhirom/blog3/db/mariadb"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
)

type Markdown struct {
	md goldmark.Markdown
}

type WikiLinkResolver struct {
	ctx context.Context
	db  *mariadb.Queries
}

func NewMarkdown(ctx context.Context, queries *mariadb.Queries) *Markdown {
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
		return "", err
	}
	return template.HTML(buf.String()), nil
}
