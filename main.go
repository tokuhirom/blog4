package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tokuhirom/blog3/db/mariadb"
	"github.com/tokuhirom/blog3/server"
	"github.com/tokuhirom/blog3/server/admin"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

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

	queries := mariadb.New(sqlDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", server.Router(queries))
	r.Mount("/admin", admin.Router(cfg))

	// Start the server
	log.Println("Starting server on http://localhost:8181/")
	err = http.ListenAndServe(":8181", r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
