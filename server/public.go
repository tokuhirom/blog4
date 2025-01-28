package server

import (
	"bytes"
	"context"
	"embed"
	"github.com/tokuhirom/blog3/db/mariadb"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gorilla/feeds"
	"github.com/tokuhirom/blog3/markdown"
)

//go:embed templates/*
var templateFS embed.FS

//go:embed static/*
var staticFS embed.FS

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
	reURL := regexp.MustCompile(`https?://[^\s]+`)
	body = reURL.ReplaceAllString(body, "")

	// Replace [[foobar]] with foobar
	reBrackets := regexp.MustCompile(`\[\[(.*?)\]\]`)
	body = reBrackets.ReplaceAllString(body, "$1")

	// Trim to the specified length without cutting multibyte characters
	if utf8.RuneCountInString(body) > length {
		runes := []rune(body)
		body = string(runes[:length])
	}

	return body
}

func RenderTopPage(w http.ResponseWriter, r *http.Request, queries *mariadb.Queries) {
	// Parse and execute the template
	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// get query parameter 'page'
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			slog.Info("failed to parse page number: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	entriesPerPage := 60
	offset := (page - 1) * entriesPerPage

	entries, err := queries.SearchEntries(r.Context(), mariadb.SearchEntriesParams{
		Limit:  int32(entriesPerPage + 1),
		Offset: int32(offset),
	})
	if err != nil {
		slog.Info("failed to search entries: %v", err)
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
			log.Printf("published_at is invalid: path=%s, published_at=%v", entry.Path, entry.PublishedAt)
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

func RenderEntryPage(w http.ResponseWriter, r *http.Request, queries *mariadb.Queries) {
	extractedPath := strings.TrimPrefix(r.URL.Path, "/entry/")

	md := markdown.NewMarkdown(r.Context(), queries)

	log.Printf("path: %s", extractedPath)
	entry, err := queries.GetEntryByPath(r.Context(), extractedPath)
	if err != nil {
		slog.Info("failed to get entry by path", err)
		http.NotFound(w, r)
		return
	}

	// Data to pass to the template
	var formattedDate string
	if entry.PublishedAt.Valid {
		formattedDate = entry.PublishedAt.Time.Format("2006-01-01(Mon) 15:04")
	} else {
		log.Printf("published_at is invalid: path=%s, published_at=%v", entry.Path, entry.PublishedAt)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	body, err := md.Render(entry.Body)
	if err != nil {
		slog.Info("failed to render markdown", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	relatedEntries, err := getRelatedEntries(r.Context(), queries, entry)
	if err != nil {
		slog.Info("failed to get related entries", err)
	}

	data := struct {
		Title             string
		Body              template.HTML
		PublishedAt       string
		HasRelatedEntries bool
		RelatedEntries    []mariadb.Entry
	}{
		Title:             entry.Title,
		Body:              body,
		PublishedAt:       formattedDate,
		HasRelatedEntries: len(relatedEntries) > 0,
		RelatedEntries:    relatedEntries,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFS(templateFS, "templates/entry.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func getRelatedEntries(context context.Context, queries *mariadb.Queries, entry mariadb.Entry) ([]mariadb.Entry, error) {
	// 現在表示しているエントリがリンクしているページ
	entries1, err := queries.GetRelatedEntries1(context, entry.Path)
	if err != nil {
		return []mariadb.Entry{}, err
	}
	// 現在表示しているページにリンクしているページ
	entries2, err := queries.GetRelatedEntries2(context, entry.Title)
	if err != nil {
		return []mariadb.Entry{}, err
	}
	entries3, err := queries.GetRelatedEntries3(context, entry.Title)
	if err != nil {
		return []mariadb.Entry{}, err
	}

	// Use a map to track unique paths
	uniqueEntriesMap := make(map[string]mariadb.Entry)

	// Helper function to add entries to the map
	addUniqueEntries := func(entries []mariadb.Entry) {
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
	uniqueEntries := make([]mariadb.Entry, 0, len(uniqueEntriesMap))
	for _, entry := range uniqueEntriesMap {
		if entry.Visibility != "public" {
			// 保険的にvisibilityがpublicでないエントリは除外
			log.Fatalf("visibility is not public: %v", entry)
		}
		uniqueEntries = append(uniqueEntries, entry)
	}

	return uniqueEntries, nil
}

func RenderFeed(writer http.ResponseWriter, request *http.Request, queries *mariadb.Queries) {
	entries, err := queries.SearchEntries(request.Context(), mariadb.SearchEntriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		slog.Info("failed to search entries: %v", err)
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
			slog.Info("failed to render markdown", err, entry.Path)
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
		slog.Info("failed to generate RSS", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write([]byte(rss))
	if err != nil {
		slog.Info("failed to write response", err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func Handle(writer http.ResponseWriter, request *http.Request, queries *mariadb.Queries) {
	if request.URL.Path == "/" {
		RenderTopPage(writer, request, queries)
	} else if request.URL.Path == "/feed" {
		RenderFeed(writer, request, queries)
	} else if strings.HasPrefix(request.URL.Path, "/entry/") {
		RenderEntryPage(writer, request, queries)
	} else if request.URL.Path == "/static/main.css" {
		// if ./server/static/main.css is available, serve it.
		// hot reload.
		if _, err := os.Stat("server/static/main.css"); err == nil {
			http.ServeFile(writer, request, "server/static/main.css")
			return
		}

		writer.Header().Set("Content-Type", "text/css")
		file, err := staticFS.ReadFile("static/main.css")
		if err != nil {
			return
		}
		http.ServeContent(writer, request, "main.css", time.Time{}, bytes.NewReader(file))
	} else {
		http.NotFound(writer, request)
	}
}
