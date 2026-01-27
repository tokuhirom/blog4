package ogimage

import (
	"bytes"
	"context"
	"io"

	"github.com/tokuhirom/blog4/internal/sobs"
)

// SobsAdapter adapts SobsClient to the S3Uploader interface
type SobsAdapter struct {
	client *sobs.SobsClient
}

// NewSobsAdapter creates a new adapter
func NewSobsAdapter(client *sobs.SobsClient) *SobsAdapter {
	return &SobsAdapter{client: client}
}

// PutObject uploads a file to S3 using SobsClient
func (a *SobsAdapter) PutObject(ctx context.Context, key string, body io.Reader, contentType string, metadata map[string]string) error {
	// Read the body to get content length (required by SobsClient)
	buf := new(bytes.Buffer)
	n, err := buf.ReadFrom(body)
	if err != nil {
		return err
	}

	// Note: metadata is not used by SobsClient's PutObjectToAttachmentBucket
	// If metadata support is needed, the SobsClient interface would need to be extended
	return a.client.PutObjectToAttachmentBucket(ctx, key, contentType, n, buf)
}
