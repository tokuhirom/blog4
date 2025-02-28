package admin

import (
	"context"
	"github.com/tokuhirom/blog4/db/admin/admindb"
	"github.com/tokuhirom/blog4/server"
	"log"
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
		return err
	}
	for _, entry := range entries {
		if err := w.service.ProcessEntry(ctx, entry); err != nil {
			log.Printf("Failed to process entry: %v", err)
		}
	}
	return nil
}
