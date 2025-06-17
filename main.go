package main

import (
	"database/sql"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/router"
	"github.com/tokuhirom/blog4/server/sobs"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target ./server/admin/openapi -package openapi --clean typespec/tsp-output/@typespec/openapi3/openapi.yaml

func main() {
	if _, err := os.Stat(".env"); err == nil {
		slog.Info("loading .env file")
		err := godotenv.Load()
		if err != nil {
			slog.Error("failed to load .env", slog.Any("error", err))
			os.Exit(1)
		}
	} else {
		slog.Info(".env file not found")
	}

	cfg, err := env.ParseAs[server.Config]()
	if err != nil {
		slog.Error("failed to parse Config", slog.Any("error", err))
		os.Exit(1)
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
		slog.Error("failed to open DB", slog.Any("error", err))
		os.Exit(1)
	}
	if err := sqlDB.Ping(); err != nil {
		slog.Error("failed to ping DB", slog.Any("error", err))
		os.Exit(1)
	}

	sobsClient, err := sobs.NewSobsClient(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, cfg.S3Region, cfg.S3AttachmentsBucketName, cfg.S3BackupBucketName, cfg.S3Endpoint)
	if err != nil {
		slog.Error("failed to create S3 client", slog.Any("error", err))
		os.Exit(1)
	}

	r := router.BuildRouter(cfg, sqlDB, sobsClient)

	go (func() {
		slog.Info("Starting backup process")
		server.StartBackup(cfg.BackupEncryptionKey, sobsClient)
	})()

	if cfg.KeepAliveUrl != "" {
		url := cfg.KeepAliveUrl
		go server.KeepAlive(url)
	}

	// Start the server
	slog.Info("Starting server", slog.String("url", "http://localhost:8181/"))
	err = http.ListenAndServe(":8181", r)
	if err != nil {
		slog.Error("Server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
}
