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

// ShortenForUser stores URL per user
func (s *URLService) ShortenForUser(username, url, customCode string) (string, error) {
	code := customCode
	var err error
	if code == "" {
		code, err = generateCode()
		if err != nil {
			return "", err
		}
	}
	userKey := username + "::" + code
	return code, s.store.Save(userKey, url)
}

// ListForUser lists URLs for a specific user
func (s *URLService) ListForUser(username string) (map[string]string, error) {
	allData, _ := s.store.List()
	userData := make(map[string]string)
	for k, v := range allData {
		prefix := username + "::"
		if len(k) > len(prefix) && k[:len(prefix)] == prefix {
			userData[k[len(prefix):]] = v
		}
	}
	return userData, nil
}

func (s *URLService) Expand(code string) (string, error) {
	return s.store.Get(code)
}

func (s *URLService) Delete(code string) error {
	return s.store.Delete(code)
}
