package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/tokuhirom/blog4/db/admin/admindb"
)

type EntryImageService struct {
	store EntryImageStore
}

func NewEntryImageService(store EntryImageStore) *EntryImageService {
	return &EntryImageService{
		store: store,
	}
}

func (w *EntryImageService) GetEntryImageNotProcessedEntries(ctx context.Context) ([]admindb.Entry, error) {
	return w.store.GetEntryImageNotProcessedEntries(ctx)
}

func (w *EntryImageService) ProcessEntry(ctx context.Context, entry admindb.Entry) error {
	slog.Info("Processing entry", slog.String("path", entry.Path), slog.String("title", entry.Title))

	image, err := w.getImageFromEntry(ctx, entry)
	if err != nil {
		return fmt.Errorf("failed to get image from entry %s: %w", entry.Path, err)
	}

	if image == nil {
		if strings.Contains(entry.Body, "[asin:") {
			slog.Warn("image is not available, maybe ASIN processing is delayed", slog.String("path", entry.Path))
			return nil
		} else {
			slog.Error("image is not available", slog.String("path", entry.Path))
			return nil
		}
	} else {
		slog.Info("image is available", slog.String("path", entry.Path), slog.String("title", entry.Title), slog.String("image", *image))
		_, err := w.store.InsertEntryImage(ctx, admindb.InsertEntryImageParams{
			Path: entry.Path,
			Url:  sql.NullString{String: *image, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to insert entry image for path %s: %w", entry.Path, err)
		}
		return nil
	}
}

func (w *EntryImageService) getImageFromEntry(ctx context.Context, entry admindb.Entry) (*string, error) {
	// extract image url from entry body.

	// image pattern is following:
	// ![Image](https://blog-attachments.64p.org/1735866754875-image.png)
	// [![Image from Gyazo](https://i.gyazo.com/d58c72d37ca373ab293184cdb5e6e6bb.jpg)](https://gyazo.com/d58c72d37ca373ab293184cdb5e6e6bb)
	// [asin:4022520221:detail]
	// <img src="https://blog-attachments.64p.org/20240318-08591291046a4b-d1c4-4282-8998-fba07edb19a6.png" style="width:100%">
	//
	// Note, ASIN pattern is not just a URL. So, get the concrete image URL from the database.
	// if it's not available yet, skip it.

	imageTagRe := regexp.MustCompile(`<img[^>]*src=['"]?(https?://[^\s)]+)\.(?:jpg|png|gif)['"]?`)
	basicImageRe := regexp.MustCompile(`!\[.*?]\((https?://[^\s)]+)\)`)
	gyazoImageRe := regexp.MustCompile(`\[!\[.*?]\((https?://[^\s)]+)\)]\((.*?)\)`)
	asinRe := regexp.MustCompile(`\[asin:([A-Z0-9]+):detail]`)

	// Match image tag pattern
	if imageTagMatch := imageTagRe.FindStringSubmatch(entry.Body); imageTagMatch != nil {
		imageURL := imageTagMatch[1]
		return &imageURL, nil
	}

	// Match basic image pattern
	if basicImageMatch := basicImageRe.FindStringSubmatch(entry.Body); basicImageMatch != nil {
		imageURL := basicImageMatch[1]
		return &imageURL, nil
	}

	// Match Gyazo image pattern
	if gyazoImageMatch := gyazoImageRe.FindStringSubmatch(entry.Body); gyazoImageMatch != nil {
		imageURL := gyazoImageMatch[1]
		return &imageURL, nil
	}

	// Match ASIN pattern and get image URL from the database
	if asinMatch := asinRe.FindStringSubmatch(entry.Body); asinMatch != nil {
		asin := asinMatch[1]
		row, err := w.store.GetAmazonImageUrlByAsin(ctx, asin)
		if err != nil {
			return nil, fmt.Errorf("failed to get Amazon image URL for ASIN %s: %w", asin, err)
		}
		return &row.String, nil
	}

	return nil, nil
}
