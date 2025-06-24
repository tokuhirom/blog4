package admin

import (
	"bytes"
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/admin/openapi"
	"github.com/tokuhirom/blog4/server/sobs"
)

type adminApiService struct {
	queries       *admindb.Queries
	db            *sql.DB
	hubUrls       []string
	paapiClient   *PAAPIClient
	S3Client      *sobs.SobsClient
	adminUser     string
	adminPassword string
	isSecure      bool
}

func (p *adminApiService) GetLatestEntries(ctx context.Context, params openapi.GetLatestEntriesParams) (openapi.GetLatestEntriesRes, error) {
	var lastEditedAt sql.NullTime
	if params.LastLastEditedAt.IsSet() {
		lastEditedAt = sql.NullTime{
			Time:  params.LastLastEditedAt.Value,
			Valid: true,
		}
	} else {
		lastEditedAt = sql.NullTime{
			Valid: false,
		}
	}
	slog.Info("GetLatestEntries", slog.Any("lastEditedAt", lastEditedAt))
	entries, err := p.queries.GetLatestEntries(ctx, admindb.GetLatestEntriesParams{
		Column1:      lastEditedAt,
		LastEditedAt: lastEditedAt,
		Limit:        100,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get latest entries: %w", err)
	}

	var result []openapi.GetLatestEntriesRow
	for _, entry := range entries {
		result = append(result, openapi.GetLatestEntriesRow{
			Path:         entry.Path,
			Title:        entry.Title,
			Body:         entry.Body,
			Visibility:   string(entry.Visibility),
			Format:       string(entry.Format),
			PublishedAt:  openapi.NewOptNilDateTime(entry.PublishedAt.Time),
			LastEditedAt: openapi.NewOptNilDateTime(entry.LastEditedAt.Time),
			CreatedAt:    openapi.NewOptNilDateTime(entry.CreatedAt.Time),
			UpdatedAt:    openapi.NewOptNilDateTime(entry.UpdatedAt.Time),
			ImageUrl:     openapi.NewOptNilString(entry.ImageUrl.String),
		})
	}
	resp := openapi.GetLatestEntriesOKApplicationJSON(result)
	return &resp, nil
}

func (p *adminApiService) GetEntryByDynamicPath(ctx context.Context, params openapi.GetEntryByDynamicPathParams) (openapi.GetEntryByDynamicPathRes, error) {
	entry, err := p.queries.AdminGetEntryByPath(ctx, params.Path)
	if err != nil {
		slog.Error("GetEntryByDynamicPath failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get entry by path %s: %w", params.Path, err)
	}

	return &openapi.GetLatestEntriesRow{
		Path:         entry.Path,
		Title:        entry.Title,
		Body:         entry.Body,
		Visibility:   string(entry.Visibility),
		Format:       string(entry.Format),
		PublishedAt:  openapi.NewOptNilDateTime(entry.PublishedAt.Time),
		LastEditedAt: openapi.NewOptNilDateTime(entry.LastEditedAt.Time),
		CreatedAt:    openapi.NewOptNilDateTime(entry.CreatedAt.Time),
		UpdatedAt:    openapi.NewOptNilDateTime(entry.UpdatedAt.Time),
		ImageUrl:     openapi.NewOptNilString(entry.ImageUrl.String),
	}, nil
}

func (p *adminApiService) GetLinkPallet(ctx context.Context, params openapi.GetLinkPalletParams) (openapi.GetLinkPalletRes, error) {
	entry, err := p.queries.AdminGetEntryByPath(ctx, params.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get entry by path %s: %w", params.Path, err)
	}

	linkPallet, err := getLinkPalletData(ctx, p.db, p.queries, params.Path, entry.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to get link pallet data for path %s: %w", params.Path, err)
	}

	return linkPallet, nil
}

func (p *adminApiService) GetLinkedEntryPaths(ctx context.Context, params openapi.GetLinkedEntryPathsParams) (openapi.GetLinkedEntryPathsRes, error) {
	links, err := p.queries.GetLinkedEntries(ctx, params.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get linked entries for path %s: %w", params.Path, err)
	}

	result := make(openapi.LinkedEntryPathsResponse)
	for _, link := range links {
		result[strings.ToLower(link.DstTitle)] = link.Path.String
	}
	return &result, nil
}

func (p *adminApiService) UpdateEntryBody(ctx context.Context, req *openapi.UpdateEntryBodyRequest, params openapi.UpdateEntryBodyParams) (openapi.UpdateEntryBodyRes, error) {
	tx, err := p.db.Begin()
	if err != nil {
		slog.Error("failed to begin transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", slog.String("error", err.Error()))
		}
	}()

	// クエリの準備
	qtx := p.queries.WithTx(tx)

	affectedRows, err := qtx.UpdateEntryBody(ctx, admindb.UpdateEntryBodyParams{
		Path: params.Path,
		Body: req.Body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update entry body for path %s: %w", params.Path, err)
	}
	if affectedRows == 0 {
		return nil, fmt.Errorf("entry not found")
	}

	newEntry, err := qtx.AdminGetEntryByPath(ctx, params.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get entry by path %s: %w", params.Path, err)
	}

	if err := updateEntryLink(ctx, tx, qtx, params.Path, newEntry.Title, req.Body); err != nil {
		return nil, fmt.Errorf("failed to update entry links for path %s: %w", params.Path, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	go func() {
		slog.Info("Starting to generate entry_image")
		err := NewEntryImageWorker(p.queries).processEntryImages(context.Background())
		if err != nil {
			slog.Error("failed to process entry images", slog.String("error", err.Error()))
		}
	}()

	return &openapi.EmptyResponse{}, nil
}

// extractLinks extracts links from the Markdown text.
func extractLinks(markdown string) []string {
	re := regexp.MustCompile(`\[\[(.+?)]]`)
	matches := re.FindAllStringSubmatch(markdown, -1)
	seen := make(map[string]struct{})
	var result []string

	for _, match := range matches {
		link := strings.TrimSpace(match[1])
		lowerLink := strings.ToLower(link)
		if _, exists := seen[lowerLink]; !exists {
			seen[lowerLink] = struct{}{}
			result = append(result, link)
		}
	}

	return result
}

// updateEntryLink updates the entry_link table within a transaction.
func updateEntryLink(ctx context.Context, tx *sql.Tx, qtx *admindb.Queries, path string, title string, body string) error {
	// Extract links from the body, filtering out the title.
	links := extractLinks(body)
	var filteredLinks []string
	for _, link := range links {
		if !strings.EqualFold(link, title) {
			filteredLinks = append(filteredLinks, link)
		}
	}

	// Delete current links for the given path.
	if _, err := qtx.DeleteEntryLinkByPath(ctx, path); err != nil {
		return fmt.Errorf("failed to delete existing links for path %s: %w", path, err)
	}

	// Insert new links into the entry_link table.
	if len(filteredLinks) > 0 {
		var values []interface{}
		var placeholders []string
		for _, link := range filteredLinks {
			values = append(values, path, link)
			placeholders = append(placeholders, "(?, ?)")
		}
		query := `
			INSERT INTO entry_link (src_path, dst_title)
			VALUES ` + strings.Join(placeholders, ", ")
		if _, err := tx.ExecContext(ctx, query, values...); err != nil {
			return fmt.Errorf("failed to insert new links for path %s: %w", path, err)
		}
	}

	return nil
}

func (p *adminApiService) UpdateEntryTitle(ctx context.Context, req *openapi.UpdateEntryTitleRequest, params openapi.UpdateEntryTitleParams) (openapi.UpdateEntryTitleRes, error) {
	_, err := p.queries.UpdateEntryTitle(ctx, admindb.UpdateEntryTitleParams{
		Path:  params.Path,
		Title: req.Title,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update entry title for path %s: %w", params.Path, err)
	}
	return &openapi.EmptyResponse{}, nil
}

func (p *adminApiService) GetAllEntryTitles(ctx context.Context) (openapi.GetAllEntryTitlesRes, error) {
	titles, err := p.queries.GetAllEntryTitles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all entry titles: %w", err)
	}

	resp := openapi.EntryTitlesResponse(titles)
	return &resp, nil
}

func getDefaultTitle() string {
	return time.Now().Format("20060102150405")
}

func (p *adminApiService) CreateEntry(ctx context.Context, req *openapi.CreateEntryRequest) (openapi.CreateEntryRes, error) {
	now := time.Now()
	path := now.Format("2006/01/02/150405")

	_, err := p.queries.CreateEmptyEntry(ctx, admindb.CreateEmptyEntryParams{
		Path:  path,
		Title: req.Title.Or(getDefaultTitle()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}
	return &openapi.CreateEntryResponse{
		Path: path,
	}, nil
}

func (p *adminApiService) DeleteEntry(ctx context.Context, params openapi.DeleteEntryParams) (openapi.DeleteEntryRes, error) {
	_, err := p.queries.DeleteEntry(ctx, params.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete entry for path %s: %w", params.Path, err)
	}
	return &openapi.EmptyResponse{}, nil
}

func (p *adminApiService) UpdateEntryVisibility(ctx context.Context, req *openapi.UpdateVisibilityRequest, params openapi.UpdateEntryVisibilityParams) (openapi.UpdateEntryVisibilityRes, error) {
	slog.Info("UpdateEntryVisibility", slog.Any("request", req))
	tx, err := p.db.Begin()
	if err != nil {
		slog.Error("failed to begin transaction", slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", slog.String("error", err.Error()))
		}
	}()

	// クエリの準備
	qtx := p.queries.WithTx(tx)

	// 現在の可視性と公開日時を取得
	entry, err := qtx.GetEntryVisibility(ctx, params.Path)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("entry not found")
		}
		return nil, fmt.Errorf("failed to query entry: %w", err)
	}

	// 可視性の更新
	if err := qtx.UpdateVisibility(ctx, admindb.UpdateVisibilityParams{
		Visibility: admindb.EntryVisibility(req.Visibility),
		Path:       params.Path,
	}); err != nil {
		return nil, fmt.Errorf("failed to update visibility: %w", err)
	}

	// 可視性がprivateからpublicに変わり、published_atがnullの場合、published_atを現在時刻に設定
	if entry.Visibility == "private" && req.Visibility == "public" && !entry.PublishedAt.Valid {
		if err := qtx.UpdatePublishedAt(ctx, params.Path); err != nil {
			return nil, fmt.Errorf("failed to update published_at: %w", err)
		}
	}

	// ここで body に amzn.to の短縮URLがあれば､amazon の商品情報を取得してキャッシュします｡
	newEntry, err := qtx.AdminGetEntryByPath(ctx, params.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get entry by path %s: %w", params.Path, err)
	}
	rewroteBody := rewriteAmazonShortUrlInMarkdown(newEntry.Body)
	if rewroteBody != newEntry.Body {
		if _, err := qtx.UpdateEntryBody(ctx, admindb.UpdateEntryBodyParams{
			Path: params.Path,
			Body: rewroteBody,
		}); err != nil {
			return nil, fmt.Errorf("failed to update entry body for path %s: %w", params.Path, err)
		}
	}

	// 次に､amazon の画像キャッシュを更新します｡
	// これはバックグラウンドで処理してかまいません｡
	go func() {
		ctx := context.Background()
		slog.Info("Starting to get amazon cache")
		err := p.getAmazonCache(newEntry.Body, ctx)
		if err != nil {
			slog.Error("failed to get amazon cache", slog.String("error", err.Error()))
		}

		// update entry_image after that.
		err = NewEntryImageWorker(p.queries).processEntryImages(ctx)
		if err != nil {
			slog.Error("failed to process entry images", slog.String("error", err.Error()))
		}
	}()

	// トランザクションのコミット
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send notification to Hub
	for _, hubUrl := range p.hubUrls {
		slog.Info("Notify Hub", slog.String("hubUrl", hubUrl))
		if err := NotifyHub(hubUrl, "https://blog.64p.org/feed"); err != nil {
			slog.Error("Failed to notify Hub", slog.String("error", err.Error()))
		}
	}

	return &openapi.UpdateVisibilityResponse{
		Visibility: req.Visibility,
	}, nil
}

func (p *adminApiService) getAmazonCache(markdown string, ctx context.Context) error {
	// extract asin from body
	re := regexp.MustCompile(`asin:([A-Z0-9]+):detail`)
	matches := re.FindAllStringSubmatch(markdown, -1)
	var asins []string
	for _, match := range matches {
		asin := match[1]
		count, err := p.queries.CountAmazonCacheByAsin(ctx, asin)
		if err != nil {
			return fmt.Errorf("failed to count amazon cache for ASIN %s: %w", asin, err)
		}
		if count > 0 {
			slog.Info("ASIN is already cached", slog.String("asin", asin))
			continue
		}

		asins = append(asins, asin)
	}

	slog.Info("getAmazonCache", slog.Any("asins", asins))

	// バッチサイズは10
	const batchSize = 10
	// バッチ処理間の待機時間は1分
	const waitDuration = 1 * time.Minute

	// ASINsをバッチに分割して処理
	for i := 0; i < len(asins); i += batchSize {
		end := i + batchSize
		if end > len(asins) {
			end = len(asins)
		}

		// 現在のバッチを取得
		currentBatch := asins[i:end]

		// バッチ処理を実行
		productDetails, err := p.paapiClient.FetchAmazonProductDetails(ctx, currentBatch)
		if err != nil {
			slog.Error("failed to fetch amazon product details", slog.String("error", err.Error()))
			return fmt.Errorf("failed to fetch amazon product details: %w", err)
		}

		for _, productDetail := range productDetails {
			slog.Info("Amazon product detail", slog.String("asin", productDetail.ASIN), slog.String("title", productDetail.Title))
			_, err := p.queries.InsertAmazonProductDetail(ctx, admindb.InsertAmazonProductDetailParams{
				Asin:           productDetail.ASIN,
				Title:          sql.NullString{String: productDetail.Title, Valid: true},
				ImageMediumUrl: sql.NullString{String: productDetail.ImageMediumURL, Valid: true},
				Link:           productDetail.Link,
			})
			if err != nil {
				return fmt.Errorf("failed to insert amazon product detail for ASIN %s: %w", productDetail.ASIN, err)
			}
		}

		// 最後のバッチでなければ待機
		if end < len(asins) {
			time.Sleep(waitDuration)
		}
	}
	slog.Info("getAmazonCache: done")
	return nil
}

func (p *adminApiService) UploadFile(ctx context.Context, req *openapi.UploadFileBodyMultipart) (openapi.UploadFileRes, error) {
	// Content-Type の取得
	contentType := req.File.Header.Get("Content-Type")
	contentLength := req.File.Size

	now := time.Now().UnixMilli()
	key := fmt.Sprintf("%d-%s", now, req.File.Name)

	slog.Info("UploadPost", slog.String("contentType", contentType), slog.String("key", key), slog.Int("size", int(req.File.Size)))

	// 先頭100バイトをキャプチャするためのバッファ
	var previewBuf bytes.Buffer
	previewReader := io.TeeReader(io.LimitReader(req.File.File, 100), &previewBuf)

	// 残りのデータと組み合わせるためのMultiReader
	reader := io.MultiReader(
		&previewBuf,
		req.File.File,
	)

	// プレビューデータを読み込んでURLエスケープ
	preview := make([]byte, req.File.Size)
	n, _ := previewReader.Read(preview)
	escapedPreview := url.QueryEscape(string(preview[:n]))

	slog.Debug("UploadPost details",
		slog.String("contentType", contentType),
		slog.String("key", key),
		slog.Int("size", int(req.File.Size)),
		slog.String("preview", escapedPreview),
		slog.Int("bufSize", len(preview)))

	// S3にアップロード
	err := p.S3Client.PutObjectToAttachmentBucket(
		ctx,
		key,
		contentType,
		contentLength,
		reader,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	url, err := url.Parse(`https://blog-attachments.64p.org/` + key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	slog.Info("UploadPost completed", slog.String("url", url.String()))

	return &openapi.UploadFileResponse{
		URL: *url,
	}, nil
}

func (p *adminApiService) NewError(_ context.Context, err error) *openapi.ErrorResponseStatusCode {
	slog.Error("NewError", slog.String("error", err.Error()))
	if errors.Is(err, sql.ErrNoRows) {
		return &openapi.ErrorResponseStatusCode{
			StatusCode: 404,
			Response: openapi.ErrorResponse{
				Message: openapi.NewOptString("Not Found"),
				Error:   openapi.NewOptString(fmt.Sprintf("Not Found: %v", err)),
			},
		}
	} else {
		return &openapi.ErrorResponseStatusCode{
			StatusCode: 500,
			Response: openapi.ErrorResponse{
				Message: openapi.NewOptString("Internal Server Error"),
				Error:   openapi.NewOptString(fmt.Sprintf("Internal Server Error: %v", err)),
			},
		}
	}
}

func (p *adminApiService) RegenerateEntryImage(ctx context.Context, params openapi.RegenerateEntryImageParams) (openapi.RegenerateEntryImageRes, error) {
	_, err := p.queries.DeleteEntryImageByPath(ctx, params.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to delete entry image for path %s: %w", params.Path, err)
	}

	service := server.NewEntryImageService(p.queries)
	entries, err := service.GetEntryImageNotProcessedEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get entry image not processed entries: %w", err)
	}
	for _, entry := range entries {
		if entry.Path == params.Path {
			err = service.ProcessEntry(ctx, entry)
			if err != nil {
				return nil, fmt.Errorf("failed to process entry image for path %s: %w", params.Path, err)
			}
		}
	}
	return &openapi.EmptyResponse{}, nil
}

// AuthLogin implements the login endpoint
func (p *adminApiService) AuthLogin(ctx context.Context, req *openapi.LoginRequest) (openapi.AuthLoginRes, error) {
	// Validate credentials using constant-time comparison
	usernameMatch := subtle.ConstantTimeCompare([]byte(req.Username), []byte(p.adminUser)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(req.Password), []byte(p.adminPassword)) == 1

	if !usernameMatch || !passwordMatch {
		// Return 401 with ErrorResponse
		return &openapi.ErrorResponse{
			Message: openapi.NewOptString("Invalid username or password"),
		}, nil
	}

	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Calculate session expiry
	sessionTimeout := defaultSessionTimeout
	expires := time.Now().Add(sessionTimeout)

	// Create session in database
	err = p.queries.CreateSession(ctx, admindb.CreateSessionParams{
		SessionID: sessionID,
		Username:  req.Username,
		ExpiresAt: expires,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Get response writer from context to set cookie
	if w, ok := GetHTTPResponse(ctx); ok {
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionID,
			Path:     "/admin",
			HttpOnly: true,
			Secure:   p.isSecure,
			SameSite: http.SameSiteStrictMode,
			Expires:  expires,
			MaxAge:   int(time.Until(expires).Seconds()),
		})
	}

	// Clean up expired sessions (1% probability)
	if rand.Float32() < 0.01 {
		go func() {
			ctx := context.Background()
			if err := p.queries.DeleteExpiredSessions(ctx); err != nil {
				slog.Error("failed to delete expired sessions", slog.String("error", err.Error()))
			}
		}()
	}

	return &openapi.LoginResponse{
		Success: true,
		Message: openapi.NewOptString("Login successful"),
	}, nil
}

// AuthLogout implements the logout endpoint
func (p *adminApiService) AuthLogout(ctx context.Context) (openapi.AuthLogoutRes, error) {
	// Get session ID from cookie
	var sessionID string
	if r, ok := GetHTTPRequest(ctx); ok {
		if cookie, err := r.Cookie(sessionCookieName); err == nil {
			sessionID = cookie.Value
		}
	}

	if sessionID != "" {
		// Delete session from database
		err := p.queries.DeleteSession(ctx, sessionID)
		if err != nil {
			slog.Error("failed to delete session", slog.String("error", err.Error()))
		}
	}

	// Clear session cookie
	if w, ok := GetHTTPResponse(ctx); ok {
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/admin",
			HttpOnly: true,
			Secure:   p.isSecure,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
		})
	}

	return &openapi.EmptyResponse{
		Message: openapi.NewOptString("Logout successful"),
	}, nil
}

// AuthCheck implements the auth check endpoint
func (p *adminApiService) AuthCheck(ctx context.Context) (openapi.AuthCheckRes, error) {
	// Get session ID from cookie
	var sessionID string
	if r, ok := GetHTTPRequest(ctx); ok {
		if cookie, err := r.Cookie(sessionCookieName); err == nil {
			sessionID = cookie.Value
		}
	}

	if sessionID == "" {
		return &openapi.ErrorResponse{
			Message: openapi.NewOptString("No session found"),
		}, nil
	}

	// Get session from database
	session, err := p.queries.GetSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &openapi.ErrorResponse{
				Message: openapi.NewOptString("Invalid session"),
			}, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Update last accessed time
	go func() {
		ctx := context.Background()
		if err := p.queries.UpdateSessionLastAccessed(ctx, sessionID); err != nil {
			slog.Error("failed to update session last accessed", slog.String("error", err.Error()))
		}
	}()

	return &openapi.AuthCheckResponse{
		Authenticated: true,
		Username:      openapi.NewOptString(session.Username),
	}, nil
}

func (p *adminApiService) GetBuildInfo(ctx context.Context) (openapi.GetBuildInfoRes, error) {
	buildInfo, err := server.ReadBuildInfo()
	if err != nil {
		return &openapi.ErrorResponse{
			Message: openapi.NewOptString("Failed to read build info"),
			Error:   openapi.NewOptString(err.Error()),
		}, nil
	}

	return &openapi.BuildInfoBuildInfo{
		BuildTime:      buildInfo.BuildTime,
		GitCommit:      buildInfo.GitCommit,
		GitShortCommit: buildInfo.GitShortCommit,
		GitBranch:      buildInfo.GitBranch,
		GitTag:         openapi.NewOptString(buildInfo.GitTag),
		GithubUrl:      buildInfo.GithubUrl,
	}, nil
}
