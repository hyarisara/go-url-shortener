package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteUserStore struct {
	db *sql.DB
}

func NewUserStore(file string) *SQLiteUserStore {
	db, _ := sql.Open("sqlite3", file)
	db.Exec("CREATE TABLE IF NOT EXISTS users(username TEXT PRIMARY KEY, hash TEXT)")
	return &SQLiteUserStore{db: db}
}

func (s *SQLiteUserStore) SaveUser(username, hash string) error {
	_, err := s.db.Exec("INSERT INTO users(username, hash) VALUES(?, ?)", username, hash)
	return err
}

func (s *SQLiteUserStore) GetUser(username string) (string, error) {
	var hash string
	row := s.db.QueryRow("SELECT hash FROM users WHERE username=?", username)
	err := row.Scan(&hash)
	if err != nil {
		return "", errors.New("user not found")
	}
	return hash, nil
}

func (s *SQLiteUserStore) UserExists(username string) bool {
	_, err := s.GetUser(username)
	return err == nil
}
