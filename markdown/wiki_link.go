package markdown

import (
	"bytes"
	"context"
	"fmt"
	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type WikiLink struct {
	Context context.Context
}

func (a WikiLink) Extend(markdown goldmark.Markdown) {
	markdown.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(&WikiParser{
				Context: a.Context,
			}, 100),
		),
	)
	markdown.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&WikiRenderer{
				Context: a.Context,
			}, 199),
		),
	)
}

type WikiParser struct {
	Context context.Context
}

func (p *WikiParser) Trigger() []byte {
	return []byte{'['}
}

var (
	wikiOpen  = []byte("[[")
	wikiClose = []byte("]]")
)

type WikiNode struct {
	ast.BaseInline

	Target []byte
	Embed  bool
}

var WikiKind = ast.NewNodeKind("WikiLink")

func (a *WikiNode) Kind() ast.NodeKind {
	return WikiKind
}

func (a *WikiNode) Dump(src []byte, level int) {
	ast.DumpHelper(a, src, level, map[string]string{
		"Target": string(a.Target),
	}, nil)
}

// [[Link]]
func (p *WikiParser) Parse(_ ast.Node, block text.Reader, _ parser.Context) ast.Node {
	line, seg := block.PeekLine()
	stop := bytes.Index(line, wikiClose)
	if stop < 0 {
		return nil // must close on the same line
	}

	var embed bool

	switch {
	case bytes.HasPrefix(line, wikiOpen):
		seg = text.NewSegment(seg.Start+len(wikiOpen), seg.Start+stop)
	default:
		println(string(line))
		return nil
	}

	n := &WikiNode{
		Target: block.Value(seg),
		Embed:  embed,
	}
	if len(n.Target) == 0 || seg.Len() == 0 {
		return nil // target and label must not be empty
	}

	block.Advance(stop + 2) // "]]".length == 2
	return n
}

type WikiRenderer struct {
	Resolver any
	Queries  *publicdb.Queries
	Context  context.Context
}

func (r *WikiRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(WikiKind, r.Render)
}

func (r *WikiRenderer) Render(writer util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, ok := node.(*WikiNode)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node %T, expected *WikiNode", node)
	}

	if entering {
		_, err := writer.WriteString(string(n.Target))
		if err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkSkipChildren, nil
	}

	return ast.WalkContinue, nil
}
