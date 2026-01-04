package admin

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
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

	"github.com/tokuhirom/blog4/internal"
	"github.com/tokuhirom/blog4/internal/sobs"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

// HtmxHandler handles htmx-based admin pages
type HtmxHandler struct {
	queries              *admindb.Queries
	sobsClient           *sobs.SobsClient
	adminUser            string
	adminPassword        string
	isSecure             bool
	s3AttachmentsBaseUrl string
}

// NewHtmxHandler creates a new HtmxHandler
func NewHtmxHandler(queries *admindb.Queries, sobsClient *sobs.SobsClient, adminUser, adminPassword string, isSecure bool, s3AttachmentsBaseUrl string) *HtmxHandler {
	return &HtmxHandler{
		queries:              queries,
		sobsClient:           sobsClient,
		adminUser:            adminUser,
		adminPassword:        adminPassword,
		isSecure:             isSecure,
		s3AttachmentsBaseUrl: s3AttachmentsBaseUrl,
	}
}

// getEntryPath extracts the entry path from query parameter
// For /entries/edit?path=getting-started -> returns "getting-started"
// For /entries/edit?path=2024/01/01/120000 -> returns "2024/01/01/120000"
func getEntryPath(c *gin.Context) string {
	return c.Query("path")
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
		tmpl, err = template.ParseFiles("admin/templates/htmx_entry_cards.html")
	} else {
		tmpl, err = template.ParseFiles(
			"admin/templates/layout.html",
			"admin/templates/htmx_entries.html",
			"admin/templates/htmx_entry_cards.html",
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
		"admin/templates/layout.html",
		"admin/templates/htmx_entry_edit.html",
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
	path := getEntryPath(c)

	// Delete existing entry image
	_, err := h.queries.DeleteEntryImageByPath(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to delete entry image", slog.String("path", path), slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Failed to delete entry image</div>`))
		return
	}

	// Process entry image in background
	go func() {
		ctx := c.Request.Context()
		service := internal.NewEntryImageService(h.queries)

		// Get the entry
		entryRow, err := h.queries.AdminGetEntryByPath(ctx, path)
		if err != nil {
			slog.Error("failed to get entry for image regeneration", slog.String("path", path), slog.Any("error", err))
			return
		}

		// Convert to Entry type
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

		// Process the entry
		err = service.ProcessEntry(ctx, entry)
		if err != nil {
			slog.Error("failed to process entry image", slog.String("path", path), slog.Any("error", err))
			return
		}

		slog.Info("successfully regenerated entry image", slog.String("path", path))
	}()

	c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-success">Image regeneration started!</div>`))
}

// UpdateEntryVisibility updates the entry visibility
func (h *HtmxHandler) UpdateEntryVisibility(c *gin.Context) {
	path := getEntryPath(c)
	visibility := c.PostForm("visibility")

	if visibility == "" {
		c.String(400, "Visibility is required")
		return
	}

	// Get current visibility and published_at
	entry, err := h.queries.GetEntryVisibility(c.Request.Context(), path)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("entry not found", slog.String("path", path))
			c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Entry not found</div>`))
			return
		}
		slog.Error("failed to get entry visibility", slog.String("path", path), slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Failed to get entry visibility</div>`))
		return
	}

	// Update visibility
	err = h.queries.UpdateVisibility(c.Request.Context(), admindb.UpdateVisibilityParams{
		Visibility: admindb.EntryVisibility(visibility),
		Path:       path,
	})
	if err != nil {
		slog.Error("failed to update visibility", slog.String("path", path), slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Failed to update visibility</div>`))
		return
	}

	// If changing from private to public and published_at is null, set it to now
	if entry.Visibility == "private" && visibility == "public" && !entry.PublishedAt.Valid {
		if err := h.queries.UpdatePublishedAt(c.Request.Context(), path); err != nil {
			slog.Error("failed to update published_at", slog.String("path", path), slog.Any("error", err))
			c.Data(200, "text/html; charset=utf-8", []byte(`<div class="feedback-error">Failed to update published_at</div>`))
			return
		}
		slog.Info("updated published_at for newly public entry", slog.String("path", path))
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

	// URL-encode the path to handle slashes correctly
	c.Header("HX-Redirect", "/admin/entries/edit?path="+url.QueryEscape(path))
	c.Status(200)
}

// HandleShareTarget handles Web Share Target API requests from Android
func (h *HtmxHandler) HandleShareTarget(c *gin.Context) {
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

// DeleteEntry deletes an entry
func (h *HtmxHandler) DeleteEntry(c *gin.Context) {
	path := getEntryPath(c)

	rows, err := h.queries.DeleteEntry(c.Request.Context(), path)
	if err != nil {
		slog.Error("failed to delete entry",
			slog.String("path", path),
			slog.Any("error", err))
		c.String(500, "Failed to delete entry")
		return
	}
	if rows == 0 {
		slog.Warn("entry to delete not found",
			slog.String("path", path))
		c.String(404, "Entry not found")
		return
	}

	c.Header("HX-Redirect", "/admin/entries/search")
	c.Status(200)
}

// RenderLoginPage displays the login page
func (h *HtmxHandler) RenderLoginPage(c *gin.Context) {
	tmpl, err := template.ParseFiles("admin/templates/htmx_login.html")
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

// HandleLogin processes the login form submission
func (h *HtmxHandler) HandleLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	rememberMe := c.PostForm("remember_me") == "true"

	// Validate credentials using constant-time comparison
	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(h.adminUser)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(password), []byte(h.adminPassword)) == 1

	if !usernameMatch || !passwordMatch {
		c.Data(200, "text/html; charset=utf-8", []byte(`Invalid username or password`))
		return
	}

	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		slog.Error("failed to generate session ID", slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`Failed to create session`))
		return
	}

	// Calculate session expiry based on remember_me option
	sessionTimeout := defaultSessionTimeout
	if rememberMe {
		sessionTimeout = extendedSessionTimeout
	}
	expires := time.Now().Add(sessionTimeout)

	// Create session in database
	err = h.queries.CreateSession(c.Request.Context(), admindb.CreateSessionParams{
		SessionID: sessionID,
		Username:  username,
		ExpiresAt: expires,
	})
	if err != nil {
		slog.Error("failed to create session", slog.Any("error", err))
		c.Data(200, "text/html; charset=utf-8", []byte(`Failed to create session`))
		return
	}

	// Set cookie
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

	// Redirect to entries page using HX-Redirect header
	c.Header("HX-Redirect", "/admin/entries/search")
	c.Status(200)
}

// UploadEntryImage handles image uploads from paste/drag-drop
func (h *HtmxHandler) UploadEntryImage(c *gin.Context) {
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
