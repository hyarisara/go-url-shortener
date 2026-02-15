package actions

import (
	"go-url-shortener/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type UserAction struct {
	Store store.UserStore
}

func NewUserAction(s store.UserStore) *UserAction {
	return &UserAction{Store: s}
}

func (a *UserAction) Register(username, password string) error {
	if a.Store.UserExists(username) {
		return nil
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return a.Store.SaveUser(username, string(hash))
}

func (a *UserAction) Login(username, password string) bool {
	hash, err := a.Store.GetUser(username)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
