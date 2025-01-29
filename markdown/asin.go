package markdown

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tokuhirom/blog3/db/public/publicdb"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
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
		return 0, err
	}
	w.WriteString("<div style='display: flex;' class='asin'>")
	w.WriteString("<p>")
	w.WriteString("<a href=\"")
	w.WriteString(asin.Link)
	w.WriteString("\">")
	w.WriteString("<img src=\"")
	w.WriteString(asin.ImageMediumUrl.String)
	w.WriteString("\" style='max-width: 100px;max-height: 100px;border-radius: 4px;'>")
	w.WriteString("</a>")
	w.WriteString("</p>")
	w.WriteString("<p>")
	w.WriteString("<a href=\"")
	w.WriteString(asin.Link)
	w.WriteString("\">")
	w.WriteString(asin.Title.String)
	w.WriteString("</a>")
	w.WriteString("</p>")
	w.WriteString("</div>")

	return ast.WalkSkipChildren, nil
}
