package utils

import (
	"bytes"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"html/template"
)

type Markdown struct {
	md goldmark.Markdown
}

func NewMarkdown() *Markdown {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,     // Enable GitHub Flavored Markdown
			extension.Linkify, // Enable auto-linking
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
