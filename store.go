package main

type UrlStore interface {
	Save(code string, url string) error
	Get(code string) (string, error)
	List() (map[string]string, error)
	Delete(code string) error
}

var store UrlStore

func initStore() {
	store = NewJSONStore("data.json")
}
