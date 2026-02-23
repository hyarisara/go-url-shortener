package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"go-url-shortener/internal/actions"
	"go-url-shortener/internal/middlewares"
	"go-url-shortener/internal/store"
)

type Handler struct {
	urlService *actions.URLService
	userStore  store.UserStore
	templates  *template.Template
}

func NewHandler(urlService *actions.URLService, userStore store.UserStore, t *template.Template) *Handler {
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
// URL FEATURES
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

	url := r.FormValue("url")
	custom := r.FormValue("custom")

	code, err := h.urlService.ShortenForUser(user, url, custom)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
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

	// ✅ Use 301 permanent redirect (mentor request)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	user := middlewares.GetSession(r)
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	q := r.URL.Query().Get("q")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "created"
	}

	page := 1
	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}
	pageSize := 20

	urls, err := h.urlService.ListForUserPaged(user, q, sort, page, pageSize)
	if err != nil {
		http.Error(w, "Failed to load URLs", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Username": user,
		"URLs":     urls,
		"Q":        q,
		"Sort":     sort,
		"Page":     page,
		"PrevPage": max(1, page-1),
		"NextPage": page + 1,
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
	userKey := user + "::" + code

	_ = h.urlService.Delete(userKey)
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}