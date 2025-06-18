package markdown

import (
	"bytes"
	"context"
	"testing"

	"github.com/yuin/goldmark/text"
)

func TestAsinParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantASIN string
		wantNil  bool
	}{
		{
			name:     "valid ASIN link",
			input:    "[asin:B0BC73K2BW:detail]",
			wantASIN: "B0BC73K2BW",
			wantNil:  false,
		},
		{
			name:     "ASIN link with text after",
			input:    "[asin:B01M2BOZDL:detail] some text",
			wantASIN: "B01M2BOZDL",
			wantNil:  false,
		},
		{
			name:     "invalid - not closed on same line",
			input:    "[asin:B0BC73K2BW:detail\n]",
			wantASIN: "",
			wantNil:  true,
		},
		{
			name:     "invalid - missing closing bracket",
			input:    "[asin:B0BC73K2BW:detail",
			wantASIN: "",
			wantNil:  true,
		},
		{
			name:     "invalid - empty ASIN",
			input:    "[asin::detail]",
			wantASIN: "",
			wantNil:  true,
		},
		{
			name:     "invalid - not starting with [asin:",
			input:    "asin:B0BC73K2BW:detail]",
			wantASIN: "",
			wantNil:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &AsinParser{
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
				t.Errorf("Parse() = nil, want AsinNode")
				return
			}
			
			asinNode, ok := node.(*AsinNode)
			if !ok {
				t.Errorf("Parse() returned %T, want *AsinNode", node)
				return
			}
			
			if string(asinNode.Target) != tt.wantASIN {
				t.Errorf("Parse() ASIN = %s, want %s", string(asinNode.Target), tt.wantASIN)
			}
		})
	}
}

func TestAsinParser_Trigger(t *testing.T) {
	parser := &AsinParser{}
	trigger := parser.Trigger()
	if !bytes.Equal(trigger, []byte{'['}) {
		t.Errorf("Trigger() = %v, want %v", trigger, []byte{'['})
	}
}

func TestAsinNode_Kind(t *testing.T) {
	node := &AsinNode{}
	if node.Kind() != AsinKind {
		t.Errorf("Kind() = %v, want %v", node.Kind(), AsinKind)
	}
}