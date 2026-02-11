package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const fileName = "data.json"

func loadData() (map[string]string, error) {
	data := make(map[string]string)

	file, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}

	json.Unmarshal(file, &data)
	return data, nil
}

func saveData(data map[string]string) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, bytes, 0644)
}

func generateCode() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func ShortenURL(url string) (string, error) {
	data, err := loadData()
	if err != nil {
		return "", err
	}

	code, err := generateCode()
	if err != nil {
		return "", err
	}

	data[code] = url
	err = saveData(data)
	return code, err
}

func ExpandURL(code string) (string, error) {
	data, err := loadData()
	if err != nil {
		return "", err
	}

	url, ok := data[code]
	if !ok {
		return "", errors.New("code not found")
	}

	return url, nil
}

func ListURLs() {
	data, err := loadData()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for k, v := range data {
		fmt.Printf("%s -> %s\n", k, v)
	}
}
