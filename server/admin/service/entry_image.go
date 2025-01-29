package service

import (
	"fmt"
	"log"
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
			log.Printf("failed to process entry %s: %v", entry.Path, err)
		}
	}

	return nil
}

// ProcessEntry processes a single entry.
func (s *EntryImageService) ProcessEntry(entry Entry) error {
	log.Printf("process entry image: %s", entry.Path)

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
		log.Printf("imageTagMatch: %s", imageTagMatch[1])
		return imageTagMatch[1], nil
	}

	basicImage := regexp.MustCompile(`!\[.*?\]\((https?:\/\/[^\s)]+)\)`).FindStringSubmatch(entry.Body)
	if len(basicImage) > 1 {
		log.Printf("basicImage: %s", basicImage[1])
		return basicImage[1], nil
	}

	gyazoImage := regexp.MustCompile(`\[!\[.*?\]\((https?:\/\/[^\s)]+)\)\]\((.*?)\)`).FindStringSubmatch(entry.Body)
	if len(gyazoImage) > 1 {
		log.Printf("gyazoImage: %s", gyazoImage[1])
		return gyazoImage[1], nil
	}

	asin := regexp.MustCompile(`\[asin:([A-Z0-9]+):detail\]`).FindStringSubmatch(entry.Body)
	if len(asin) > 1 {
		imageUrl, err := s.amazonRepository.GetAsinImage(asin[1])
		if err != nil {
			return "", fmt.Errorf("failed to get ASIN image: %w", err)
		}
		log.Printf("asin: asin=%s imageUrl=%s", asin[1], imageUrl)
		return imageUrl, nil
	}

	return "", nil
}
