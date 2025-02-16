package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tokuhirom/blog4/db/public/publicdb"
	"github.com/tokuhirom/blog4/server"
	"github.com/tokuhirom/blog4/server/admin"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target ./server/admin/openapi -package openapi --clean openapi.yml

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/admin", admin.Router(cfg, sqlDB))
	r.Mount("/", server.Router(publicQueries))

	// Start the server
	log.Println("Starting server on http://localhost:8181/")
	err = http.ListenAndServe(":8181", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
