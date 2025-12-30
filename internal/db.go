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
	DB, err = sql.Open("postgres", "user=postgres password=mypassword dbname=pastebin sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to PostgreSQL successfully")
}
