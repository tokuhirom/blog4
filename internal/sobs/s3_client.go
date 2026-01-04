package sobs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type SobsClient struct {
	minioClient             *minio.Client
	s3AttachmentsBucketName string
	s3BackupBucketName      string
}

func NewSobsClient(s3AccessKeyId, s3SecretAccessKey, s3Region, s3AttachmentsBucketName, s3BackupBucketName, s3Endpoint string, useSSL bool) (*SobsClient, error) {
	if s3AccessKeyId == "" || s3SecretAccessKey == "" {
		return nil, fmt.Errorf("S3 credentials are not set: access key ID or secret access key is empty")
	}

	slog.Info("Creating S3 client", slog.String("endpoint", s3Endpoint), slog.Bool("useSSL", useSSL))

	// Initialize minio client object.
	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyId, s3SecretAccessKey, ""),
		Secure: useSSL,
		Region: s3Region,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize minio client for endpoint %s: %w", s3Endpoint, err)
	}

	return &SobsClient{
		minioClient:             minioClient,
		s3AttachmentsBucketName: s3AttachmentsBucketName,
		s3BackupBucketName:      s3BackupBucketName,
	}, nil
}

func (c *SobsClient) PutObjectToAttachmentBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	_, err := c.minioClient.PutObject(ctx, c.s3AttachmentsBucketName, key, body, contentLength, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to put object to attachment bucket %s with key %s: %w", c.s3AttachmentsBucketName, key, err)
	}

	return nil
}

func (c *SobsClient) PutObjectToBackupBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	slog.Info("Uploading file to Sobs", slog.String("bucket", c.s3BackupBucketName), slog.String("key", key))
	_, err := c.minioClient.PutObject(ctx, c.s3BackupBucketName, key, body, contentLength, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to put object to backup bucket %s with key %s: %w", c.s3BackupBucketName, key, err)
	}

	return nil
}

// DeleteOldBackups deletes backup files older than the specified number of days
// Only files matching the *.sql.enc pattern are deleted
func (c *SobsClient) DeleteOldBackups(ctx context.Context, daysToKeep int) error {
	slog.Info("Checking for old backups to delete",
		slog.String("bucket", c.s3BackupBucketName),
		slog.Int("daysToKeep", daysToKeep))

	// Calculate cutoff time
	cutoffTime := time.Now().AddDate(0, 0, -daysToKeep)

	// List all objects in backup bucket
	objectCh := c.minioClient.ListObjects(ctx, c.s3BackupBucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	deletedCount := 0
	for object := range objectCh {
		if object.Err != nil {
			slog.Error("Error listing objects", slog.Any("error", object.Err))
			continue
		}

		// Only process *.sql.enc files
		if !strings.HasSuffix(object.Key, ".sql.enc") {
			continue
		}

		// Check if object is older than cutoff time
		if object.LastModified.Before(cutoffTime) {
			slog.Info("Deleting old backup file",
				slog.String("key", object.Key),
				slog.Time("lastModified", object.LastModified),
				slog.Time("cutoffTime", cutoffTime))

			err := c.minioClient.RemoveObject(ctx, c.s3BackupBucketName, object.Key, minio.RemoveObjectOptions{})
			if err != nil {
				slog.Error("Failed to delete old backup",
					slog.String("key", object.Key),
					slog.Any("error", err))
				continue
			}
			deletedCount++
		}
	}

	if deletedCount > 0 {
		slog.Info("Old backups deleted",
			slog.Int("count", deletedCount),
			slog.Int("daysToKeep", daysToKeep))
	} else {
		slog.Info("No old backups to delete")
	}

	return nil
}
