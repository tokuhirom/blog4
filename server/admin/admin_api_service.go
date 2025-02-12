package admin

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server/admin/openapi"
	"log"
)

type adminApiService struct {
	queries *admindb.Queries
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

func (p *adminApiService) NewError(_ context.Context, err error) *openapi.ErrorResponseStatusCode {
	log.Printf("NewError %v", err)
	return &openapi.ErrorResponseStatusCode{
		StatusCode: 500,
		Response: openapi.ErrorResponse{
			Message: openapi.NewOptString("Internal Server Error"),
			Error:   openapi.NewOptString(fmt.Sprintf("Internal Server Error: %v", err)),
		},
	}
}
