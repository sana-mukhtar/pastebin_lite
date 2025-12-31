package internal

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func InitDB() {
	var err error
	// minimal DSN: user, password, dbname, sslmode
	DB, err = sql.Open("postgres", "user=postgres password=postgres dbname=pastebin sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL successfully")
}

func AutoMigrate() {
	_, err := DB.Exec(`
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS pastes (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		content TEXT NOT NULL,
		ttl_seconds INT,
		max_views INT,
		created_at TIMESTAMP DEFAULT NOW(),
		expires_at TIMESTAMP,
		views INT DEFAULT 0
	);
	`)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("Table 'pastes' ensured")
}
