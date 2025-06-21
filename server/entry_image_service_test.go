package server

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server/mocks"
	"go.uber.org/mock/gomock"
)

func TestEntryImageService_GetEntryImageNotProcessedEntries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockEntryImageStore(ctrl)
	service := NewEntryImageService(mockStore)

	ctx := context.Background()
	expectedEntries := []admindb.Entry{
		{Path: "/test1", Title: "Test 1", Body: "Body 1"},
		{Path: "/test2", Title: "Test 2", Body: "Body 2"},
	}

	mockStore.EXPECT().
		GetEntryImageNotProcessedEntries(ctx).
		Return(expectedEntries, nil)

	entries, err := service.GetEntryImageNotProcessedEntries(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(entries) != len(expectedEntries) {
		t.Fatalf("expected %d entries, got %d", len(expectedEntries), len(entries))
	}
}

func TestEntryImageService_ProcessEntry(t *testing.T) {
	tests := []struct {
		name          string
		entry         admindb.Entry
		setupMocks    func(*mocks.MockEntryImageStore)
		expectedError bool
	}{
		{
			name: "successful processing with basic image",
			entry: admindb.Entry{
				Path:  "/test",
				Title: "Test Entry",
				Body:  "Some text ![Image](https://example.com/image.jpg) more text",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					InsertEntryImage(gomock.Any(), admindb.InsertEntryImageParams{
						Path: "/test",
						Url:  sql.NullString{String: "https://example.com/image.jpg", Valid: true},
					}).
					Return(int64(1), nil)
			},
			expectedError: false,
		},
		{
			name: "successful processing with image tag",
			entry: admindb.Entry{
				Path:  "/test",
				Title: "Test Entry",
				Body:  `<img src="https://blog-attachments.64p.org/image.png" style="width:100%">`,
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					InsertEntryImage(gomock.Any(), admindb.InsertEntryImageParams{
						Path: "/test",
						Url:  sql.NullString{String: "https://blog-attachments.64p.org/image", Valid: true},
					}).
					Return(int64(1), nil)
			},
			expectedError: false,
		},
		{
			name: "successful processing with gyazo image",
			entry: admindb.Entry{
				Path:  "/test",
				Title: "Test Entry",
				Body:  "[![Image from Gyazo](https://i.gyazo.com/abc123.jpg)](https://gyazo.com/abc123)",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					InsertEntryImage(gomock.Any(), admindb.InsertEntryImageParams{
						Path: "/test",
						Url:  sql.NullString{String: "https://i.gyazo.com/abc123.jpg", Valid: true},
					}).
					Return(int64(1), nil)
			},
			expectedError: false,
		},
		{
			name: "successful processing with ASIN",
			entry: admindb.Entry{
				Path:  "/test",
				Title: "Test Entry",
				Body:  "Check this product [asin:B00EXAMPLE:detail]",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					GetAmazonImageUrlByAsin(gomock.Any(), "B00EXAMPLE").
					Return(sql.NullString{String: "https://amazon.com/image.jpg", Valid: true}, nil)
				mockStore.EXPECT().
					InsertEntryImage(gomock.Any(), admindb.InsertEntryImageParams{
						Path: "/test",
						Url:  sql.NullString{String: "https://amazon.com/image.jpg", Valid: true},
					}).
					Return(int64(1), nil)
			},
			expectedError: false,
		},
		{
			name: "ASIN lookup fails",
			entry: admindb.Entry{
				Path:  "/test",
				Title: "Test Entry",
				Body:  "Product [asin:B00EXAMPLE:detail] is great",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					GetAmazonImageUrlByAsin(gomock.Any(), "B00EXAMPLE").
					Return(sql.NullString{Valid: false}, errors.New("not found"))
			},
			expectedError: true, // ASIN lookup error propagates
		},
		{
			name: "no image found without ASIN",
			entry: admindb.Entry{
				Path:  "/test",
				Title: "Test Entry",
				Body:  "Just plain text without any images",
			},
			setupMocks:    func(mockStore *mocks.MockEntryImageStore) {},
			expectedError: false,
		},
		{
			name: "insert error",
			entry: admindb.Entry{
				Path:  "/test",
				Title: "Test Entry",
				Body:  "![Image](https://example.com/image.jpg)",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					InsertEntryImage(gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("database error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockEntryImageStore(ctrl)
			service := NewEntryImageService(mockStore)

			tt.setupMocks(mockStore)

			err := service.ProcessEntry(context.Background(), tt.entry)
			if tt.expectedError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestEntryImageService_getImageFromEntry(t *testing.T) {
	tests := []struct {
		name        string
		entry       admindb.Entry
		setupMocks  func(*mocks.MockEntryImageStore)
		wantImage   *string
		wantError   bool
	}{
		{
			name: "extract image tag",
			entry: admindb.Entry{
				Body: `<img src="https://example.com/image.jpg" alt="test">`,
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {},
			wantImage:  stringPtr("https://example.com/image"),
			wantError:  false,
		},
		{
			name: "extract basic markdown image",
			entry: admindb.Entry{
				Body: "![Alt text](https://example.com/image.png)",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {},
			wantImage:  stringPtr("https://example.com/image.png"),
			wantError:  false,
		},
		{
			name: "extract gyazo image",
			entry: admindb.Entry{
				Body: "[![Image from Gyazo](https://i.gyazo.com/123.jpg)](https://gyazo.com/123)",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {},
			wantImage:  stringPtr("https://i.gyazo.com/123.jpg"),
			wantError:  false,
		},
		{
			name: "extract ASIN image",
			entry: admindb.Entry{
				Body: "[asin:B00ABC123:detail]",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					GetAmazonImageUrlByAsin(gomock.Any(), "B00ABC123").
					Return(sql.NullString{String: "https://amazon.com/product.jpg", Valid: true}, nil)
			},
			wantImage: stringPtr("https://amazon.com/product.jpg"),
			wantError: false,
		},
		{
			name: "ASIN lookup error",
			entry: admindb.Entry{
				Body: "[asin:B00ABC123:detail]",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {
				mockStore.EXPECT().
					GetAmazonImageUrlByAsin(gomock.Any(), "B00ABC123").
					Return(sql.NullString{}, errors.New("database error"))
			},
			wantImage: nil,
			wantError: true,
		},
		{
			name: "no image found",
			entry: admindb.Entry{
				Body: "Just plain text",
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {},
			wantImage:  nil,
			wantError:  false,
		},
		{
			name: "multiple images - takes first",
			entry: admindb.Entry{
				Body: `<img src="https://first.com/image.png"> and ![second](https://second.com/image.jpg)`,
			},
			setupMocks: func(mockStore *mocks.MockEntryImageStore) {},
			wantImage:  stringPtr("https://first.com/image"),
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStore := mocks.NewMockEntryImageStore(ctrl)
			service := NewEntryImageService(mockStore)

			tt.setupMocks(mockStore)

			image, err := service.getImageFromEntry(context.Background(), tt.entry)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if (image == nil) != (tt.wantImage == nil) {
					t.Fatalf("image nil mismatch: got %v, want %v", image, tt.wantImage)
				}
				if image != nil && tt.wantImage != nil && *image != *tt.wantImage {
					t.Fatalf("image mismatch: got %q, want %q", *image, *tt.wantImage)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}