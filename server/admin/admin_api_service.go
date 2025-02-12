package admin

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server/admin/openapi"
	"log"
	"time"
)

type adminApiService struct {
	queries *admindb.Queries
}

func (p *adminApiService) GetLatestEntries(ctx context.Context) ([]openapi.GetLatestEntriesRow, error) {
	entries, err := p.queries.GetLatestEntries(ctx, admindb.GetLatestEntriesParams{
		Column1: nil,
		LastEditedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Limit: 100,
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
