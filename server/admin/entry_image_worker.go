package admin

import (
	"context"
	"fmt"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
	"log/slog"
)

type EntryImageWorker struct {
	service *server.EntryImageService
}

func NewEntryImageWorker(queries *admindb.Queries) *EntryImageWorker {
	return &EntryImageWorker{
		service: server.NewEntryImageService(queries),
	}
}

func (w *EntryImageWorker) processEntryImages(ctx context.Context) error {
	entries, err := w.service.GetEntryImageNotProcessedEntries(ctx)
	if err != nil {
		return fmt.Errorf("failed to get unprocessed entries for image generation: %w", err)
	}
	for _, entry := range entries {
		if err := w.service.ProcessEntry(ctx, entry); err != nil {
			slog.Error("Failed to process entry image", slog.String("path", entry.Path), slog.Any("error", err))
		}
	}
	return nil
}
