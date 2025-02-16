package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server/admin/openapi"
	"log"
	"strings"
	"time"
)

type adminApiService struct {
	queries *admindb.Queries
	db      *sql.DB
}

func (p *adminApiService) GetLatestEntries(ctx context.Context, params openapi.GetLatestEntriesParams) ([]openapi.GetLatestEntriesRow, error) {
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
	log.Printf("GetLatestEntries %v", lastEditedAt)
	entries, err := p.queries.GetLatestEntries(ctx, admindb.GetLatestEntriesParams{
		Column1:      lastEditedAt,
		LastEditedAt: lastEditedAt,
		Limit:        100,
	})
	if err != nil {
		return nil, err
	}

	var result []openapi.GetLatestEntriesRow
	for _, entry := range entries {
		result = append(result, openapi.GetLatestEntriesRow{
			Path:         openapi.NewOptString(entry.Path),
			Title:        openapi.NewOptString(entry.Title),
			Body:         openapi.NewOptString(entry.Body),
			Visibility:   openapi.NewOptString(string(entry.Visibility)),
			Format:       openapi.NewOptString(string(entry.Format)),
			PublishedAt:  openapi.NewOptNilDateTime(entry.PublishedAt.Time),
			LastEditedAt: openapi.NewOptNilDateTime(entry.LastEditedAt.Time),
			CreatedAt:    openapi.NewOptNilDateTime(entry.CreatedAt.Time),
			UpdatedAt:    openapi.NewOptNilDateTime(entry.UpdatedAt.Time),
			ImageUrl:     openapi.NewOptNilString(entry.ImageUrl.String),
		})
	}
	return result, nil
}

func (p *adminApiService) GetEntryByDynamicPath(ctx context.Context, params openapi.GetEntryByDynamicPathParams) (*openapi.GetLatestEntriesRow, error) {
	entry, err := p.queries.AdminGetEntryByPath(ctx, params.Path)
	if err != nil {
		log.Printf("GetEntryByDynamicPath %v", err)
		return nil, err
	}

	return &openapi.GetLatestEntriesRow{
		Path:         openapi.NewOptString(entry.Path),
		Title:        openapi.NewOptString(entry.Title),
		Body:         openapi.NewOptString(entry.Body),
		Visibility:   openapi.NewOptString(string(entry.Visibility)),
		Format:       openapi.NewOptString(string(entry.Format)),
		PublishedAt:  openapi.NewOptNilDateTime(entry.PublishedAt.Time),
		LastEditedAt: openapi.NewOptNilDateTime(entry.LastEditedAt.Time),
		CreatedAt:    openapi.NewOptNilDateTime(entry.CreatedAt.Time),
		UpdatedAt:    openapi.NewOptNilDateTime(entry.UpdatedAt.Time),
		ImageUrl:     openapi.NewOptNilString(entry.ImageUrl.String),
	}, nil
}

func (p *adminApiService) GetLinkPallet(ctx context.Context, params openapi.GetLinkPalletParams) (*openapi.LinkPalletData, error) {
	entry, err := p.queries.AdminGetEntryByPath(ctx, params.Path)
	if err != nil {
		return nil, err
	}

	linkPallet, err := getLinkPalletData(ctx, p.db, p.queries, params.Path, entry.Title)
	if err != nil {
		return nil, err
	}

	return linkPallet, nil
}

func (p *adminApiService) GetLinkedEntryPaths(ctx context.Context, params openapi.GetLinkedEntryPathsParams) (openapi.GetLinkedEntryPathsRes, error) {
	links, err := p.queries.GetLinkedEntries(ctx, params.Path)
	if err != nil {
		return nil, err
	}

	result := make(openapi.LinkedEntryPathsResponse)
	for _, link := range links {
		result[strings.ToLower(link.DstTitle)] = link.Path.String
	}
	return &result, nil
}

func (p *adminApiService) UpdateEntryBody(ctx context.Context, req *openapi.UpdateEntryBodyRequest, params openapi.UpdateEntryBodyParams) (openapi.UpdateEntryBodyRes, error) {
	_, err := p.queries.UpdateEntryBody(ctx, admindb.UpdateEntryBodyParams{
		Path: params.Path,
		Body: req.Body,
	})
	if err != nil {
		return nil, err
	}
	return &openapi.EmptyResponse{}, nil
}

func (p *adminApiService) UpdateEntryTitle(ctx context.Context, req *openapi.UpdateEntryTitleRequest, params openapi.UpdateEntryTitleParams) (openapi.UpdateEntryTitleRes, error) {
	_, err := p.queries.UpdateEntryTitle(ctx, admindb.UpdateEntryTitleParams{
		Path:  params.Path,
		Title: req.Title,
	})
	if err != nil {
		return nil, err
	}
	return &openapi.EmptyResponse{}, nil
}

func (p *adminApiService) GetAllEntryTitles(ctx context.Context) (openapi.EntryTitlesResponse, error) {
	titles, err := p.queries.GetAllEntryTitles(ctx)
	if err != nil {
		return nil, err
	}

	return titles, nil
}

func (p *adminApiService) CreateEntry(ctx context.Context, req *openapi.CreateEntryRequest) (openapi.CreateEntryRes, error) {
	now := time.Now()
	path := now.Format("2006/01/02/150405")

	_, err := p.queries.CreateEmptyEntry(ctx, admindb.CreateEmptyEntryParams{
		Path:  path,
		Title: req.Title,
	})
	if err != nil {
		return nil, err
	}
	return &openapi.CreateEntryResponse{
		Path: path,
	}, nil
}

func (p *adminApiService) DeleteEntry(ctx context.Context, params openapi.DeleteEntryParams) (openapi.DeleteEntryRes, error) {
	_, err := p.queries.DeleteEntry(ctx, params.Path)
	if err != nil {
		return nil, err
	}
	return &openapi.EmptyResponse{}, nil
}

func (p *adminApiService) UpdateEntryVisibility(ctx context.Context, req *openapi.UpdateVisibilityRequest, params openapi.UpdateEntryVisibilityParams) (openapi.UpdateEntryVisibilityRes, error) {
	// トランザクションがうまく動かないので一旦コメントアウト

	tx, err := p.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %v", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("failed to rollback transaction: %v", err)
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

	// トランザクションのコミット
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// TODO ここで blog ping を送るとベターだが後でやる。
	/*
		for (const hubUrl of HUB_URLS) {
			console.log(`Notify Hub: ${hubUrl}`);
			await notifyHub(hubUrl, `https://blog.64p.org/feed`);
		}
	*/

	return &openapi.UpdateVisibilityResponse{
		Visibility: req.Visibility,
	}, nil
}

func (p *adminApiService) NewError(_ context.Context, err error) *openapi.ErrorResponseStatusCode {
	log.Printf("NewError %v", err)
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
