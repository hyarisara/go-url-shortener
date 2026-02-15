package handlers

import (
	"html/template"
	"net/http"

	"go-url-shortener/internal/actions"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/store"
)

type Handler struct {
	urlService *actions.URLService
	userStore  store.UserStore
	templates  *template.Template
}

func NewHandler(
	urlService *actions.URLService,
	userStore store.UserStore,
	t *template.Template,
) *Handler {
	return &Handler{
		urlService: urlService,
		userStore:  userStore,
		templates:  t,
	}
}

//////////////////////
// AUTH
//////////////////////

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if h.userStore.UserExists(username) {
			h.templates.ExecuteTemplate(w, "register.html", "User already exists")
			return
		}

		hash, err := middlewares.HashPassword(password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		if err := h.userStore.SaveUser(username, hash); err != nil {
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		middlewares.SetSession(w, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.templates.ExecuteTemplate(w, "register.html", nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hash, err := h.userStore.GetUser(username)
		if err != nil || hash == "" || !middlewares.CheckPassword(hash, password) {
			h.templates.ExecuteTemplate(w, "login.html", "Invalid credentials")
			return
		}

		middlewares.SetSession(w, username)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	h.templates.ExecuteTemplate(w, "login.html", nil)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	middlewares.ClearSession(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

//////////////////////
// URL
//////////////////////

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"Username": user,
	}

	h.templates.ExecuteTemplate(w, "index.html", data)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	originalURL := r.FormValue("url")
	custom := r.FormValue("custom")

	code, err := h.urlService.ShortenForUser(user, originalURL, custom)
	if err != nil {
		http.Error(w, "Error shortening", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Username": user,
		"ShortURL": "/r/" + code,
	}

	h.templates.ExecuteTemplate(w, "index.html", data)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]

	url, err := h.urlService.Expand(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	urls, err := h.urlService.ListForUser(user)
	if err != nil {
		http.Error(w, "Error loading URLs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Username": user,
		"URLs":     urls,
	}

	h.templates.ExecuteTemplate(w, "list.html", data)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	code := r.URL.Path[len("/delete/"):]
	key := user + "::" + code

	if err := h.urlService.Delete(key); err != nil {
		http.Error(w, "Error deleting", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/list", http.StatusSeeOther)
}
