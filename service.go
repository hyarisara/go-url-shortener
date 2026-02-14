package main

import (
	"crypto/rand"
	"encoding/hex"
)

func generateCode() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func ShortenURL(url string) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}

	return code, store.Save(code, url)
}

func ExpandURL(code string) (string, error) {
	return store.Get(code)
}
