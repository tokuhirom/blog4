package markdown

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"github.com/tokuhirom/blog4/db/public/publicdb"
)

type AsinLink struct {
	Context context.Context
	Queries *publicdb.Queries
}

func (a AsinLink) Extend(markdown goldmark.Markdown) {
	markdown.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&AsinParser{
				Context: a.Context,
			}, 100),
		),
	)
	markdown.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&AsinRenderer{
				Context: a.Context,
				Queries: a.Queries,
			}, 199),
		),
	)
}

type AsinParser struct {
	Context context.Context
}

func (p *AsinParser) Trigger() []byte {
	return []byte{'['}
}

var (
	asinOpen = []byte("[asin:")
	asnClose = []byte(":detail]")
)

type AsinNode struct {
	ast.BaseInline

	Target []byte
	Embed  bool
}

var AsinKind = ast.NewNodeKind("AsinLink")

func (a *AsinNode) Kind() ast.NodeKind {
	return AsinKind
}

func (a *AsinNode) Dump(src []byte, level int) {
	ast.DumpHelper(a, src, level, map[string]string{
		"Target": string(a.Target),
	}, nil)
}

// [asin:B0BC73K2BW:detail]
func (p *AsinParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	stop := bytes.Index(line, asnClose)
	if stop < 0 {
		return nil // must close on the same line
	}

	var embed bool

	switch {
	case bytes.HasPrefix(line, asinOpen):
		seg = text.NewSegment(seg.Start+len(asinOpen), seg.Start+stop)
	default:
		return nil
	}

	n := &AsinNode{
		Target: block.Value(seg),
		Embed:  embed,
	}
	if len(n.Target) == 0 || seg.Len() == 0 {
		return nil // target and label must not be empty
	}

	block.Advance(stop + 8)
	return n
}

type AsinRenderer struct {
	Resolver any
	Queries  *publicdb.Queries
	Context  context.Context
}

func (r *AsinRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(AsinKind, r.Render)
}

func (r *AsinRenderer) Render(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, ok := node.(*AsinNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *AsinNode", node)
	}

	if entering {
		return r.enter(writer, n, source)
	}

	return ast.WalkContinue, nil
}
func (r *AsinRenderer) enter(w util.BufWriter, n *AsinNode, src []byte) (ast.WalkStatus, error) {
	asin, err := r.Queries.GetAsin(r.Context, string(n.Target))
	if err != nil {
		// If ASIN not found, render a placeholder or the raw ASIN link
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("ASIN not found in database, rendering as plain text", slog.String("asin", string(n.Target)))
		} else {
			slog.Error("Failed to query ASIN", slog.String("asin", string(n.Target)), slog.Any("error", err))
		}

		// Build fallback text
		var buf bytes.Buffer
		buf.WriteString("[asin:")
		buf.Write(n.Target)
		buf.WriteString(":detail]")

		// Write the complete string at once
		if _, err := w.Write(buf.Bytes()); err != nil {
			return 0, fmt.Errorf("failed to write ASIN fallback text: %w", err)
		}

		return ast.WalkSkipChildren, nil
	}

	// Build HTML string
	var buf bytes.Buffer
	buf.WriteString("<div style='display: flex;' class='asin'>")
	buf.WriteString("<p>")
	buf.WriteString("<a href=\"")
	buf.WriteString(asin.Link)
	buf.WriteString("\">")
	buf.WriteString("<img src=\"")
	buf.WriteString(asin.ImageMediumUrl.String)
	buf.WriteString("\" style='max-width: 100px;max-height: 100px;border-radius: 4px;'>")
	buf.WriteString("</a>")
	buf.WriteString("</p>")
	buf.WriteString("<p>")
	buf.WriteString("<a href=\"")
	buf.WriteString(asin.Link)
	buf.WriteString("\">")
	buf.WriteString(asin.Title.String)
	buf.WriteString("</a>")
	buf.WriteString("</p>")
	buf.WriteString("</div>")

	// Write the complete string at once
	if _, err := w.Write(buf.Bytes()); err != nil {
		return 0, fmt.Errorf("failed to write ASIN HTML: %w", err)
	}

	return ast.WalkSkipChildren, nil
}
