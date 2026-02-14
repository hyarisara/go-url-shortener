package main

type UserStore interface {
	SaveUser(username, passwordHash string) error
	GetUser(username string) (string, error) // returns hashed password
	UserExists(username string) bool
}
