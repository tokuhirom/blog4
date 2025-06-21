package admin

import (
	"context"
	"database/sql"
	"io"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/internal/admin/openapi"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/sobs"
)

// AdminStoreAdapter wraps admindb.Queries to implement AdminStore interface
type AdminStoreAdapter struct {
	*admindb.Queries
}

// WithTx implements AdminStore.WithTx
func (a *AdminStoreAdapter) WithTx(tx *sql.Tx) AdminStore {
	return &AdminStoreAdapter{
		Queries: a.Queries.WithTx(tx),
	}
}

// TxManagerAdapter wraps sql.DB to implement TxManager interface
type TxManagerAdapter struct {
	db *sql.DB
}

func NewTxManagerAdapter(db *sql.DB) TxManager {
	return &TxManagerAdapter{db: db}
}

func (t *TxManagerAdapter) Begin() (*sql.Tx, error) {
	return t.db.Begin()
}

// AmazonClientAdapter wraps PAAPIClient to implement AmazonClient interface
type AmazonClientAdapter struct {
	client *PAAPIClient
}

func NewAmazonClientAdapter(client *PAAPIClient) AmazonClient {
	return &AmazonClientAdapter{client: client}
}

func (a *AmazonClientAdapter) FetchAmazonProductDetails(ctx context.Context, asins []string) ([]ProductDetail, error) {
	details, err := a.client.FetchAmazonProductDetails(ctx, asins)
	if err != nil {
		return nil, err
	}

	// Convert AmazonProductDetail to ProductDetail
	result := make([]ProductDetail, len(details))
	for i, detail := range details {
		result[i] = ProductDetail{
			ASIN:           detail.ASIN,
			Title:          detail.Title,
			ImageMediumURL: detail.ImageMediumURL,
			Link:           detail.Link,
		}
	}
	return result, nil
}

// StorageClientAdapter wraps sobs.SobsClient to implement StorageClient interface
type StorageClientAdapter struct {
	client *sobs.SobsClient
}

func NewStorageClientAdapter(client *sobs.SobsClient) StorageClient {
	return &StorageClientAdapter{client: client}
}

func (s *StorageClientAdapter) PutObjectToAttachmentBucket(ctx context.Context, key string, contentType string, contentLength int64, reader io.Reader) error {
	return s.client.PutObjectToAttachmentBucket(ctx, key, contentType, contentLength, reader)
}

// HubNotifierAdapter implements HubNotifier interface
type HubNotifierAdapter struct{}

func NewHubNotifierAdapter() HubNotifier {
	return &HubNotifierAdapter{}
}

func (h *HubNotifierAdapter) NotifyHub(hubUrl string, topicUrl string) error {
	return NotifyHub(hubUrl, topicUrl)
}

// EntryImageProcessorAdapter implements EntryImageProcessor interface
type EntryImageProcessorAdapter struct {
	service *server.EntryImageService
}

func NewEntryImageProcessorAdapter(store server.EntryImageStore) EntryImageProcessor {
	return &EntryImageProcessorAdapter{
		service: server.NewEntryImageService(store),
	}
}

func (e *EntryImageProcessorAdapter) GetEntryImageNotProcessedEntries(ctx context.Context) ([]admindb.Entry, error) {
	return e.service.GetEntryImageNotProcessedEntries(ctx)
}

func (e *EntryImageProcessorAdapter) ProcessEntry(ctx context.Context, entry admindb.Entry) error {
	return e.service.ProcessEntry(ctx, entry)
}

// LinkPalletServiceAdapter implements LinkPalletService interface
type LinkPalletServiceAdapter struct{}

func NewLinkPalletServiceAdapter() LinkPalletService {
	return &LinkPalletServiceAdapter{}
}

func (l *LinkPalletServiceAdapter) GetLinkPalletData(ctx context.Context, txManager TxManager, store AdminStore, path string, title string) (*openapi.LinkPalletData, error) {
	// We need to access the underlying sql.DB from TxManager
	// This is a limitation of the current design, but for now we'll cast it
	if adapter, ok := txManager.(*TxManagerAdapter); ok {
		return getLinkPalletData(ctx, adapter.db, store.(*AdminStoreAdapter).Queries, path, title)
	}
	// Fallback - this shouldn't happen in production
	return nil, nil
}
