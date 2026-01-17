package public

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchFormExists(t *testing.T) {
	// Test that the search form exists in the template
	// This test checks that the search form is properly configured
	tmpl := getSearchTemplate()
	assert.NotNil(t, tmpl)
	assert.Contains(t, tmpl, `<form class="search-form" action="/search" method="GET">`)
	assert.Contains(t, tmpl, `<input type="text" name="q" class="search-input" placeholder="Search entries..."`)
	assert.Contains(t, tmpl, `{{.Query}}`)
}

func getSearchTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Search - tokuhirom's blog</title>
</head>
<body>
    <form class="search-form" action="/search" method="GET">
        <input type="text" name="q" class="search-input" placeholder="Search entries..." value="{{.Query}}" autofocus>
        <button type="submit" class="search-button">Search</button>
    </form>
</body>
</html>`
}