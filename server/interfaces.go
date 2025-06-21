package server

import (
	"context"
	"database/sql"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

// EntryImageStore defines the database operations needed by EntryImageService
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_interfaces.go -package=mocks
type EntryImageStore interface {
	GetEntryImageNotProcessedEntries(ctx context.Context) ([]admindb.Entry, error)
	GetAmazonImageUrlByAsin(ctx context.Context, asin string) (sql.NullString, error)
	InsertEntryImage(ctx context.Context, params admindb.InsertEntryImageParams) (int64, error)
}