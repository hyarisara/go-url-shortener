package store

type UserStore interface {
	SaveUser(username, hash string) error
	GetUser(username string) (string, error)
	UserExists(username string) bool
}
