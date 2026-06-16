package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

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

var db *sql.DB

func Init(dbFile string) error {
	var install bool

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

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}
	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
