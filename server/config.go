package server

import (
	"strings"
)

type Config struct {
	LocalDev bool `env:"LOCAL_DEV" envDefault:"false"`

	Port int `env:"BLOG_PORT" envDefault:"9191"`

	DBUser     string `env:"DATABASE_USER"`
	DBPassword string `env:"DATABASE_PASSWORD"`
	DBHostname string `env:"DATABASE_HOST"`
	DBPort     int    `env:"DATABASE_PORT" envDefault:"3306"`
	DBName     string `env:"DATABASE_DB"   envDefault:"blog3"`

	AdminUser     string `env:"ADMIN_USER"   envDefault:"admin"`
	AdminPassword string `env:"ADMIN_PW"   envDefault:"admin"`

	HubUrls string `env:"HUB_URLS"`

	AmazonPaapi5AccessKey string `env:"AMAZON_PAAPI5_ACCESS_KEY"`
	AmazonPaapi5SecretKey string `env:"AMAZON_PAAPI5_SECRET_KEY"`

	S3AccessKeyId           string `env:"S3_ACCESS_KEY_ID"`
	S3SecretAccessKey       string `env:"S3_SECRET_ACCESS_KEY"`
	S3Region                string `env:"S3_REGION" envDefault:"jp-north-1"`
	S3AttachmentsBucketName string `env:"S3_ATTACHMENTS_BUCKET_NAME" envDefault:"blog3-attachments"`
	S3BackupBucketName      string `env:"S3_BACKUP_BUCKET_NAME" envDefault:"blog3-backup"`
	S3Endpoint              string `env:"S3_ENDPOINT" envDefault:"s3.isk01.sakurastorage.jp"`

	BackupEncryptionKey string `env:"BACKUP_ENCRYPTION_KEY"`

	// 9*60*60=32400 is JST
	TimeZoneOffset int `env:"TIMEZONE_OFFSET" envDefault:"32400"`
}

func (c *Config) GetHubUrls() []string {
	if c.HubUrls != "" {
		return strings.Split(c.HubUrls, ",")
	}
	return []string{
		"https://pubsubhubbub.appspot.com/",
		"https://pubsubhubbub.superfeedr.com/",
	}
}
