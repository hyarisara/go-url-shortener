package main

import (
	"html/template"
	"net/http"
)

type Handler struct {
	service   *URLService
	templates *template.Template
}

func NewHandler(s *URLService, t *template.Template) *Handler {
	return &Handler{service: s, templates: t}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "index.html", nil)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	url := r.FormValue("url")
	custom := r.FormValue("custom")
	code, err := h.service.Shorten(url, custom)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	result := "/r/" + code
	h.templates.ExecuteTemplate(w, "index.html", result)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/r/"):]
	url, err := h.service.Expand(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	data, _ := h.service.List()
	h.templates.ExecuteTemplate(w, "list.html", data)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[len("/delete/"):]
	h.service.Delete(code)
	http.Redirect(w, r, "/list", http.StatusSeeOther)
}
