package admin

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/internal"

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
