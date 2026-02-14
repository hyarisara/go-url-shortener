package main

import (
	"html/template"
	"net/http"
)

type Handler struct {
	service   *URLService
	userStore UserStore
	templates *template.Template
}

func NewHandler(s *URLService, u UserStore, t *template.Template) *Handler {
	return &Handler{service: s, userStore: u, templates: t}
}

// Home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	data := map[string]interface{}{
		"Username": user,
	}
	h.templates.ExecuteTemplate(w, "index.html", data)
}

// Shorten URL
func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	url := r.FormValue("url")
	custom := r.FormValue("custom")
	code, _ := h.service.ShortenForUser(user, url, custom)
	data := map[string]interface{}{
		"Username": user,
		"ShortURL": "/r/" + code,
	}
	h.templates.ExecuteTemplate(w, "index.html", data)
}

// Redirect short URL
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]
	url, err := h.service.Expand(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

// List user URLs
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	urls, _ := h.service.ListForUser(user)
	data := map[string]interface{}{
		"Username": user,
		"URLs":     urls,
	}
	h.templates.ExecuteTemplate(w, "list.html", data)
}

// Delete user URL
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	code := r.URL.Path[len("/delete/"):]
	userKey := user + "::" + code
	h.service.Delete(userKey)
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

// Register page
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		if h.userStore.UserExists(username) {
			w.Write([]byte("User already exists"))
			return
		}
		hash, _ := HashPassword(password)
		h.userStore.SaveUser(username, hash)
		SetSession(w, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.templates.ExecuteTemplate(w, "register.html", nil)
}

// Login page
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")
		hash, err := h.userStore.GetUser(username)
		if err != nil || !CheckPassword(hash, password) {
			w.Write([]byte("Invalid credentials"))
			return
		}
		SetSession(w, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	h.templates.ExecuteTemplate(w, "login.html", nil)
}
