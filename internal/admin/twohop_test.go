package admin

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/tokuhirom/blog4/internal/admin/openapi"
)

func Test_toOptNilString(t *testing.T) {
	tests := []struct {
		name string
		src  sql.NullString
		want openapi.OptNilString
	}{
		{
			name: "valid string",
			src:  sql.NullString{String: "test", Valid: true},
			want: openapi.OptNilString{Null: false, Value: "test"},
		},
		{
			name: "null string",
			src:  sql.NullString{String: "", Valid: false},
			want: openapi.OptNilString{Null: true},
		},
		{
			name: "empty valid string",
			src:  sql.NullString{String: "", Valid: true},
			want: openapi.OptNilString{Null: false, Value: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toOptNilString(tt.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toOptNilString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uniqueStrings(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "empty slice",
			input: []string{},
			want:  nil, // uniqueStrings returns nil for empty input, not empty slice
		},
		{
			name:  "no duplicates",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "with duplicates",
			input: []string{"a", "b", "a", "c", "b", "a"},
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "all same",
			input: []string{"x", "x", "x"},
			want:  []string{"x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := uniqueStrings(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uniqueStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uniqueEntries(t *testing.T) {
	tests := []struct {
		name  string
		input []*openapi.EntryWithImage
		want  int // number of unique entries expected
	}{
		{
			name:  "empty slice",
			input: []*openapi.EntryWithImage{},
			want:  0,
		},
		{
			name: "no duplicates",
			input: []*openapi.EntryWithImage{
				{Path: "/path1", Title: "Title 1"},
				{Path: "/path2", Title: "Title 2"},
				{Path: "/path3", Title: "Title 3"},
			},
			want: 3,
		},
		{
			name: "with duplicates",
			input: []*openapi.EntryWithImage{
				{Path: "/path1", Title: "Title 1"},
				{Path: "/path2", Title: "Title 2"},
				{Path: "/path1", Title: "Title 1 Modified"}, // duplicate path
				{Path: "/path3", Title: "Title 3"},
			},
			want: 3, // only 3 unique paths
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uniqueEntries(tt.input)
			if len(got) != tt.want {
				t.Errorf("uniqueEntries() returned %d entries, want %d", len(got), tt.want)
			}
			
			// Check that all entries have unique paths
			pathMap := make(map[string]bool)
			for _, entry := range got {
				if pathMap[entry.Path] {
					t.Errorf("uniqueEntries() returned duplicate path: %s", entry.Path)
				}
				pathMap[entry.Path] = true
			}
		})
	}
}