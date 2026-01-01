package internal

import (
	"database/sql"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("environment variable DATABASE_URL is not set")
	}

	if !strings.Contains(dbURL, "sslmode=") {
		if os.Getenv("ENV") == "LOCAL" {
			dbURL += "?sslmode=disable"
		} else {
			dbURL += "?sslmode=require"
		}
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("sql.Open error: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	log.Println("PostgreSQL connected")
}
