package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*.html"))


func main() {
	initStore()
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/r/", redirectHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/delete/", deleteHandler)

	println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", nil)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	url := r.FormValue("url")
	code, err := ShortenURL(url)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	result := "http://localhost:8080/r/" + code
	templates.ExecuteTemplate(w, "index.html", result)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]

	url, err := ExpandURL(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := store.List()
	templates.ExecuteTemplate(w, "list.html", data)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/delete/"):]
	store.Delete(code)
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}