package main

import (
	"database/sql"
	"github.com/tokuhirom/blog3/db/mariadb"
	"github.com/tokuhirom/blog3/middleware"
	"github.com/tokuhirom/blog3/server"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type config struct {
	Port       int    `env:"BLOG3_PORT" envDefault:"9191"`
	DBUser     string `env:"DATABASE_USER"`
	DBPassword string `env:"DATABASE_PASSWORD"`
	DBHostname string `env:"DATABASE_HOST"`
	DBPort     int    `env:"DATABASE_PORT" envDefault:"3306"`
	DBName     string `env:"DATABASE_DB"   envDefault:"blog3"`
	// 9*60*60=32400 is JST
	TimeZoneOffset int `env:"TIMEZONE_OFFSET" envDefault:"32400"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	cfg, err := env.ParseAs[config]()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		server.RenderTopPage(writer, request, queries)
	})
	mux.HandleFunc("/entry/", func(writer http.ResponseWriter, request *http.Request) {
		server.RenderEntryPage(writer, request, queries)
	})

	loggedMux := middleware.LoggingMiddleware(mux)

	// Start the server
	log.Println("Starting server on http://localhost:8181/")
	err = http.ListenAndServe(":8181", loggedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
