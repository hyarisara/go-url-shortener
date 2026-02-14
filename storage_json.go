package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type JSONStore struct {
	file string
	mu   sync.Mutex
}

func NewJSONStore(file string) *JSONStore {
	return &JSONStore{file: file}
}

func (s *JSONStore) load() (map[string]string, error) {
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

func (s *JSONStore) save(data map[string]string) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, b, 0644)
}

func (s *JSONStore) Save(code string, url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	data[code] = url
	return s.save(data)
}

func (s *JSONStore) Get(code string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return "", err
	}

	u, ok := data[code]
	if !ok {
		return "", errors.New("not found")
	}
	return u, nil
}

func (s *JSONStore) List() (map[string]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.load()
}

func (s *JSONStore) Delete(code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.load()
	if err != nil {
		return err
	}

	delete(data, code)
	return s.save(data)
}