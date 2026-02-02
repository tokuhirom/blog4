package admin

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/internal"
	"github.com/tokuhirom/blog4/internal/markdown"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

// APIUpdateTitleRequest is the JSON request body for updating entry title
type APIUpdateTitleRequest struct {
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
}

// APIUpdateBodyRequest is the JSON request body for updating entry body
type APIUpdateBodyRequest struct {
	Body      string `json:"body"`
	UpdatedAt string `json:"updated_at"`
}

// APIUpdateVisibilityRequest is the JSON request body for updating entry visibility
type APIUpdateVisibilityRequest struct {
	Visibility string `json:"visibility"`
}

// APIResponse is the standard JSON response
type APIResponse struct {
	OK        bool   `json:"ok"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
	Redirect  string `json:"redirect,omitempty"`
}

// APIUpdateTitle updates the entry title and returns JSON
func (h *HtmxHandler) APIUpdateTitle(c *gin.Context) {
	path := getEntryPath(c)

	var req APIUpdateTitleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid request body"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Title cannot be empty"})
		return
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, req.UpdatedAt)
	if err != nil {
		slog.Error("failed to parse updated_at", slog.String("updated_at", req.UpdatedAt), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid updated_at"})
		return
	}

	rows, err := h.queries.UpdateEntryTitle(c.Request.Context(), admindb.UpdateEntryTitleParams{
		Title:     req.Title,
		Path:      path,
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
	})
	if err != nil {
		slog.Error("failed to update title", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to update title"})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusConflict, APIResponse{Error: "他のタブで更新されています。ページをリロードしてください。"})
		return
	}

	entry, err := h.queries.AdminGetEntryByPath(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to get entry after update", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusOK, APIResponse{OK: true, Message: "Title updated!"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		OK:        true,
		UpdatedAt: entry.UpdatedAt.Time.Format(time.RFC3339Nano),
		Message:   "Title updated!",
	})
}

// APIUpdateBody updates the entry body and returns JSON
func (h *HtmxHandler) APIUpdateBody(c *gin.Context) {
	path := getEntryPath(c)

	var req APIUpdateBodyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid request body"})
		return
	}

	if req.Body == "" {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Body cannot be empty"})
		return
	}

	updatedAt, err := time.Parse(time.RFC3339Nano, req.UpdatedAt)
	if err != nil {
		slog.Error("failed to parse updated_at", slog.String("updated_at", req.UpdatedAt), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid updated_at"})
		return
	}

	rows, err := h.queries.UpdateEntryBody(c.Request.Context(), admindb.UpdateEntryBodyParams{
		Body:      req.Body,
		Path:      path,
		UpdatedAt: sql.NullTime{Time: updatedAt, Valid: true},
	})
	if err != nil {
		slog.Error("failed to update body", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to update body"})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusConflict, APIResponse{Error: "他のタブで更新されています。ページをリロードしてください。"})
		return
	}

	entry, err := h.queries.AdminGetEntryByPath(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to get entry after update", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusOK, APIResponse{OK: true, Message: "Body updated!"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		OK:        true,
		UpdatedAt: entry.UpdatedAt.Time.Format(time.RFC3339Nano),
		Message:   "Body updated!",
	})
}

// APIUpdateVisibility updates the entry visibility and returns JSON
func (h *HtmxHandler) APIUpdateVisibility(c *gin.Context) {
	path := getEntryPath(c)

	var req APIUpdateVisibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid request body"})
		return
	}

	if req.Visibility == "" {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Visibility is required"})
		return
	}

	entry, err := h.queries.GetEntryVisibility(c.Request.Context(), path)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, APIResponse{Error: "Entry not found"})
			return
		}
		slog.Error("failed to get entry visibility", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to get entry visibility"})
		return
	}

	err = h.queries.UpdateVisibility(c.Request.Context(), admindb.UpdateVisibilityParams{
		Visibility: admindb.EntryVisibility(req.Visibility),
		Path:       path,
	})
	if err != nil {
		slog.Error("failed to update visibility", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to update visibility"})
		return
	}

	if entry.Visibility == "private" && req.Visibility == "public" {
		if !entry.PublishedAt.Valid {
			if err := h.queries.UpdatePublishedAt(c.Request.Context(), path); err != nil {
				slog.Error("failed to update published_at", slog.String("path", path), slog.Any("error", err))
				c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to update published_at"})
				return
			}
			slog.Info("updated published_at for newly public entry", slog.String("path", path))
		}

		if h.ogImageService != nil {
			go func() {
				ctx := context.Background()
				if err := h.ogImageService.EnsureOGImage(ctx, path); err != nil {
					slog.Error("failed to ensure OG image",
						slog.String("path", path), slog.Any("error", err))
				}
			}()
		}
	}

	c.JSON(http.StatusOK, APIResponse{
		OK:      true,
		Message: "Visibility updated to " + req.Visibility,
	})
}

// APIDeleteEntry deletes an entry and returns JSON
func (h *HtmxHandler) APIDeleteEntry(c *gin.Context) {
	path := getEntryPath(c)

	rows, err := h.queries.DeleteEntry(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to delete entry", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to delete entry"})
		return
	}
	if rows == 0 {
		c.JSON(http.StatusNotFound, APIResponse{Error: "Entry not found"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		OK:       true,
		Redirect: "/admin/entries/search",
	})
}

// APIRegenerateEntryImage triggers image regeneration and returns JSON
func (h *HtmxHandler) APIRegenerateEntryImage(c *gin.Context) {
	path := getEntryPath(c)

	_, err := h.queries.DeleteEntryImageByPath(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to delete entry image", slog.String("path", path), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to delete entry image"})
		return
	}

	go func() {
		ctx := context.Background()
		service := internal.NewEntryImageService(h.queries)

		entryRow, err := h.queries.AdminGetEntryByPath(ctx, path)
		if err != nil {
			slog.Error("failed to get entry for image regeneration", slog.String("path", path), slog.Any("error", err))
			return
		}

		entry := admindb.Entry{
			Path:         entryRow.Path,
			Title:        entryRow.Title,
			Body:         entryRow.Body,
			Visibility:   entryRow.Visibility,
			Format:       entryRow.Format,
			PublishedAt:  entryRow.PublishedAt,
			LastEditedAt: entryRow.LastEditedAt,
			CreatedAt:    entryRow.CreatedAt,
			UpdatedAt:    entryRow.UpdatedAt,
		}

		err = service.ProcessEntry(ctx, entry)
		if err != nil {
			slog.Error("failed to process entry image", slog.String("path", path), slog.Any("error", err))
			return
		}

		slog.Info("successfully regenerated entry image", slog.String("path", path))
	}()

	c.JSON(http.StatusOK, APIResponse{
		OK:      true,
		Message: "Image regeneration started!",
	})
}

// APIPreviewMarkdownRequest is the JSON request body for markdown preview
type APIPreviewMarkdownRequest struct {
	Body string `json:"body"`
}

// APIPreviewMarkdownResponse is the JSON response for markdown preview
type APIPreviewMarkdownResponse struct {
	HTML string `json:"html"`
}

// APIPreviewMarkdown renders markdown and returns the HTML
func (h *HtmxHandler) APIPreviewMarkdown(c *gin.Context) {
	var req APIPreviewMarkdownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid request body"})
		return
	}

	md := markdown.NewPreviewMarkdown(c.Request.Context())
	html, err := md.Render(req.Body)
	if err != nil {
		slog.Error("failed to render markdown preview", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to render markdown"})
		return
	}

	c.JSON(http.StatusOK, APIPreviewMarkdownResponse{HTML: string(html)})
}

// APIEntryCard represents a single entry in the list API response
type APIEntryCard struct {
	Path        string `json:"path"`
	Title       string `json:"title"`
	BodyPreview string `json:"body_preview"`
	Visibility  string `json:"visibility"`
	ImageURL    string `json:"image_url"`
}

// APIListEntriesResponse is the JSON response for the entries list API
type APIListEntriesResponse struct {
	Entries    []APIEntryCard `json:"entries"`
	HasMore    bool           `json:"has_more"`
	LastCursor string         `json:"last_cursor"`
}

// APIListEntries returns entries as JSON for the Preact entry-list app
func (h *HtmxHandler) APIListEntries(c *gin.Context) {
	searchQuery := c.Query("q")
	lastEditedAtStr := c.Query("last_last_edited_at")

	var cards []APIEntryCard
	var hasMore bool
	var lastCursor string

	if searchQuery != "" {
		searchResults, err := h.queries.AdminFullTextSearchEntries(c.Request.Context(), admindb.AdminFullTextSearchEntriesParams{
			Column1: searchQuery,
			Column2: searchQuery,
			Limit:   100,
		})
		if err != nil {
			slog.Error("failed to search entries", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search entries"})
			return
		}

		for _, e := range searchResults {
			cards = append(cards, APIEntryCard{
				Path:        e.Path,
				Title:       e.Title,
				BodyPreview: simplifyMarkdown(e.Body),
				Visibility:  string(e.Visibility),
				ImageURL:    e.ImageUrl.String,
			})
		}
		hasMore = false
	} else {
		var lastEditedAt sql.NullTime
		if lastEditedAtStr != "" {
			if t, err := time.Parse(time.RFC3339, lastEditedAtStr); err == nil {
				lastEditedAt = sql.NullTime{Time: t, Valid: true}
			}
		}

		entries, err := h.queries.GetLatestEntries(c.Request.Context(), admindb.GetLatestEntriesParams{
			Column1:      lastEditedAt,
			LastEditedAt: lastEditedAt,
			Limit:        100,
		})
		if err != nil {
			slog.Error("failed to get latest entries", slog.Any("error", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get entries"})
			return
		}

		for _, e := range entries {
			cards = append(cards, APIEntryCard{
				Path:        e.Path,
				Title:       e.Title,
				BodyPreview: simplifyMarkdown(e.Body),
				Visibility:  string(e.Visibility),
				ImageURL:    e.ImageUrl.String,
			})
		}

		hasMore = len(entries) >= 100
		if hasMore && len(entries) > 0 {
			lastCursor = entries[len(entries)-1].LastEditedAt.Time.Format(time.RFC3339)
		}
	}

	c.JSON(http.StatusOK, APIListEntriesResponse{
		Entries:    cards,
		HasMore:    hasMore,
		LastCursor: lastCursor,
	})
}

// APICreateEntryRequest is the JSON request body for creating an entry
type APICreateEntryRequest struct {
	Title string `json:"title"`
}

// APICreateEntry creates a new entry and returns JSON with redirect URL
func (h *HtmxHandler) APICreateEntry(c *gin.Context) {
	var req APICreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid request body"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Title is required"})
		return
	}

	now := time.Now()
	path := now.Format("2006/01/02/150405")

	_, err := h.queries.CreateEmptyEntry(c.Request.Context(), admindb.CreateEmptyEntryParams{
		Path:  path,
		Title: req.Title,
	})
	if err != nil {
		slog.Error("failed to create entry", slog.String("title", req.Title), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to create entry"})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		OK:       true,
		Redirect: "/admin/entries/edit?path=" + url.QueryEscape(path),
	})
}

// APILoginRequest is the JSON request body for login
type APILoginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

// APILogin handles JSON-based login and returns JSON response
func (h *HtmxHandler) APILogin(c *gin.Context) {
	var req APILoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: "Invalid request body"})
		return
	}

	usernameMatch := subtle.ConstantTimeCompare([]byte(req.Username), []byte(h.adminUser)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(req.Password), []byte(h.adminPassword)) == 1

	if !usernameMatch || !passwordMatch {
		c.JSON(http.StatusUnauthorized, APIResponse{Error: "Invalid username or password"})
		return
	}

	sessionID, err := generateSessionID()
	if err != nil {
		slog.Error("failed to generate session ID", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to create session"})
		return
	}

	sessionTimeout := defaultSessionTimeout
	if req.RememberMe {
		sessionTimeout = extendedSessionTimeout
	}
	expires := time.Now().Add(sessionTimeout)

	err = h.queries.CreateSession(c.Request.Context(), admindb.CreateSessionParams{
		SessionID: sessionID,
		Username:  req.Username,
		ExpiresAt: expires,
	})
	if err != nil {
		slog.Error("failed to create session", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, APIResponse{Error: "Failed to create session"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/admin",
		HttpOnly: true,
		Secure:   h.isSecure,
		SameSite: http.SameSiteStrictMode,
		Expires:  expires,
		MaxAge:   int(time.Until(expires).Seconds()),
	})

	c.JSON(http.StatusOK, APIResponse{
		OK:       true,
		Redirect: "/admin/entries/search",
	})
}
