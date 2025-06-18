package markdown

import (
	"bytes"
	"context"
	"testing"

	"github.com/yuin/goldmark/text"
)

func TestWikiParser_Parse(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantTarget string
		wantNil    bool
	}{
		{
			name:       "valid wiki link",
			input:      "[[HomePage]]",
			wantTarget: "HomePage",
			wantNil:    false,
		},
		{
			name:       "wiki link with spaces",
			input:      "[[My Page Title]]",
			wantTarget: "My Page Title",
			wantNil:    false,
		},
		{
			name:       "wiki link with text after",
			input:      "[[PageName]] some text",
			wantTarget: "PageName",
			wantNil:    false,
		},
		{
			name:       "invalid - not closed on same line",
			input:      "[[PageName\n]]",
			wantTarget: "",
			wantNil:    true,
		},
		{
			name:       "invalid - missing closing brackets",
			input:      "[[PageName",
			wantTarget: "",
			wantNil:    true,
		},
		{
			name:       "invalid - empty target",
			input:      "[[]]",
			wantTarget: "",
			wantNil:    true,
		},
		{
			name:       "invalid - single bracket",
			input:      "[PageName]",
			wantTarget: "",
			wantNil:    true,
		},
		{
			name:       "invalid - not starting with [[",
			input:      "PageName]]",
			wantTarget: "",
			wantNil:    true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &WikiParser{
				Context: context.Background(),
			}
			
			reader := text.NewReader([]byte(tt.input))
			node := parser.Parse(nil, reader, nil)
			
			if tt.wantNil {
				if node != nil {
					t.Errorf("Parse() = %v, want nil", node)
				}
				return
			}
			
			if node == nil {
				t.Errorf("Parse() = nil, want WikiNode")
				return
			}
			
			wikiNode, ok := node.(*WikiNode)
			if !ok {
				t.Errorf("Parse() returned %T, want *WikiNode", node)
				return
			}
			
			if string(wikiNode.Target) != tt.wantTarget {
				t.Errorf("Parse() Target = %s, want %s", string(wikiNode.Target), tt.wantTarget)
			}
		})
	}
}

func TestWikiParser_Trigger(t *testing.T) {
	parser := &WikiParser{}
	trigger := parser.Trigger()
	if !bytes.Equal(trigger, []byte{'['}) {
		t.Errorf("Trigger() = %v, want %v", trigger, []byte{'['})
	}
}

func TestWikiNode_Kind(t *testing.T) {
	node := &WikiNode{}
	if node.Kind() != WikiKind {
		t.Errorf("Kind() = %v, want %v", node.Kind(), WikiKind)
	}
}

// Note: WikiRenderer.Render is difficult to test in isolation because it requires
// a proper util.BufWriter interface implementation from goldmark.
// This is documented in ISSUES.md for future refactoring.
// Integration tests would be more appropriate for testing the rendering functionality.