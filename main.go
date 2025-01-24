package main

import (
	"database/sql"
	"embed"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

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
}

// Entry represents a blog article.
type Entry struct {
	Title   string
	Content string
}

// Dummy data for demonstration. In real use, retrieve from a database.
var articles = map[string]Entry{
	"example/path": {Title: "Example Title", Content: "This is an example article content."},
}

func renderTopPage(w http.ResponseWriter, r *http.Request) {
	// Parse and execute the template
	tmpl, err := template.ParseFS(templateFS, "templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, articles)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Render the article page
func renderArticlePage(w http.ResponseWriter, r *http.Request) {
	// Extract the PATH from the URL
	path := strings.TrimPrefix(r.URL.Path, "/entry/")
	if path == "" {
		http.NotFound(w, r)
		return
	}

	// Retrieve the article from the database or dummy data
	article, exists := articles[path]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// Parse and execute the template
	tmpl, err := template.ParseFS(templateFS, "templates/entry.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, article)
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
	}
	sqlDB, err := sql.Open("mysql", mysqlConfig.FormatDSN())
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("failed to ping DB: %v", err)
	}

	http.HandleFunc("/", renderTopPage)
	http.HandleFunc("/entry/", renderArticlePage)

	// Start the server
	log.Println("Starting server on http://localhost:8181/")
	err = http.ListenAndServe(":8181", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
