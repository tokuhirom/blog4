package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/router"
	"github.com/tokuhirom/blog4/server/sobs"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --config ../../ogen.yml --target ../../server/admin/openapi -package openapi --clean ../../typespec/tsp-output/@typespec/openapi3/openapi.yaml

func main() {
	if err := DoMain(); err != nil {
		slog.Error("failed to start server", slog.Any("error", err))
		os.Exit(1)
	}
	os.Exit(0)
}

func DoMain() error {
	cfg, err := env.ParseAs[server.Config]()
	if err != nil {
		return fmt.Errorf("failed to parse Config: %w", err)
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
		return fmt.Errorf("failed to open DB connection: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	// Use SSL for S3 connections unless in local development mode
	useSSL := !cfg.LocalDev
	sobsClient, err := sobs.NewSobsClient(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, cfg.S3Region, cfg.S3AttachmentsBucketName, cfg.S3BackupBucketName, cfg.S3Endpoint, useSSL)
	if err != nil {
		return fmt.Errorf("failed to create SobsClient: %w", err)
	}

	r, err := router.BuildRouter(cfg, sqlDB, sobsClient)
	if err != nil {
		return fmt.Errorf("failed to build router: %w", err)
	}

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
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
