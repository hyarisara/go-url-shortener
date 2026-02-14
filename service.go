package main

import (
	"crypto/rand"
	"encoding/hex"
)

type URLService struct {
	store UrlStore
}

func NewURLService(s UrlStore) *URLService {
	return &URLService{store: s}
}

func generateCode() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Shorten with optional custom code
func (s *URLService) Shorten(url string, customCode string) (string, error) {
	code := customCode
	var err error

	if code == "" {
		code, err = generateCode()
		if err != nil {
			return "", err
		}
	}

	return code, s.store.Save(code, url)
}

func (s *URLService) Expand(code string) (string, error) {
	return s.store.Get(code)
}

func (s *URLService) List() (map[string]string, error) {
	return s.store.List()
}

func (s *URLService) Delete(code string) error {
	return s.store.Delete(code)
}
