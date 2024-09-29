package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("sqlite3", "puny-url.db")
	if err != nil {
		log.Fatal(err)
	}

	runMigrations()
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func runMigrations() {
	migrationSQL, err := os.ReadFile("./db/migrations.sql")
	if err != nil {
		log.Fatalf("Error reading migrations: %v", err)
	}

	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func StoreURL(id string, longURL string) error {
	_, err := db.Exec("INSERT INTO urls (id, long_url) VALUES (?, ?)", id, longURL)
	return err
}

func GetLong(id string) (string, error) {
	var longURL string
	err := db.QueryRow("SELECT long_url FROM urls WHERE id = ?", id).Scan(&longURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return longURL, nil
}

func GetShortIDByLongURL(longURL string) (string, error) {
	var shortID string
	err := db.QueryRow("SELECT id FROM urls WHERE long_url = ?", longURL).Scan(&shortID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return shortID, err
}
