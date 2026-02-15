package main

import (
	"html/template"
	"log"
	"net/http"

	"go-url-shortener/internal/actions"
	"go-url-shortener/internal/handlers"
	"go-url-shortener/internal/sqlite"
)

func main() {
	// 1️⃣ Stores
	userStore := sqlite.NewUserStore("data.db")
	urlStore := sqlite.NewURLStore("data.db")

	// 2️⃣ Service
	urlService := actions.NewURLService(urlStore)

	// 3️⃣ Templates
	templates := template.Must(template.ParseGlob("internal/templates/*.html"))

	// 4️⃣ Handlers
	h := handlers.NewHandler(urlService, userStore, templates)

	// 5️⃣ Static
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("internal/static")),
		),
	)

	// 6️⃣ Routes
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/shorten", h.Shorten)
	http.HandleFunc("/r/", h.Redirect)
	http.HandleFunc("/list", h.List)
	http.HandleFunc("/login", h.Login)
	http.HandleFunc("/register", h.Register)
	http.HandleFunc("/logout", h.Logout)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
