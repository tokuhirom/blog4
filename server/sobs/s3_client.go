package sobs

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
)

type SobsClient struct {
	minioClient             *minio.Client
	s3AttachmentsBucketName string
	s3BackupBucketName      string
}

func NewSobsClient(s3AccessKeyId, s3SecretAccessKey, s3Region, s3AttachmentsBucketName, s3BackupBucketName, s3Endpoint string) *SobsClient {
	if s3AccessKeyId == "" || s3SecretAccessKey == "" {
		log.Fatal("S3 credentials are not set")
	}

	log.Printf("Creating S3 client for %s", s3Endpoint)

	// Initialize minio client object.
	minioClient, err := minio.New(s3Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(s3AccessKeyId, s3SecretAccessKey, ""),
		Secure: true,
		Region: s3Region,
	})
	if err != nil {
		log.Fatalf("unable to initialize minio client: %v", err)
	}

	return &SobsClient{
		minioClient:             minioClient,
		s3AttachmentsBucketName: s3AttachmentsBucketName,
		s3BackupBucketName:      s3BackupBucketName,
	}
}

func (c *SobsClient) PutObjectToAttachmentBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	_, err := c.minioClient.PutObject(ctx, c.s3AttachmentsBucketName, key, body, contentLength, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *SobsClient) PutObjectToBackupBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	log.Printf("Uploading file to Sobs: bucket=%s key=%s", c.s3BackupBucketName, key)
	_, err := c.minioClient.PutObject(ctx, c.s3BackupBucketName, key, body, contentLength, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}
