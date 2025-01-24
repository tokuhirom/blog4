package server

import (
	"embed"
	"github.com/tokuhirom/blog3/db/mariadb"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/tokuhirom/blog3/utils"
)

//go:embed templates/*
var templateFS embed.FS

type EntryViewData struct {
	Path        string
	Title       string
	PublishedAt string
}

func RenderTopPage(w http.ResponseWriter, r *http.Request, queries *mariadb.Queries) {
	// Parse and execute the template
	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	entries, err := queries.SearchEntries(r.Context(), 0)
	if err != nil {
		slog.Info("failed to search entries: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
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
		})
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, viewData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func RenderEntryPage(w http.ResponseWriter, r *http.Request, queries *mariadb.Queries) {
	extractedPath := strings.TrimPrefix(r.URL.Path, "/entry/")

	md := utils.NewMarkdown()

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

	data := struct {
		Title       string
		Body        template.HTML
		PublishedAt string
	}{
		Title:       entry.Title,
		Body:        body,
		PublishedAt: formattedDate,
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
