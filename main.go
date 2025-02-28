package main

import (
	"database/sql"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/router"
	"github.com/tokuhirom/blog4/server/sobs"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target ./server/admin/openapi -package openapi --clean openapi.yml

func main() {
	if _, err := os.Stat(".env"); err == nil {
		log.Printf("loading .env file")
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("failed to load .env: %v", err)
		}
	} else {
		log.Printf(".env file not found")
	}

	cfg, err := env.ParseAs[server.Config]()
	if err != nil {
		log.Fatalf("failed to parse Config: %v", err)
	}

	mysqlConfig := mysql.Config{
		User:                 cfg.DBUser,
		Passwd:               cfg.DBPassword,
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(cfg.DBHostname, strconv.Itoa(cfg.DBPort)),
		DBName:               cfg.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
		Loc:                  time.FixedZone("Asia/Tokyo", cfg.TimeZoneOffset), // Set time zone to JST
	}
	sqlDB, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	sobsClient := sobs.NewSobsClient(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, cfg.S3Region, cfg.S3AttachmentsBucketName, cfg.S3BackupBucketName, cfg.S3Endpoint)

	r := router.BuildRouter(cfg, sqlDB, sobsClient)

	go (func() {
		log.Printf("Starting backup process...")
		server.StartBackup(cfg.BackupEncryptionKey, sobsClient)
	})()

	if cfg.KeepAliveUrl != "" {
		url := cfg.KeepAliveUrl
		go server.KeepAlive(url)
	}

	// Start the server
	log.Println("Starting server on http://localhost:8181/")
	err = http.ListenAndServe(":8181", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
