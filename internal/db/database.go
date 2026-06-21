package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// schema contains SQL statements to create database tables and indexes
const (
	schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(128) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
    );

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
    `
)

// db is the global database connection
var db *sql.DB

// Init initializes database connection and creates schema if needed
func Init(dbFile string) error {
	var install bool

	// Check if database file exists
	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("failed to check database file: %w", err)
		}
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create schema if this is a new database
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}
	return nil
}

// Close closes the database connection
func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
