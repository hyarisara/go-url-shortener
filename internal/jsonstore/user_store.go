package jsonstore

import (
	"encoding/json"
	"errors"
	"os"
)

type JSONUserStore struct {
	file string
}

func NewJSONUserStore(file string) *JSONUserStore {
	return &JSONUserStore{file: file}
}

func (s *JSONUserStore) loadUsers() (map[string]string, error) {
	data := make(map[string]string)
	file, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}
	json.Unmarshal(file, &data)
	return data, nil
}

func (s *JSONUserStore) saveUsers(users map[string]string) error {
	bytes, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, bytes, 0644)
}

func (s *JSONUserStore) SaveUser(username, hash string) error {
	users, _ := s.loadUsers()
	users[username] = hash
	return s.saveUsers(users)
}

func (s *JSONUserStore) GetUser(username string) (string, error) {
	users, _ := s.loadUsers()
	hash, ok := users[username]
	if !ok {
		return "", errors.New("user not found")
	}
	return hash, nil
}

func (s *JSONUserStore) UserExists(username string) bool {
	users, _ := s.loadUsers()
	_, ok := users[username]
	return ok
}
