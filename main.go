package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/admin"
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

	publicQueries := publicdb.New(sqlDB)
	sobsClient := sobs.NewSobsClient(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, cfg.S3Region, cfg.S3AttachmentsBucketName, cfg.S3BackupBucketName, cfg.S3Endpoint)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/admin", admin.Router(cfg, sqlDB, sobsClient))
	r.Mount("/", server.Router(publicQueries))
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})
	r.Get("/git_hash", func(w http.ResponseWriter, r *http.Request) {
		gitHash := os.Getenv("GIT_HASH")
		if gitHash == "" {
			http.Error(w, "GIT_HASH not set", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(gitHash))
		if err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})

	go (func() {
		log.Printf("Starting backup process...")
		server.StartBackup(cfg.BackupEncryptionKey, sobsClient)
	})()

	// apprun が今 min-scale 0 なので､寝ないように自分で自分を起こし続ける
	go (func() {
		url := "https://app-ceb34d94-8ffd-4ade-ba21-1dd9b124f836.ingress.apprun.sakura.ne.jp/healthz"
		for {
			time.Sleep(10 * time.Second)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("failed to request %s: %v", url, err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("unexpected status code: %d", resp.StatusCode)
			}
			_ = resp.Body.Close()
		}
	})()

	// Start the server
	log.Println("Starting server on http://localhost:8181/")
	err = http.ListenAndServe(":8181", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
