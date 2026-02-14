package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	// --- 1️⃣ Create stores ---
	jsonStore := NewJSONStore("data.json")       // URL storage
	userStore := NewJSONUserStore("users.json") // User storage

	// --- 2️⃣ Create service ---
	urlService := NewURLService(jsonStore)

	// --- 3️⃣ Parse templates ---
	templates := template.Must(template.ParseGlob("templates/*.html"))

	// --- 4️⃣ Initialize handlers (pass both URL and User stores) ---
	h := NewHandler(urlService, userStore, templates)

	// --- 5️⃣ Serve static files ---
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// --- 6️⃣ Register routes ---
	http.HandleFunc("/", h.Home)
	http.HandleFunc("/shorten", h.Shorten)
	http.HandleFunc("/r/", h.Redirect)
	http.HandleFunc("/list", h.List)
	http.HandleFunc("/delete/", h.Delete)

	// --- 6️⃣ User auth routes ---
	http.HandleFunc("/register", h.Register)
	http.HandleFunc("/login", h.Login)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		ClearSession(w)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

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
