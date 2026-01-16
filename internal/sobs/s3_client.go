package sobs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type SobsClient struct {
	s3Client                *s3.Client
	s3AttachmentsBucketName string
	s3BackupBucketName      string
}

func NewSobsClient(s3AccessKeyId, s3SecretAccessKey, s3Region, s3AttachmentsBucketName, s3BackupBucketName, s3Endpoint string, useSSL bool) (*SobsClient, error) {
	if s3AccessKeyId == "" || s3SecretAccessKey == "" {
		return nil, fmt.Errorf("S3 credentials are not set: access key ID or secret access key is empty")
	}

	slog.Info("Creating S3 client", slog.String("endpoint", s3Endpoint), slog.Bool("useSSL", useSSL))

	protocol := "https"
	if !useSSL {
		protocol = "http"
	}
	endpointURL := fmt.Sprintf("%s://%s", protocol, s3Endpoint)

	s3Client := s3.New(s3.Options{
		Region: s3Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			s3AccessKeyId,
			s3SecretAccessKey,
			"",
		),
		BaseEndpoint: aws.String(endpointURL),
		UsePathStyle: true,
	})

	return &SobsClient{
		s3Client:                s3Client,
		s3AttachmentsBucketName: s3AttachmentsBucketName,
		s3BackupBucketName:      s3BackupBucketName,
	}, nil
}

func (c *SobsClient) PutObjectToAttachmentBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.s3AttachmentsBucketName),
		Key:           aws.String(key),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	})
	if err != nil {
		return fmt.Errorf("failed to put object to attachment bucket %s with key %s: %w", c.s3AttachmentsBucketName, key, err)
	}

	return nil
}

func (c *SobsClient) PutObjectToBackupBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	slog.Info("Uploading file to Sobs", slog.String("bucket", c.s3BackupBucketName), slog.String("key", key))
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.s3BackupBucketName),
		Key:           aws.String(key),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
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
	paginator := s3.NewListObjectsV2Paginator(c.s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(c.s3BackupBucketName),
	})

	deletedCount := 0
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			slog.Error("Error listing objects", slog.Any("error", err))
			return fmt.Errorf("failed to list objects in backup bucket: %w", err)
		}

		for _, object := range page.Contents {
			key := aws.ToString(object.Key)
			slog.Debug("Checking object for deletion",
				slog.String("bucket", c.s3BackupBucketName),
				slog.String("key", key))

			// Only process *.sql.enc files
			if !strings.HasSuffix(key, ".sql.enc") {
				continue
			}

			// Check if object is older than cutoff time
			if object.LastModified != nil && object.LastModified.Before(cutoffTime) {
				slog.Info("Deleting old backup file",
					slog.String("key", key),
					slog.Time("lastModified", *object.LastModified),
					slog.Time("cutoffTime", cutoffTime))

				_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
					Bucket: aws.String(c.s3BackupBucketName),
					Key:    aws.String(key),
				})
				if err != nil {
					slog.Error("Failed to delete old backup",
						slog.String("key", key),
						slog.Any("error", err))
					continue
				}
				deletedCount++
			}
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
