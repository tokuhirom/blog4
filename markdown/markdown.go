package markdown

import (
	"bytes"
	"context"
	"github.com/tokuhirom/blog3/db/mariadb"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/wikilink"
	"html/template"
	"log"
)

type Markdown struct {
	md goldmark.Markdown
}

type WikiLinkResolver struct {
	ctx context.Context
	db  *mariadb.Queries
}

func (w *WikiLinkResolver) ResolveWikilink(n *wikilink.Node) ([]byte, error) {
	entry, err := w.db.GetEntryByTitle(w.ctx, string(n.Target))
	if err != nil {
		log.Printf("failed to get entry by title: %v", err)
		return []byte(""), nil
	} else {
		return []byte(entry.Path), nil
	}
}

func NewMarkdown(ctx context.Context, queries *mariadb.Queries) *Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,     // Enable GitHub Flavored Markdown
			extension.Linkify, // Enable auto-linking
			&AsinLink{
				Context: ctx,
				Queries: queries,
			},
			&wikilink.Extender{
				Resolver: &WikiLinkResolver{
					ctx: ctx,
					db:  queries,
				},
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
