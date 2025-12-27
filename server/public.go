package server

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/feeds"

	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/internal/markdown"
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

func RenderTopPage(c *gin.Context, queries *publicdb.Queries) {
	// Parse and execute the template
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// get query parameter 'page'
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			slog.Error("failed to parse page number", slog.String("page", pageStr), slog.Any("error", err))
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}

	entriesPerPage := 60
	offset := (page - 1) * entriesPerPage

	entries, err := queries.SearchEntries(c.Request.Context(), publicdb.SearchEntriesParams{
		Limit:  int32(entriesPerPage + 1),
		Offset: int32(offset),
	})
	if err != nil {
		slog.Error("failed to search entries", slog.Int("page", page), slog.Any("error", err))
		c.String(http.StatusInternalServerError, "Internal Server Error")
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
			c.String(http.StatusInternalServerError, "Internal Server Error")
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

	c.Status(http.StatusOK)
	err = tmpl.Execute(c.Writer, TopPageData{
		Page:    page,
		Prev:    page - 1,
		Next:    page + 1,
		HasPrev: page > 1,
		HasNext: hasNext,
		Entries: viewData,
	})
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}

func RenderEntryPage(c *gin.Context, queries *publicdb.Queries) {
	extractedPath := c.Param("filepath")
	// Strip leading slash from wildcard parameter
	extractedPath = strings.TrimPrefix(extractedPath, "/")

	md := markdown.NewMarkdown(c.Request.Context(), queries)

	slog.Info("rendering entry page", slog.String("path", extractedPath))
	entry, err := queries.GetEntryByPath(c.Request.Context(), extractedPath)
	if err != nil {
		slog.Error("failed to get entry by path", slog.String("path", extractedPath), slog.Any("error", err))
		c.Status(http.StatusNotFound)
		return
	}

	// Data to pass to the template
	var formattedDate string
	if entry.PublishedAt.Valid {
		formattedDate = formatDateTime(entry.PublishedAt.Time)
	} else {
		slog.Error("published_at is invalid", slog.String("path", entry.Path), slog.Any("published_at", entry.PublishedAt))
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	body, err := md.Render(entry.Body)
	if err != nil {
		slog.Error("failed to render markdown", slog.String("path", entry.Path), slog.Any("error", err))
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	relatedEntries, err := getRelatedEntries(c.Request.Context(), queries, entry)
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
	tmpl, err := template.ParseFiles("web/templates/entry.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = tmpl.Execute(c.Writer, data)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal Server Error")
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

func RenderFeed(c *gin.Context, queries *publicdb.Queries) {
	entries, err := queries.SearchEntries(c.Request.Context(), publicdb.SearchEntriesParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		slog.Error("failed to search entries for feed", slog.Any("error", err))
		c.String(http.StatusInternalServerError, "Internal Server Error")
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
	md := markdown.NewMarkdown(c.Request.Context(), queries)
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
		c.String(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	c.String(http.StatusOK, rss)
}

func RenderStaticMainCss(c *gin.Context) {
	// Read main.css from filesystem
	file, err := os.ReadFile("web/static/main.css")
	if err != nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	c.Data(http.StatusOK, "text/css", file)
}
