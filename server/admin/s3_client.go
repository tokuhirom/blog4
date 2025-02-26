package admin

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"log"
	"os"
)

// See https://cloud.sakura.ad.jp/news/2025/02/04/objectstorage_defectversion/

type S3Client struct {
	s3Client                *s3.Client
	s3AttachmentsBucketName string
}

func NewS3Client(s3AccessKeyId, s3SecretAccessKey, s3Region, s3AttachmentsBucketName, s3Endpoint string) *S3Client {
	if s3AccessKeyId == "" || s3SecretAccessKey == "" {
		log.Fatal("S3 credentials are not set")
	}

	requestChecksumCalculation := os.Getenv("AWS_REQUEST_CHECKSUM_CALCULATION")
	log.Printf("AWS_REQUEST_CHECKSUM_CALCULATION: '%s'", requestChecksumCalculation)

	// AWS_RESPONSE_CHECKSUM_VALIDATION
	responseChecksumValidation := os.Getenv("AWS_RESPONSE_CHECKSUM_VALIDATION")
	log.Printf("AWS_RESPONSE_CHECKSUM_VALIDATION: '%s'", responseChecksumValidation)

	// カスタムエンドポイントの設定
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(s3Region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(
			func(ctx context.Context) (aws.Credentials, error) {
				return aws.Credentials{
					AccessKeyID:     s3AccessKeyId,
					SecretAccessKey: s3SecretAccessKey,
				}, nil
			},
		)),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// S3クライアントの初期化
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s3Endpoint)
	})

	return &S3Client{
		s3Client:                s3Client,
		s3AttachmentsBucketName: s3AttachmentsBucketName,
	}
}

func (c *S3Client) PutObjectToAttachmentBucket(ctx context.Context, key string, contentType string, contentLength int64, body io.Reader) error {
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:            aws.String(c.s3AttachmentsBucketName),
		Key:               aws.String(key),
		ContentType:       aws.String(contentType),
		ContentLength:     &contentLength,
		Body:              body,
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
	})
	if err != nil {
		return err
	}

	return nil
}
