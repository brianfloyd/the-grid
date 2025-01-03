package handler

import (
	"net/http"

	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/brianfloyd/the-grid/view/template"
	"github.com/go-chi/chi"
)

func GetUid(w http.ResponseWriter, r *http.Request) (string, bool) {
	userId := r.CookiesNamed(m.COOKIE_UID)

	if len(userId) == 0 {
		w.Header().Add("Location", "/login")
		w.WriteHeader(301)
		return "", false
	}

	return userId[0].Value, true
}

type AppViewHandler struct {
}

func NewAppViewHandler() *AppViewHandler {
	return &AppViewHandler{}
}

func (a *AppViewHandler) Register(r *chi.Mux) {
	r.Get("/", a.getApp)
}

func (a *AppViewHandler) getApp(w http.ResponseWriter, r *http.Request) {
	if _, ok := GetUid(w, r); ok {
		date := "01/01/2025"
		template.App(date).Render(r.Context(), w)
	}
}
