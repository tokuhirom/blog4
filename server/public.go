package server

import (
	"embed"
	"github.com/tokuhirom/blog3/db/mariadb"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/feeds"
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

	md := utils.NewMarkdown(r.Context(), queries)

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

func RenderFeed(writer http.ResponseWriter, request *http.Request, queries *mariadb.Queries) {
	entries, err := queries.SearchEntries(request.Context(), 0)
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
	md := utils.NewMarkdown(request.Context(), queries)
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
	} else {
		http.NotFound(writer, request)
	}
}
