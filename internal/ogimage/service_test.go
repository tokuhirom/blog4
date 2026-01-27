package ogimage

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/internal/ogimage/mocks"
	"go.uber.org/mock/gomock"
)

func TestEnsureOGImage_ImageExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockEntryImageStore(ctrl)
	generator := NewGenerator(&mockS3Uploader{}, "https://example.com", "")
	service := NewService(generator, mockStore)

	// Mock: image already exists
	mockStore.EXPECT().
		GetEntryImageByPath(gomock.Any(), "2026/01/test").
		Return(admindb.EntryImage{
			Path: "2026/01/test",
			Url:  sql.NullString{String: "https://example.com/og-images/test.png", Valid: true},
		}, nil)

	err := service.EnsureOGImage(context.Background(), "2026/01/test")
	require.NoError(t, err)
}

func TestEnsureOGImage_GeneratesImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockEntryImageStore(ctrl)
	mockS3 := &mockS3Uploader{}
	generator := NewGenerator(mockS3, "https://example.com", "")
	service := NewService(generator, mockStore)

	// Mock: no image exists
	mockStore.EXPECT().
		GetEntryImageByPath(gomock.Any(), "2026/01/test").
		Return(admindb.EntryImage{}, sql.ErrNoRows)

	// Mock: get entry info
	mockStore.EXPECT().
		AdminGetEntryByPath(gomock.Any(), "2026/01/test").
		Return(admindb.AdminGetEntryByPathRow{
			Path:        "2026/01/test",
			Title:       "Test Entry",
			PublishedAt: sql.NullTime{Time: time.Date(2026, 1, 27, 12, 0, 0, 0, time.UTC), Valid: true},
		}, nil)

	// Mock: insert entry image
	mockStore.EXPECT().
		InsertEntryImage(gomock.Any(), gomock.Any()).
		Return(int64(1), nil)

	err := service.EnsureOGImage(context.Background(), "2026/01/test")
	require.NoError(t, err)

	// Verify S3 upload happened
	assert.Contains(t, mockS3.uploadedKey, "og-images/")
	assert.Equal(t, "image/png", mockS3.uploadedContentType)
}

func TestEnsureOGImage_EntryNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockEntryImageStore(ctrl)
	generator := NewGenerator(&mockS3Uploader{}, "https://example.com", "")
	service := NewService(generator, mockStore)

	// Mock: no image exists
	mockStore.EXPECT().
		GetEntryImageByPath(gomock.Any(), "2026/01/nonexistent").
		Return(admindb.EntryImage{}, sql.ErrNoRows)

	// Mock: entry not found
	mockStore.EXPECT().
		AdminGetEntryByPath(gomock.Any(), "2026/01/nonexistent").
		Return(admindb.AdminGetEntryByPathRow{}, sql.ErrNoRows)

	err := service.EnsureOGImage(context.Background(), "2026/01/nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get entry")
}

func TestEnsureOGImage_InsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockEntryImageStore(ctrl)
	mockS3 := &mockS3Uploader{}
	generator := NewGenerator(mockS3, "https://example.com", "")
	service := NewService(generator, mockStore)

	// Mock: no image exists
	mockStore.EXPECT().
		GetEntryImageByPath(gomock.Any(), "2026/01/test").
		Return(admindb.EntryImage{}, sql.ErrNoRows)

	// Mock: get entry info
	mockStore.EXPECT().
		AdminGetEntryByPath(gomock.Any(), "2026/01/test").
		Return(admindb.AdminGetEntryByPathRow{
			Path:        "2026/01/test",
			Title:       "Test Entry",
			PublishedAt: sql.NullTime{Time: time.Date(2026, 1, 27, 12, 0, 0, 0, time.UTC), Valid: true},
		}, nil)

	// Mock: insert fails
	mockStore.EXPECT().
		InsertEntryImage(gomock.Any(), gomock.Any()).
		Return(int64(0), errors.New("database error"))

	err := service.EnsureOGImage(context.Background(), "2026/01/test")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to insert entry image")
}
