package internal

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	DB, err = sql.Open(
		"postgres",
		"user=postgres password=postgres dbname=pastebin sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}

	log.Println("PostgreSQL connected")
}
