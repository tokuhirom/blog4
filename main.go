package main

import (
	"database/sql"
	"embed"
	"github.com/tokuhirom/blog3/db/mariadb"
	"github.com/tokuhirom/blog3/middleware"
	"html/template"
	"log"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

//go:embed templates/*
var templateFS embed.FS

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

type EntryViewData struct {
	Path        string
	Title       string
	PublishedAt string
}

func renderTopPage(w http.ResponseWriter, r *http.Request, queries *mariadb.Queries) {
	// Parse and execute the template
	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	entries, err := queries.SearchEntries(r.Context(), 0)
	if err != nil {
		slog.Info("failed to search entries: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Prepare data for the template
	var viewData []EntryViewData
	for _, entry := range entries {
		// Format the PublishedAt date
		var formattedDate string
		if entry.PublishedAt.Valid {
			formattedDate = entry.PublishedAt.Time.Format("2006-01-02(Mon)")
		} else {
			log.Printf("published_at is invalid: path=%s, published_at=%v", entry.Path, entry.PublishedAt)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		viewData = append(viewData, EntryViewData{
			Path:        entry.Path,
			Title:       entry.Title,
			PublishedAt: formattedDate,
		})
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, viewData)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func renderEntryPage(w http.ResponseWriter, r *http.Request, queries *mariadb.Queries) {
	extractedPath := strings.TrimPrefix(r.URL.Path, "/entry/")

	log.Printf("path: %s", extractedPath)
	entry, err := queries.GetEntryByPath(r.Context(), extractedPath)
	if err != nil {
		slog.Info("failed to get entry by path", err)
		http.NotFound(w, r)
		return
	}

	// Data to pass to the template
	var formattedDate string
	if entry.PublishedAt.Valid {
		formattedDate = entry.PublishedAt.Time.Format("2006-01-01(Mon) 15:04")
	} else {
		log.Printf("published_at is invalid: path=%s, published_at=%v", entry.Path, entry.PublishedAt)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := struct {
		Title       string
		Body        string
		PublishedAt string
	}{
		Title:       entry.Title,
		Body:        entry.Body,
		PublishedAt: formattedDate,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFS(templateFS, "templates/entry.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
		log.Println("top page")
		renderTopPage(writer, request, queries)
	})
	mux.HandleFunc("/entry/", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("entry page")
		renderEntryPage(writer, request, queries)
	})

	loggedMux := middleware.LoggingMiddleware(mux)

	// Start the server
	log.Println("Starting server on http://localhost:8181/")
	err = http.ListenAndServe(":8181", loggedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
