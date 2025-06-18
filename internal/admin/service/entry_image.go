package service

import (
	"fmt"
	"log/slog"
	"regexp"
)

// Entry represents the structure of an entry.
type Entry struct {
	Path string
	Body string
}

// EntryImageRepository defines the interface for entry image operations.
type EntryImageRepository interface {
	GetNotProcessedEntries() ([]Entry, error)
	InsertEntryImage(path, image string) error
}

// AmazonRepository defines the interface for Amazon operations.
type AmazonRepository interface {
	GetAsinImage(asin string) (string, error)
}

// EntryImageService provides methods to process entry images.
type EntryImageService struct {
	entryImageRepository EntryImageRepository
	amazonRepository     AmazonRepository
}

// NewEntryImageService creates a new EntryImageService.
func NewEntryImageService(entryImageRepo EntryImageRepository, amazonRepo AmazonRepository) *EntryImageService {
	return &EntryImageService{
		entryImageRepository: entryImageRepo,
		amazonRepository:     amazonRepo,
	}
}

// ProcessEntryImages processes all not processed entries.
func (s *EntryImageService) ProcessEntryImages() error {
	entries, err := s.entryImageRepository.GetNotProcessedEntries()
	if err != nil {
		return fmt.Errorf("failed to get not processed entries: %w", err)
	}

	for _, entry := range entries {
		if err := s.ProcessEntry(entry); err != nil {
			slog.Error("failed to process entry",
				slog.String("path", entry.Path),
				slog.Any("error", err))
		}
	}

	return nil
}

// ProcessEntry processes a single entry.
func (s *EntryImageService) ProcessEntry(entry Entry) error {
	slog.Info("processing entry image",
		slog.String("path", entry.Path))

	image, err := s.GetImageFromEntry(entry)
	if err != nil {
		return fmt.Errorf("failed to get image from entry: %w", err)
	}

	if image == "" && regexp.MustCompile(`\[asin:`).MatchString(entry.Body) {
		// ASIN processing is delayed; skip it.
		return nil
	}

	return s.entryImageRepository.InsertEntryImage(entry.Path, image)
}

// GetImageFromEntry extracts the image URL from the entry body.
func (s *EntryImageService) GetImageFromEntry(entry Entry) (string, error) {
	imageTagMatch := regexp.MustCompile(`<img[^>]*src=['"]?(https?:\/\/[^\s)]+)\.(?:jpg|png|gif)['"]?`).FindStringSubmatch(entry.Body)
	if len(imageTagMatch) > 1 {
		slog.Debug("found image tag match", slog.String("url", imageTagMatch[1]))
		return imageTagMatch[1], nil
	}

	basicImage := regexp.MustCompile(`!\[.*?\]\((https?:\/\/[^\s)]+)\)`).FindStringSubmatch(entry.Body)
	if len(basicImage) > 1 {
		slog.Debug("found basic image", slog.String("url", basicImage[1]))
		return basicImage[1], nil
	}

	gyazoImage := regexp.MustCompile(`\[!\[.*?\]\((https?:\/\/[^\s)]+)\)\]\((.*?)\)`).FindStringSubmatch(entry.Body)
	if len(gyazoImage) > 1 {
		slog.Debug("found gyazo image", slog.String("url", gyazoImage[1]))
		return gyazoImage[1], nil
	}

	asin := regexp.MustCompile(`\[asin:([A-Z0-9]+):detail\]`).FindStringSubmatch(entry.Body)
	if len(asin) > 1 {
		imageUrl, err := s.amazonRepository.GetAsinImage(asin[1])
		if err != nil {
			return "", fmt.Errorf("failed to get ASIN image: %w", err)
		}
		slog.Debug("found ASIN image", slog.String("asin", asin[1]), slog.String("image_url", imageUrl))
		return imageUrl, nil
	}

	return "", nil
}
