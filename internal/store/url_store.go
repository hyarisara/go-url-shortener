package store

type URLStore interface {
	// Save a URL under a key (username::code)
	SaveURL(key, url string) error

	// Get the original URL from a short code (find key ending with ::code)
	GetURLByCode(code string) (string, error)

	// List all URLs for a specific username
	ListByUser(username string) (map[string]string, error)

	// Delete a URL by key
	DeleteURL(key string) error
}
