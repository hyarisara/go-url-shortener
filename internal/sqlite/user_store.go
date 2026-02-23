package sqlite

import (
	"database/sql"
	"errors"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(dbPath string) *UserStore {
	db, err := openDB(dbPath)
	if err != nil {
		panic(err)
	}
	return &UserStore{db: db}
}

// UserExists checks if username already exists
func (s *UserStore) UserExists(username string) bool {
	var id int64
	err := s.db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	return err == nil
}

// SaveUser inserts a new user (username must be unique)
func (s *UserStore) SaveUser(username, passwordHash string) error {
	_, err := s.db.Exec(`
		INSERT INTO users (username, password_hash, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, username, passwordHash)
	return err
}

// GetUser returns password hash for login validation
func (s *UserStore) GetUser(username string) (string, error) {
	var hash string
	err := s.db.QueryRow(`SELECT password_hash FROM users WHERE username = ?`, username).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("user not found")
		}
		return "", err
	}
	return hash, nil
}

// GetUserID is needed by URL store to link URLs with user_id FK
func (s *UserStore) GetUserID(username string) (int64, error) {
	var id int64
	err := s.db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("user not found")
		}
		return 0, err
	}
	return id, nil
}