package main

import (
	"encoding/json"
	"os"
	"sync"
)

type JSONUserStore struct {
	file string
	mu   sync.RWMutex
}

func NewJSONUserStore(file string) *JSONUserStore {
	return &JSONUserStore{file: file}
}

func (s *JSONUserStore) load() (map[string]string, error) {
	data := make(map[string]string)
	b, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}
	json.Unmarshal(b, &data)
	return data, nil
}

func (s *JSONUserStore) save(data map[string]string) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, b, 0644)
}

func (s *JSONUserStore) SaveUser(username, passwordHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, _ := s.load()
	data[username] = passwordHash
	return s.save(data)
}

func (s *JSONUserStore) GetUser(username string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, _ := s.load()
	pwd, ok := data[username]
	if !ok {
		return "", os.ErrNotExist
	}
	return pwd, nil
}

func (s *JSONUserStore) UserExists(username string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	data, _ := s.load()
	_, ok := data[username]
	return ok
}
