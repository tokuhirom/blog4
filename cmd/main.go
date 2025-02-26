package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/tokuhirom/blog4/server"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type S3Client struct {
	s3Client                *s3.Client
	s3AttachmentsBucketName string
}

func NewS3Client(s3AccessKeyId, s3SecretAccessKey, s3Region, s3AttachmentsBucketName, s3Endpoint string) *S3Client {
	if s3AccessKeyId == "" || s3SecretAccessKey == "" {
		log.Fatal("S3 credentials are not set")
	}

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

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s3Endpoint)
		o.HTTPClient = &http.Client{
			Transport: &debugTransport{transport: http.DefaultTransport},
		}
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

	// test
	object, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.s3AttachmentsBucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	log.Printf("Object: %v, contentLength:%v", *object.ContentLength, contentLength)
	defer object.Body.Close()
	buf, err := io.ReadAll(object.Body)
	if err != nil {
		log.Fatalf("failed to read object: %v", err)
	}
	log.Printf("length of buf: %v", len(buf))
	// write buf to /tmp/xxx.png
	err = os.WriteFile("/tmp/xxx.png", buf, 0644)
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	return nil
}

type debugTransport struct {
	transport http.RoundTripper
}

func (d *debugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	dump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("Request:\n%s\n", dump)

	resp, err := d.transport.RoundTrip(req)

	if err == nil {
		dump, _ = httputil.DumpResponse(resp, true)
		fmt.Printf("Response:\n%s\n", dump)
	}

	return resp, err
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	cfg, err := env.ParseAs[server.Config]()
	if err != nil {
		log.Fatalf("failed to parse Config: %v", err)
	}

	s3Client := NewS3Client(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, cfg.S3Region, cfg.S3AttachmentsBucketName, cfg.S3Endpoint)
	object, err := s3Client.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:            aws.String(cfg.S3AttachmentsBucketName),
		Key:               aws.String("test.png"),
		Body:              strings.NewReader("test12345"),
		ContentType:       aws.String("text/plain"),
		ChecksumAlgorithm: types.ChecksumAlgorithmCrc32,
	})
	if err != nil {
		log.Fatalf("failed to put object: %v", err)
	}

	getObject, err := s3Client.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(cfg.S3AttachmentsBucketName),
		Key:    aws.String("test.png"),
	})
	if err != nil {
		log.Fatalf("failed to get object: %v", err)
	}
	all, err := io.ReadAll(getObject.Body)
	if err != nil {
		log.Fatalf("failed to read object: %v", err)
	}
	log.Printf("Object: %v", url.QueryEscape(string(all)))

	log.Printf("Object: %v", object)
}
