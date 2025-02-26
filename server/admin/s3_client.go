package admin

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
)

type S3Client struct {
	minioClient             *minio.Client
	s3AttachmentsBucketName string
}

func NewS3Client(s3AccessKeyId, s3SecretAccessKey, s3Region, s3AttachmentsBucketName, s3Endpoint string) *S3Client {
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

	return &S3Client{
		minioClient:             minioClient,
		s3AttachmentsBucketName: s3AttachmentsBucketName,
	}
}

func (c *S3Client) PutObjectToAttachmentBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	_, err := c.minioClient.PutObject(ctx, c.s3AttachmentsBucketName, key, body, contentLength, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}
