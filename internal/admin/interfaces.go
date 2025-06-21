package admin

import (
	"context"
	"database/sql"
	"io"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/internal/admin/openapi"
)

// AdminStore defines the database operations needed by adminApiService
//
//go:generate mockgen -source=interfaces.go -destination=mocks/mock_interfaces.go -package=mocks -mock_names=AdminStore=MockAdminStore,TxManager=MockTxManager,AmazonClient=MockAmazonClient,StorageClient=MockStorageClient,HubNotifier=MockHubNotifier,EntryImageProcessor=MockEntryImageProcessor,LinkPalletService=MockLinkPalletService
type AdminStore interface {
	// Entry operations
	GetLatestEntries(ctx context.Context, params admindb.GetLatestEntriesParams) ([]admindb.GetLatestEntriesRow, error)
	AdminGetEntryByPath(ctx context.Context, path string) (admindb.AdminGetEntryByPathRow, error)
	GetLinkedEntries(ctx context.Context, path string) ([]admindb.GetLinkedEntriesRow, error)
	UpdateEntryBody(ctx context.Context, params admindb.UpdateEntryBodyParams) (int64, error)
	UpdateEntryTitle(ctx context.Context, params admindb.UpdateEntryTitleParams) (int64, error)
	GetAllEntryTitles(ctx context.Context) ([]string, error)
	CreateEmptyEntry(ctx context.Context, params admindb.CreateEmptyEntryParams) (int64, error)
	DeleteEntry(ctx context.Context, path string) (int64, error)
	GetEntryVisibility(ctx context.Context, path string) (admindb.GetEntryVisibilityRow, error)
	UpdateVisibility(ctx context.Context, params admindb.UpdateVisibilityParams) error
	UpdatePublishedAt(ctx context.Context, path string) error
	DeleteEntryLinkByPath(ctx context.Context, path string) (int64, error)
	DeleteEntryImageByPath(ctx context.Context, path string) (int64, error)

	// Amazon cache operations
	CountAmazonCacheByAsin(ctx context.Context, asin string) (int64, error)
	InsertAmazonProductDetail(ctx context.Context, params admindb.InsertAmazonProductDetailParams) (int64, error)

	// Transaction support
	WithTx(tx *sql.Tx) AdminStore
}

// TxManager handles database transactions
type TxManager interface {
	Begin() (*sql.Tx, error)
}

// AmazonClient defines Amazon product API operations
type AmazonClient interface {
	FetchAmazonProductDetails(ctx context.Context, asins []string) ([]ProductDetail, error)
}

// ProductDetail represents Amazon product information
type ProductDetail struct {
	ASIN           string
	Title          string
	ImageMediumURL string
	Link           string
}

// StorageClient defines file storage operations
type StorageClient interface {
	PutObjectToAttachmentBucket(ctx context.Context, key string, contentType string, contentLength int64, reader io.Reader) error
}

// HubNotifier defines hub notification operations
type HubNotifier interface {
	NotifyHub(hubUrl string, topicUrl string) error
}

// EntryImageProcessor defines entry image processing operations
type EntryImageProcessor interface {
	GetEntryImageNotProcessedEntries(ctx context.Context) ([]admindb.Entry, error)
	ProcessEntry(ctx context.Context, entry admindb.Entry) error
}

// LinkPalletService defines link pallet operations
type LinkPalletService interface {
	GetLinkPalletData(ctx context.Context, txManager TxManager, store AdminStore, path string, title string) (*openapi.LinkPalletData, error)
}
