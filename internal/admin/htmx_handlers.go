package admin

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

type HtmxEntriesData struct {
	Entries    []EntryCard
	HasMore    bool
	LastCursor string
	Query      string
}

type EntryCard struct {
	Path        string
	Title       string
	BodyPreview string
	Visibility  string
	ImageUrl    string
}

func simplifyMarkdown(text string) string {
	// Remove newlines
	text = strings.ReplaceAll(text, "\n", " ")

	// Remove markdown links [text](url) -> text
	re := regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	text = re.ReplaceAllString(text, "$1")

	// Remove wiki links [[text]] -> text
	text = strings.ReplaceAll(text, "[[", "")
	text = strings.ReplaceAll(text, "]]", "")

	// Remove inline code
	re = regexp.MustCompile("`.*?`")
	text = re.ReplaceAllString(text, "")

	// Remove headers
	text = strings.ReplaceAll(text, "#", "")

	// Remove URLs
	re = regexp.MustCompile(`https?://\S+`)
	text = re.ReplaceAllString(text, "")

	// Normalize whitespace
	re = regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")

	text = strings.TrimSpace(text)

	// Truncate to 100 characters
	runes := []rune(text)
	if len(runes) > 100 {
		return string(runes[:100]) + "..."
	}
	return text
}

func RenderHtmxEntriesPage(queries *admindb.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		query := r.URL.Query().Get("q")
		lastEditedAtStr := r.URL.Query().Get("last_last_edited_at")

		var lastEditedAt sql.NullTime
		if lastEditedAtStr != "" {
			t, err := time.Parse(time.RFC3339, lastEditedAtStr)
			if err == nil {
				lastEditedAt = sql.NullTime{Time: t, Valid: true}
			}
		}

		// Perform search
		entries, err := queries.SearchEntriesAdmin(r.Context(), admindb.SearchEntriesAdminParams{
			Column1:      lastEditedAt,
			LastEditedAt: lastEditedAt,
			Column3:      query,
			Column4:      query,
			Limit:        100,
		})
		if err != nil {
			slog.Error("failed to search entries", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert to EntryCard
		var cards []EntryCard
		for _, e := range entries {
			cards = append(cards, EntryCard{
				Path:        e.Path,
				Title:       e.Title,
				BodyPreview: simplifyMarkdown(e.Body),
				Visibility:  string(e.Visibility),
				ImageUrl:    e.ImageUrl.String,
			})
		}

		// Determine if there are more entries
		hasMore := len(entries) >= 100
		var lastCursor string
		if hasMore && len(entries) > 0 {
			lastCursor = entries[len(entries)-1].LastEditedAt.Time.Format(time.RFC3339)
		}

		data := HtmxEntriesData{
			Entries:    cards,
			HasMore:    hasMore,
			LastCursor: lastCursor,
			Query:      query,
		}

		// Check if this is an htmx request
		isHtmxRequest := r.Header.Get("HX-Request") == "true"

		// Parse templates
		var tmpl *template.Template
		var err2 error

		if isHtmxRequest {
			// For htmx requests, return just the fragment
			tmpl, err2 = template.ParseFiles(
				"web/templates/admin/htmx_entry_cards.html",
			)
		} else {
			// For regular requests, return the full page
			tmpl, err2 = template.ParseFiles(
				"web/templates/admin/htmx_entries.html",
				"web/templates/admin/htmx_entry_cards.html",
			)
		}

		if err2 != nil {
			slog.Error("failed to parse template", slog.Any("error", err2))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if isHtmxRequest {
			err2 = tmpl.ExecuteTemplate(w, "entry-cards", data)
		} else {
			err2 = tmpl.Execute(w, data)
		}

		if err2 != nil {
			slog.Error("failed to execute template", slog.Any("error", err2))
		}
	}
}
