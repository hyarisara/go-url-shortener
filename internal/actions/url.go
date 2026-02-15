package actions

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"go-url-shortener/internal/store"
)

type URLService struct {
	store store.URLStore
}

func NewURLService(s store.URLStore) *URLService {
	return &URLService{store: s}
}

//////////////////////
// Helpers
//////////////////////

func generateCode() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

//////////////////////
// Business Logic
//////////////////////

func (s *URLService) ShortenForUser(username, originalURL, custom string) (string, error) {
	code := custom

	if code == "" {
		var err error
		code, err = generateCode()
		if err != nil {
			return "", err
		}
	}

	key := username + "::" + code

	err := s.store.SaveURL(key, originalURL)
	return code, err
}

func (s *URLService) Expand(code string) (string, error) {
	// we don't know the user → store should find by code
	return s.store.GetURLByCode(code)
}

func (s *URLService) ListForUser(username string) (map[string]string, error) {
	return s.store.ListByUser(username)
}

func (s *URLService) Delete(key string) error {
	if key == "" {
		return errors.New("empty key")
	}
	return s.store.DeleteURL(key)
}
