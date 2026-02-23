package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func openDB(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_busy_timeout=5000", path)

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initSchema(db *sql.DB) error {
	schema := `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  username      TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at    DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at    DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE IF NOT EXISTS urls (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id      INTEGER NOT NULL,
  code         TEXT NOT NULL UNIQUE,
  original_url TEXT NOT NULL,
  created_at   DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at   DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),

  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_urls_user_created ON urls(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_urls_user_code ON urls(user_id, code);

-- Faster lookups and sorting
CREATE UNIQUE INDEX IF NOT EXISTS idx_urls_code ON urls(code);
CREATE INDEX IF NOT EXISTS idx_urls_user_updated ON urls(user_id, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_urls_user_original ON urls(user_id, original_url);

`
	_, err := db.Exec(schema)
	return err
}