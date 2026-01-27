package ogimage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

//go:generate go run go.uber.org/mock/mockgen -source=service.go -destination=mocks/mock_service.go -package=mocks

// EntryImageStore defines the database operations needed for entry images
type EntryImageStore interface {
	GetEntryImageByPath(ctx context.Context, path string) (admindb.EntryImage, error)
	AdminGetEntryByPath(ctx context.Context, path string) (admindb.AdminGetEntryByPathRow, error)
	InsertEntryImage(ctx context.Context, arg admindb.InsertEntryImageParams) (int64, error)
}

// Service orchestrates OG image generation
type Service struct {
	generator *Generator
	store     EntryImageStore
}

// NewService creates a new Service
func NewService(generator *Generator, store EntryImageStore) *Service {
	return &Service{
		generator: generator,
		store:     store,
	}
}

// EnsureOGImage checks if an entry has an image, and generates one if needed
func (s *Service) EnsureOGImage(ctx context.Context, path string) error {
	// Check if entry_image already exists
	entryImage, err := s.store.GetEntryImageByPath(ctx, path)
	if err == nil && entryImage.Url.Valid && entryImage.Url.String != "" {
		// Image already exists
		slog.Info("Entry image already exists", slog.String("path", path), slog.String("url", entryImage.Url.String))
		return nil
	}

	if err != nil && err != sql.ErrNoRows {
		// Unexpected database error
		return fmt.Errorf("failed to check entry image: %w", err)
	}

	// Get entry information
	entry, err := s.store.AdminGetEntryByPath(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to get entry: %w", err)
	}

	// Generate OG image
	entryInfo := EntryInfo{
		Path:        path,
		Title:       entry.Title,
		PublishedAt: entry.PublishedAt.Time,
	}

	url, err := s.generator.GenerateOGImage(ctx, entryInfo)
	if err != nil {
		return fmt.Errorf("failed to generate OG image: %w", err)
	}

	// Insert into database
	_, err = s.store.InsertEntryImage(ctx, admindb.InsertEntryImageParams{
		Path: path,
		Url:  sql.NullString{String: url, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to insert entry image: %w", err)
	}

	slog.Info("Generated OG image", slog.String("path", path), slog.String("url", url))
	return nil
}
