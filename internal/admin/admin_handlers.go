package admin

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tokuhirom/blog4/internal/ogimage"
	"github.com/tokuhirom/blog4/internal/sobs"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

// AdminHandler handles admin pages and JSON APIs
type AdminHandler struct {
	queries              *admindb.Queries
	sobsClient           *sobs.SobsClient
	adminUser            string
	adminPassword        string
	isSecure             bool
	s3AttachmentsBaseUrl string
	ogImageService       *ogimage.Service
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(queries *admindb.Queries, sobsClient *sobs.SobsClient, adminUser, adminPassword string, isSecure bool, s3AttachmentsBaseUrl string, ogImageService *ogimage.Service) *AdminHandler {
	return &AdminHandler{
		queries:              queries,
		sobsClient:           sobsClient,
		adminUser:            adminUser,
		adminPassword:        adminPassword,
		isSecure:             isSecure,
		s3AttachmentsBaseUrl: s3AttachmentsBaseUrl,
		ogImageService:       ogImageService,
	}
}

// getEntryPath extracts the entry path from query parameter
// For /entries/edit?path=getting-started -> returns "getting-started"
// For /entries/edit?path=2024/01/01/120000 -> returns "2024/01/01/120000"
func getEntryPath(c *gin.Context) string {
	return c.Query("path")
}

type EntriesPageData struct {
	Entries    []EntryCard
	HasMore    bool
	LastCursor string
	InitJSON   template.JS
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
func (h *AdminHandler) RenderEntriesPage(c *gin.Context) {
	// Get query parameters
	searchQuery := c.Query("q")
	lastEditedAtStr := c.Query("last_last_edited_at")

	var cards []EntryCard
	var hasMore bool
	var lastCursor string

	if searchQuery != "" {
		// Use full-text search
		searchResults, err := h.queries.AdminFullTextSearchEntries(c.Request.Context(), admindb.AdminFullTextSearchEntriesParams{
			Column1: searchQuery,
			Column2: searchQuery,
			Limit:   100,
		})
		if err != nil {
			slog.Error("failed to search entries", slog.Any("error", err))
			c.String(500, "Internal Server Error")
			return
		}

		// Convert to EntryCard
		for _, e := range searchResults {
			cards = append(cards, EntryCard{
				Path:        e.Path,
				Title:       e.Title,
				BodyPreview: simplifyMarkdown(e.Body),
				Visibility:  string(e.Visibility),
				ImageUrl:    e.ImageUrl.String,
			})
		}

		// Search results don't support pagination for now
		hasMore = false
	} else {
		// Get latest entries
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
			c.String(500, "Internal Server Error")
			return
		}

		// Convert to EntryCard
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
		hasMore = len(entries) >= 100
		if hasMore && len(entries) > 0 {
			lastCursor = entries[len(entries)-1].LastEditedAt.Time.Format(time.RFC3339)
		}
	}

	// Build JSON data for Preact app
	apiCards := make([]APIEntryCard, 0, len(cards))
	for _, card := range cards {
		apiCards = append(apiCards, APIEntryCard{
			Path:        card.Path,
			Title:       card.Title,
			BodyPreview: card.BodyPreview,
			Visibility:  card.Visibility,
			ImageURL:    card.ImageUrl,
		})
	}
	initData := APIListEntriesResponse{
		Entries:    apiCards,
		HasMore:    hasMore,
		LastCursor: lastCursor,
	}
	jsonBytes, err := json.Marshal(initData)
	if err != nil {
		slog.Error("failed to marshal entries data", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}

	data := EntriesPageData{
		Entries:    cards,
		HasMore:    hasMore,
		LastCursor: lastCursor,
		InitJSON:   template.JS(jsonBytes),
	}

	tmpl, err := template.ParseFiles(
		"admin/templates/layout.html",
		"admin/templates/entries.html",
	)
	if err != nil {
		slog.Error("failed to parse template", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.ExecuteTemplate(c.Writer, "layout", data)
}

// EntryEditData holds data for the entry edit page
type EntryEditData struct {
	Path       string
	Title      string
	Body       string
	Visibility string
	UpdatedAt  string
	InitJSON   template.JS
}

// RenderEntryEditPage displays the entry edit page
func (h *AdminHandler) RenderEntryEditPage(c *gin.Context) {
	path := getEntryPath(c)

	entry, err := h.queries.AdminGetEntryByPath(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to get entry", slog.String("path", path), slog.Any("error", err))
		c.String(404, "Entry not found")
		return
	}

	// Build JSON data for Preact app
	initData := map[string]string{
		"path":       entry.Path,
		"title":      entry.Title,
		"body":       entry.Body,
		"visibility": string(entry.Visibility),
		"updated_at": entry.UpdatedAt.Time.Format(time.RFC3339Nano),
	}
	jsonBytes, err := json.Marshal(initData)
	if err != nil {
		slog.Error("failed to marshal entry data", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}

	data := EntryEditData{
		Path:       entry.Path,
		Title:      entry.Title,
		Body:       entry.Body,
		Visibility: string(entry.Visibility),
		UpdatedAt:  entry.UpdatedAt.Time.Format(time.RFC3339Nano),
		InitJSON:   template.JS(jsonBytes),
	}

	tmpl, err := template.ParseFiles(
		"admin/templates/layout.html",
		"admin/templates/entry_edit.html",
	)
	if err != nil {
		slog.Error("failed to parse template", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	_ = tmpl.ExecuteTemplate(c.Writer, "layout", data)
}

// HandleShareTarget handles Web Share Target API requests from Android
func (h *AdminHandler) HandleShareTarget(c *gin.Context) {
	ctx := c.Request.Context()
	now := time.Now()

	// Get shared content from POST form
	sharedTitle := c.PostForm("title")
	sharedText := c.PostForm("text")
	sharedURL := c.PostForm("url")

	// Generate unique title
	var title string
	if sharedTitle == "" {
		// No title provided, generate from timestamp
		title = "Shared " + now.Format("2006-01-02 15:04:05")
	} else {
		// Use shared title with timestamp suffix to ensure uniqueness
		title = sharedTitle + " - " + now.Format("15:04:05")
	}

	// Build entry body from shared content
	var bodyParts []string
	if sharedText != "" {
		bodyParts = append(bodyParts, sharedText)
	}
	if sharedURL != "" {
		bodyParts = append(bodyParts, "\n[Source]("+sharedURL+")")
	}
	body := strings.Join(bodyParts, "\n\n")

	// Generate path based on current time
	path := now.Format("2006/01/02/150405")

	// Create entry with body
	_, err := h.queries.CreateEntryWithBody(ctx, admindb.CreateEntryWithBodyParams{
		Path:  path,
		Title: title,
		Body:  body,
	})
	if err != nil {
		slog.Error("failed to create shared entry",
			slog.String("title", title),
			slog.String("path", path),
			slog.Any("error", err))

		// Render error template
		tmpl, tmplErr := template.ParseFiles("admin/templates/share_error.html")
		if tmplErr != nil {
			slog.Error("failed to parse share error template", slog.Any("error", tmplErr))
			c.String(500, "Failed to save shared content")
			return
		}

		c.Status(500)
		if err := tmpl.Execute(c.Writer, gin.H{
			"error":      "Failed to save shared content. Please try again.",
			"entriesUrl": "/admin/entries/search",
		}); err != nil {
			slog.Error("failed to execute share error template", slog.Any("error", err))
			c.String(500, "Failed to save shared content")
		}
		return
	}

	slog.Info("created shared entry",
		slog.String("title", title),
		slog.String("path", path),
		slog.String("sharedFrom", sharedURL))

	// Redirect to edit page (regular HTTP redirect, not HX-Redirect)
	// URL-encode the path to handle slashes correctly
	c.Redirect(http.StatusSeeOther, "/admin/entries/edit?path="+url.QueryEscape(path))
}

// RenderLoginPage displays the login page
func (h *AdminHandler) RenderLoginPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("admin/templates/login.html")
	if err != nil {
		slog.Error("failed to parse login template", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		slog.Error("failed to execute login template", slog.Any("error", err))
		c.String(500, "Internal Server Error")
		return
	}
}

// UploadEntryImage handles image uploads from paste/drag-drop
func (h *AdminHandler) UploadEntryImage(c *gin.Context) {
	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		slog.Error("failed to get uploaded file", slog.Any("error", err))
		c.JSON(400, gin.H{"error": "Invalid file"})
		return
	}

	// Validate MIME type
	contentType := file.Header.Get("Content-Type")
	if !isValidImageMimeType(contentType) {
		slog.Warn("invalid file type", slog.String("contentType", contentType))
		c.JSON(400, gin.H{"error": "Only image files are allowed"})
		return
	}

	// Validate file size (10MB limit)
	const maxUploadSize int64 = 10 * 1024 * 1024
	if file.Size > maxUploadSize {
		slog.Warn("file too large", slog.Int64("size", file.Size))
		c.JSON(400, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// Generate S3 key
	key, err := generateImageKey(file)
	if err != nil {
		slog.Error("failed to generate image key", slog.Any("error", err))
		c.JSON(500, gin.H{"error": "Failed to generate file name"})
		return
	}

	// Open file
	fileContent, err := file.Open()
	if err != nil {
		slog.Error("failed to open uploaded file", slog.Any("error", err))
		c.JSON(500, gin.H{"error": "Failed to read file"})
		return
	}
	defer func() {
		_ = fileContent.Close()
	}()

	// Upload to S3
	err = h.sobsClient.PutObjectToAttachmentBucket(
		c.Request.Context(),
		key,
		contentType,
		file.Size,
		fileContent,
	)
	if err != nil {
		slog.Error("failed to upload to S3", slog.String("key", key), slog.Any("error", err))
		c.JSON(500, gin.H{"error": "Upload failed"})
		return
	}

	// Generate URL using configured base URL
	url := fmt.Sprintf("%s/%s", strings.TrimSuffix(h.s3AttachmentsBaseUrl, "/"), key)

	slog.Info("image uploaded successfully", slog.String("key", key), slog.String("url", url))
	c.JSON(200, gin.H{"url": url})
}

// isValidImageMimeType checks if the MIME type is a valid image type
func isValidImageMimeType(contentType string) bool {
	allowedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/svg+xml",
	}

	for _, allowed := range allowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// generateImageKey generates a unique S3 key for the uploaded image
func generateImageKey(file *multipart.FileHeader) (string, error) {
	// Get file extension
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// Infer extension from MIME type
		ext = getExtensionFromMimeType(file.Header.Get("Content-Type"))
	}

	// Generate timestamp + random string for uniqueness
	now := time.Now()
	timeStr := strconv.FormatInt(now.UnixMilli(), 10)
	randomStr, err := generateRandomString(6)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf(
		"attachments/%04d/%02d/%02d/%s-%s%s",
		now.Year(),
		now.Month(),
		now.Day(),
		timeStr,
		randomStr,
		ext,
	)

	return key, nil
}

// getExtensionFromMimeType returns the file extension for a given MIME type
func getExtensionFromMimeType(mimeType string) string {
	mimeToExt := map[string]string{
		"image/jpeg":    ".jpg",
		"image/jpg":     ".jpg",
		"image/png":     ".png",
		"image/gif":     ".gif",
		"image/webp":    ".webp",
		"image/svg+xml": ".svg",
	}

	if ext, ok := mimeToExt[mimeType]; ok {
		return ext
	}
	return ".bin"
}

// generateRandomString generates a random alphanumeric string of given length
func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b), nil
}
