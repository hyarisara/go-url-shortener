package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	// --- 1️⃣ Create store ---
	jsonStore := NewJSONStore("data.json")

	// --- 2️⃣ Create service ---
	urlService := NewURLService(jsonStore)

	// --- 3️⃣ Parse templates ---
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// --- 4️⃣ Initialize handlers ---
	h := NewHandler(urlService, templates)

	// --- 5️⃣ Serve static files ---
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// --- 6️⃣ Register routes ---
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/shorten", h.Shorten)
	http.HandleFunc("/r/", h.Redirect)
	http.HandleFunc("/list", h.List)
	http.HandleFunc("/delete/", h.Delete)

	// --- 7️⃣ Get port from environment (for hosting) ---
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}

	log.Println("Server running on port :" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
