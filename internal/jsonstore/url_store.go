package jsonstore

import (
	"encoding/json"
	"os"
)

type URLRecord struct {
	URL   string `json:"url"`
	Owner string `json:"owner"`
}

type JSONURLStore struct {
	file string
}

func NewJSONURLStore(file string) *JSONURLStore {
	return &JSONURLStore{file: file}
}

func (s *JSONURLStore) loadData() (map[string]URLRecord, error) {
	data := make(map[string]URLRecord)
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

func (s *JSONURLStore) saveData(data map[string]URLRecord) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, bytes, 0644)
}

func (s *JSONURLStore) SaveURL(code, url, owner string) error {
	data, _ := s.loadData()
	data[code] = URLRecord{URL: url, Owner: owner}
	return s.saveData(data)
}

func (s *JSONURLStore) ExpandURL(code string) (string, error) {
	data, _ := s.loadData()
	rec, ok := data[code]
	if !ok {
		return "", nil
	}
	return rec.URL, nil
}

func (s *JSONURLStore) ListForUser(owner string) (map[string]string, error) {
	data, _ := s.loadData()
	userURLs := make(map[string]string)
	for code, rec := range data {
		if rec.Owner == owner {
			userURLs[code] = rec.URL
		}
	}
	return userURLs, nil
}

func (s *JSONURLStore) DeleteURL(code string) error {
	data, _ := s.loadData()
	delete(data, code)
	return s.saveData(data)
}
