package public

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchTemplate(t *testing.T) {
	b, err := os.ReadFile("../../public/templates/search.html")
	assert.NoError(t, err)
	tmpl := string(b)

	// 検索フォームの要素は残す (JS が submit を制御する)。
	assert.Contains(t, tmpl, `<form class="search-form" action="/search" method="GET">`)
	assert.Contains(t, tmpl, `name="q"`)
	// クライアント検索 (search.js) が描画先とスクリプトを必要とする。
	assert.Contains(t, tmpl, `id="search-results"`)
	assert.Contains(t, tmpl, `/static/search.js`)
}
