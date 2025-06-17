package sobs

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log/slog"
)

type SobsClient struct {
	minioClient             *minio.Client
	s3AttachmentsBucketName string
	s3BackupBucketName      string
}

func NewSobsClient(s3AccessKeyId, s3SecretAccessKey, s3Region, s3AttachmentsBucketName, s3BackupBucketName, s3Endpoint string) (*SobsClient, error) {
	if s3AccessKeyId == "" || s3SecretAccessKey == "" {
		return nil, fmt.Errorf("S3 credentials are not set: access key ID or secret access key is empty")
	}

	slog.Info("Creating S3 client", slog.String("endpoint", s3Endpoint))

	// Initialize minio client object.
	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyId, s3SecretAccessKey, ""),
		Secure: true,
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
