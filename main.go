package main

import (
	"html/template"
	"net/http"
	"os"
)

func main() {
	// create store
	jsonStore := NewJSONStore("data.json")

	// create service
	urlService := NewURLService(jsonStore)

	// templates
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// handlers
	h := NewHandler(urlService, templates)

	// static
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// routes
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/shorten", h.Shorten)
	http.HandleFunc("/r/", h.Redirect)
	http.HandleFunc("/list", h.List)
	http.HandleFunc("/delete/", h.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	println("Server running on :" + port)
	http.ListenAndServe(":"+port, nil)
}
