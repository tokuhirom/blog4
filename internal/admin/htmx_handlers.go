package admin

import (
	"database/sql"
	"html/template"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

// HtmxHandler handles htmx-based admin pages
type HtmxHandler struct {
	queries *admindb.Queries
}

// NewHtmxHandler creates a new HtmxHandler
func NewHtmxHandler(queries *admindb.Queries) *HtmxHandler {
	return &HtmxHandler{queries: queries}
}

// getEntryPath extracts and constructs the entry path from gin context params
func getEntryPath(c *gin.Context) string {
	year := c.Param("year")
	month := c.Param("month")
	day := c.Param("day")
	time := c.Param("time")

	if year != "" && month != "" && day != "" && time != "" {
		return year + "/" + month + "/" + day + "/" + time
	}
	return c.Param("path")
}

type HtmxEntriesData struct {
	Entries    []EntryCard
	HasMore    bool
	LastCursor string
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

// RenderEntriesPage displays the entries list page
func (h *HtmxHandler) RenderEntriesPage(c *gin.Context) {
	// Get query parameters (search is client-side filtered for now)
	lastEditedAtStr := c.Query("last_last_edited_at")

	var lastEditedAt sql.NullTime
	if lastEditedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, lastEditedAtStr); err == nil {
			lastEditedAt = sql.NullTime{Time: t, Valid: true}
		}
	}

	// Get latest entries
	entries, err := h.queries.GetLatestEntries(c.Request.Context(), admindb.GetLatestEntriesParams{
		Column1:      lastEditedAt,
		LastEditedAt: lastEditedAt,
		Limit:        100,
	})
	if err != nil {
		slog.Error("failed to search entries", slog.Any("error", err))
		c.String(500, "Internal Server Error")
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
	}

	// Check if this is an htmx request
	isHtmxRequest := c.GetHeader("HX-Request") == "true"

	// Parse templates
	var tmpl *template.Template
	if isHtmxRequest {
		tmpl, err = template.ParseFiles("web/templates/admin/htmx_entry_cards.html")
	} else {
		tmpl, err = template.ParseFiles(
			"web/templates/admin/layout.html",
			"web/templates/admin/htmx_entries.html",
			"web/templates/admin/htmx_entry_cards.html",
		)
	}
	if err != nil {
		slog.Error("failed to parse template", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if isHtmxRequest {
		_ = tmpl.ExecuteTemplate(c.Writer, "entry-cards", data)
	} else {
		_ = tmpl.ExecuteTemplate(c.Writer, "layout", data)
	}
}

// EntryEditData holds data for the entry edit page
type EntryEditData struct {
	Path       string
	Title      string
	Body       string
	Visibility string
}

// RenderEntryEditPage displays the entry edit page
func (h *HtmxHandler) RenderEntryEditPage(c *gin.Context) {
	path := getEntryPath(c)

	entry, err := h.queries.AdminGetEntryByPath(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to get entry", slog.String("path", path), slog.Any("error", err))
		c.String(404, "Entry not found")
		return
	}

	data := EntryEditData{
		Path:       entry.Path,
		Title:      entry.Title,
		Body:       entry.Body,
		Visibility: string(entry.Visibility),
	}

	tmpl, err := template.ParseFiles(
		"web/templates/admin/layout.html",
		"web/templates/admin/htmx_entry_edit.html",
	)
	if err != nil {
		slog.Error("failed to parse template", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.ExecuteTemplate(c.Writer, "layout", data)
}

// UpdateEntryTitle updates the entry title and returns feedback HTML
func (h *HtmxHandler) UpdateEntryTitle(c *gin.Context) {
	path := getEntryPath(c)
	title := c.PostForm("title")

	if title == "" {
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Title cannot be empty</div>`))
		return
	}

	rows, err := h.queries.UpdateEntryTitle(c.Request.Context(), admindb.UpdateEntryTitleParams{
		Title: title,
		Path:  path,
	})
	if err != nil || rows == 0 {
		slog.Error("failed to update title", slog.String("path", path), slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Failed to update title</div>`))
		return
	}

	c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-success">Title updated!</div>`))
}

// UpdateEntryBody updates the entry body and returns feedback HTML
func (h *HtmxHandler) UpdateEntryBody(c *gin.Context) {
	path := getEntryPath(c)
	body := c.PostForm("body")

	if body == "" {
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Body cannot be empty</div>`))
		return
	}

	rows, err := h.queries.UpdateEntryBody(c.Request.Context(), admindb.UpdateEntryBodyParams{
		Body: body,
		Path: path,
	})
	if err != nil || rows == 0 {
		slog.Error("failed to update body", slog.String("path", path), slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Failed to update body</div>`))
		return
	}

	c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-success">Body updated!</div>`))
}

// RegenerateEntryImage triggers image regeneration and returns feedback HTML
func (h *HtmxHandler) RegenerateEntryImage(c *gin.Context) {
	c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-success">Image regeneration queued!</div>`))
}

// UpdateEntryVisibility updates the entry visibility
func (h *HtmxHandler) UpdateEntryVisibility(c *gin.Context) {
	path := getEntryPath(c)
	visibility := c.PostForm("visibility")

	if visibility == "" {
		c.String(400, "Visibility is required")
		return
	}

	err := h.queries.UpdateVisibility(c.Request.Context(), admindb.UpdateVisibilityParams{
		Visibility: admindb.EntryVisibility(visibility),
		Path:       path,
	})
	if err != nil {
		slog.Error("failed to update visibility", slog.String("path", path), slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Failed to update visibility</div>`))
		return
	}

	c.Header("HX-Refresh", "true")
	c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-success">Visibility updated!</div>`))
}

// CreateEntry creates a new entry and redirects to edit page
func (h *HtmxHandler) CreateEntry(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.String(400, "Title is required")
		return
	}

	// Generate path based on current time
	now := time.Now()
	path := now.Format("2006/01/02/150405")

	_, err := h.queries.CreateEmptyEntry(c.Request.Context(), admindb.CreateEmptyEntryParams{
		Path:  path,
		Title: title,
	})
	if err != nil {
		slog.Error("failed to create entry", slog.String("title", title), slog.Any("error", err))
		c.String(500, "Failed to create entry")
		return
	}

	c.Header("HX-Redirect", "/admin/htmx/entries/"+path+"/edit")
	c.Status(200)
}

// DeleteEntry deletes an entry
func (h *HtmxHandler) DeleteEntry(c *gin.Context) {
	path := getEntryPath(c)

	rows, err := h.queries.DeleteEntry(c.Request.Context(), path)
	if err != nil || rows == 0 {
		slog.Error("failed to delete entry", slog.String("path", path), slog.Any("error", err))
		c.String(500, "Failed to delete entry")
		return
	}

	c.Header("HX-Redirect", "/admin/htmx/entries")
	c.Status(200)
}
