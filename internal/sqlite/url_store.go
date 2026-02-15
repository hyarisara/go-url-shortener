package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type URLStore struct {
	db *sql.DB
}

// Constructor
func NewURLStore(path string) *URLStore {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}

	// Create table if not exists
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		key TEXT PRIMARY KEY,
		url TEXT NOT NULL
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	return &URLStore{db: db}
}

// SaveURL implements store.URLStore
func (s *URLStore) SaveURL(key, url string) error {
	_, err := s.db.Exec(`INSERT OR REPLACE INTO urls(key, url) VALUES(?, ?)`, key, url)
	return err
}

// GetURLByCode implements store.URLStore
func (s *URLStore) GetURLByCode(code string) (string, error) {
	row := s.db.QueryRow(`SELECT url FROM urls WHERE key LIKE ?`, "%::"+code)
	var url string
	err := row.Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("code not found")
		}
		return "", err
	}
	return url, nil
}

// ListByUser implements store.URLStore
func (s *URLStore) ListByUser(username string) (map[string]string, error) {
	rows, err := s.db.Query(`SELECT key, url FROM urls WHERE key LIKE ?`, username+"::%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, url string
		rows.Scan(&key, &url)
		// remove username prefix
		code := key[len(username+"::"):]
		result[code] = url
	}
	return result, nil
}

// DeleteURL implements store.URLStore
func (s *URLStore) DeleteURL(key string) error {
	_, err := s.db.Exec(`DELETE FROM urls WHERE key = ?`, key)
	return err
}
