package server

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	"github.com/tokuhirom/blog4/db/public/publicdb"

	"github.com/gorilla/feeds"

	"github.com/tokuhirom/blog4/internal/markdown"
	"github.com/tokuhirom/blog4/web"
)

type TopPageData struct {
	Page    int
	Entries []EntryViewData
	HasPrev bool
	HasNext bool
	Prev    int
	Next    int
}

type EntryViewData struct {
	Path        string
	Title       string
	PublishedAt string
	TextPreview string
	ImageUrl    string
}

// summarizeEntry takes a string and returns a processed summary.
func summarizeEntry(body string, length int) string {
	// Remove URLs
	reURL := regexp.MustCompile(`https?://\S+`)
	body = reURL.ReplaceAllString(body, "")

	// Replace [[foobar]] with foobar
	reBrackets := regexp.MustCompile(`\[\[(.*?)]]`)
	body = reBrackets.ReplaceAllString(body, "$1")

	// Trim to the specified length without cutting multibyte characters
	if utf8.RuneCountInString(body) > length {
		runes := []rune(body)
		body = string(runes[:length])
	}

	return body
}

func RenderTopPage(w http.ResponseWriter, r *http.Request, queries *publicdb.Queries) {
	// Parse and execute the template
	tmpl, err := template.ParseFS(web.TemplateFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// get query parameter 'page'
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			slog.Error("failed to parse page number", slog.String("page", pageStr), slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	entriesPerPage := 60
	offset := (page - 1) * entriesPerPage

	entries, err := queries.SearchEntries(r.Context(), publicdb.SearchEntriesParams{
		Limit:  int32(entriesPerPage + 1),
		Offset: int32(offset),
	})
	if err != nil {
		slog.Error("failed to search entries", slog.Int("page", page), slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// remove last entry if there are more entries
	var hasNext = false
	if len(entries) > entriesPerPage {
		entries = entries[:entriesPerPage]
		hasNext = true
	}

	// Prepare data for the template
	var viewData []EntryViewData
	for _, entry := range entries {
		// Format the PublishedAt date
		var formattedDate string
		if entry.PublishedAt.Valid {
			formattedDate = entry.PublishedAt.Time.Format("2006-01-02(Mon)")
		} else {
			slog.Error("published_at is invalid", slog.String("path", entry.Path), slog.Any("published_at", entry.PublishedAt))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		viewData = append(viewData, EntryViewData{
			Path:        entry.Path,
			Title:       entry.Title,
			PublishedAt: formattedDate,
			TextPreview: summarizeEntry(entry.Body, 100),
			ImageUrl:    entry.ImageUrl.String,
		})
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, TopPageData{
		Page:    page,
		Prev:    page - 1,
		Next:    page + 1,
		HasPrev: page > 1,
		HasNext: hasNext,
		Entries: viewData,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RenderEntryPage(w http.ResponseWriter, r *http.Request, queries *publicdb.Queries) {
	extractedPath := chi.URLParam(r, "*")

	md := markdown.NewMarkdown(r.Context(), queries)

	slog.Info("rendering entry page", slog.String("path", extractedPath))
	entry, err := queries.GetEntryByPath(r.Context(), extractedPath)
	if err != nil {
		slog.Error("failed to get entry by path", slog.String("path", extractedPath), slog.Any("error", err))
		http.NotFound(w, r)
		return
	}

	// Data to pass to the template
	var formattedDate string
	if entry.PublishedAt.Valid {
		formattedDate = formatDateTime(entry.PublishedAt.Time)
	} else {
		slog.Error("published_at is invalid", slog.String("path", entry.Path), slog.Any("published_at", entry.PublishedAt))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	body, err := md.Render(entry.Body)
	if err != nil {
		slog.Error("failed to render markdown", slog.String("path", entry.Path), slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	relatedEntries, err := getRelatedEntries(r.Context(), queries, entry)
	if err != nil {
		slog.Error("failed to get related entries", slog.String("path", entry.Path), slog.Any("error", err))
	}

	data := struct {
		Title             string
		Body              template.HTML
		PublishedAt       string
		HasRelatedEntries bool
		RelatedEntries    []publicdb.Entry
	}{
		Title:             entry.Title,
		Body:              body,
		PublishedAt:       formattedDate,
		HasRelatedEntries: len(relatedEntries) > 0,
		RelatedEntries:    relatedEntries,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFS(web.TemplateFS, "templates/entry.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getRelatedEntries(context context.Context, queries *publicdb.Queries, entry publicdb.Entry) ([]publicdb.Entry, error) {
	// 現在表示しているエントリがリンクしているページ
	entries1, err := queries.GetRelatedEntries1(context, entry.Path)
	if err != nil {
		return []publicdb.Entry{}, fmt.Errorf("failed to get related entries1 for path %s: %w", entry.Path, err)
	}
	// 現在表示しているページにリンクしているページ
	entries2, err := queries.GetRelatedEntries2(context, entry.Title)
	if err != nil {
		return []publicdb.Entry{}, fmt.Errorf("failed to get related entries2 for title %s: %w", entry.Title, err)
	}
	entries3, err := queries.GetRelatedEntries3(context, entry.Title)
	if err != nil {
		return []publicdb.Entry{}, fmt.Errorf("failed to get related entries3 for title %s: %w", entry.Title, err)
	}

	// Use a map to track unique paths
	uniqueEntriesMap := make(map[string]publicdb.Entry)

	// Helper function to add entries to the map
	addUniqueEntries := func(entries []publicdb.Entry) {
		for _, e := range entries {
			if _, exists := uniqueEntriesMap[e.Path]; !exists {
				uniqueEntriesMap[e.Path] = e
			}
		}
	}

	// Add entries from each slice
	addUniqueEntries(entries1)
	addUniqueEntries(entries2)
	addUniqueEntries(entries3)

	// Convert map values to a slice
	uniqueEntries := make([]publicdb.Entry, 0, len(uniqueEntriesMap))
	for _, entry := range uniqueEntriesMap {
		if entry.Visibility != "public" {
			// 保険的にvisibilityがpublicでないエントリは除外
			slog.Error("unexpected non-public entry in related entries", slog.String("path", entry.Path), slog.String("visibility", string(entry.Visibility)))
			continue // Skip this entry instead of exiting
		}
		uniqueEntries = append(uniqueEntries, entry)
	}

	return uniqueEntries, nil
}

func RenderFeed(writer http.ResponseWriter, request *http.Request, queries *publicdb.Queries) {
	entries, err := queries.SearchEntries(request.Context(), publicdb.SearchEntriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		slog.Error("failed to search entries for feed", slog.Any("error", err))
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	now := time.Now()
	feed := &feeds.Feed{
		Title:       "tokuhirom's blog",
		Link:        &feeds.Link{Href: "https://blog.64p.org/"},
		Description: "tokuhirom's thoughts",
		Author:      &feeds.Author{Name: "Tokuhiro Matsuno", Email: "tokuhirom+blog-gmail.com"},
		Created:     now,
	}
	md := markdown.NewMarkdown(request.Context(), queries)
	for _, entry := range entries {
		render, err := md.Render(entry.Body)
		if err != nil {
			slog.Error("failed to render markdown for feed", slog.String("path", entry.Path), slog.Any("error", err))
			// skip this entry
			continue
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Title:       entry.Title,
			Link:        &feeds.Link{Href: "https://blog.64p.org/entry/" + entry.Path},
			Description: entry.Body,
			Content:     string(render),
			Created:     entry.PublishedAt.Time,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		slog.Error("failed to generate RSS", slog.Any("error", err))
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(rss))
	if err != nil {
		slog.Error("failed to write RSS response", slog.Any("error", err))
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func RenderStaticMainCss(writer http.ResponseWriter, request *http.Request) {
	// if ./web/static/main.css is available, serve it.
	// hot reload.
	if _, err := os.Stat("web/static/main.css"); err == nil {
		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		writer.Header().Set("Pragma", "no-cache")
		writer.Header().Set("Expires", "0")

		http.ServeFile(writer, request, "web/static/main.css")
		return
	}

	writer.Header().Set("Content-Type", "text/css")
	file, err := web.StaticFS.ReadFile("static/main.css")
	if err != nil {
		return
	}
	http.ServeContent(writer, request, "main.css", time.Time{}, bytes.NewReader(file))
}

func Router(queries *publicdb.Queries) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		RenderTopPage(w, r, queries)
	})
	r.Get("/feed", func(w http.ResponseWriter, r *http.Request) {
		RenderFeed(w, r, queries)
	})
	r.Get("/entry/*", func(w http.ResponseWriter, r *http.Request) {
		RenderEntryPage(w, r, queries)
	})
	r.Get("/static/main.css", RenderStaticMainCss)
	return r
}
