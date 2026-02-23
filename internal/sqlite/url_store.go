package sqlite

import (
	"database/sql"
	"errors"
	"strings"
)

type URLStore struct {
	db *sql.DB
}

func NewURLStore(dbPath string) *URLStore {
	db, err := openDB(dbPath)
	if err != nil {
		panic(err)
	}
	return &URLStore{db: db}
}

// SaveURL expects key format: "username::code"
func (s *URLStore) SaveURL(key, originalURL string) error {
	username, code, err := splitKey(key)
	if err != nil {
		return err
	}

	// Lookup user_id
	var userID int64
	err = s.db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO urls (user_id, code, original_url, created_at, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, code, originalURL)
	return err
}

// GetURLByCode finds original_url by short code
func (s *URLStore) GetURLByCode(code string) (string, error) {
	var original string
	err := s.db.QueryRow(`SELECT original_url FROM urls WHERE code = ?`, code).Scan(&original)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("code not found")
		}
		return "", err
	}
	return original, nil
}

// ListByUser returns all URLs for a user as map[code]original_url
func (s *URLStore) ListByUser(username string) (map[string]string, error) {
	var userID int64
	err := s.db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]string{}, nil
		}
		return nil, err
	}

	rows, err := s.db.Query(`
		SELECT code, original_url
		FROM urls
		WHERE user_id = ?
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string)
	for rows.Next() {
		var code, url string
		if err := rows.Scan(&code, &url); err != nil {
			return nil, err
		}
		out[code] = url
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteURL expects key format: "username::code"
func (s *URLStore) DeleteURL(key string) error {
	username, code, err := splitKey(key)
	if err != nil {
		return err
	}

	var userID int64
	err = s.db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	res, err := s.db.Exec(`DELETE FROM urls WHERE user_id = ? AND code = ?`, userID, code)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("url not found")
	}
	return nil
}

func splitKey(key string) (username, code string, err error) {
	parts := strings.Split(key, "::")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", errors.New("invalid key format (expected username::code)")
	}
	return parts[0], parts[1], nil
}
func (s *URLStore) ListByUserPaged(username string, q string, sort string, page int, pageSize int) (map[string]string, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	// Lookup user_id
	var userID int64
	err := s.db.QueryRow(`SELECT id FROM users WHERE username = ?`, username).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return map[string]string{}, nil
		}
		return nil, err
	}

	// Safe sort options
	orderBy := "created_at DESC"
	if sort == "updated" {
		orderBy = "updated_at DESC"
	}

	// Search (basic)
	// NOTE: LIKE is not perfect; Step 6 can add FTS later.
	query := `
		SELECT code, original_url
		FROM urls
		WHERE user_id = ?
	`
	args := []any{userID}

	if q != "" {
		query += ` AND (code LIKE ? OR original_url LIKE ?)`
		args = append(args, "%"+q+"%", "%"+q+"%")
	}

	query += ` ORDER BY ` + orderBy + ` LIMIT ? OFFSET ?`
	args = append(args, pageSize, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string)
	for rows.Next() {
		var code, url string
		if err := rows.Scan(&code, &url); err != nil {
			return nil, err
		}
		out[code] = url
	}
	return out, rows.Err()
}