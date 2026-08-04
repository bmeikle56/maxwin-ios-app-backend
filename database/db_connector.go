package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// singleton Postgres instance — call GetDB() after ConnectDB().
var db *sql.DB

func GetDB() *sql.DB {
	return db
}

// ConnectDB opens and pings the Postgres connection using DB_URL.
// Call this when USE_MOCK_DB is false.
func ConnectDB() error {
	url := os.Getenv("DB_URL")
	if url == "" {
		return fmt.Errorf("DB_URL is not set")
	}

	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		return fmt.Errorf("sql.Open error: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return fmt.Errorf("db.Ping error: %w", err)
	}

	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
