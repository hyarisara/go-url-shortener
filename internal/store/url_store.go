package store

type URLStore interface {
	SaveURL(key, originalURL string) error
	GetURLByCode(code string) (string, error)
	ListByUser(username string) (map[string]string, error)
	DeleteURL(key string) error

	// ✅ Step 2: fast search/sort/paging
	ListByUserPaged(username, q, sort string, page, pageSize int) (map[string]string, error)
}