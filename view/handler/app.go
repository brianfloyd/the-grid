package handler

import (
	"net/http"

	"github.com/brianfloyd/the-grid/view/template"
	"github.com/go-chi/chi"
)

type AppViewHandler struct {
}

func NewAppViewHandler() *AppViewHandler {
	return &AppViewHandler{}
}

func (a *AppViewHandler) Register(r *chi.Mux) {
	r.Get("/", a.getApp)
}

func (a *AppViewHandler) getApp(w http.ResponseWriter, r *http.Request) {
	userId := r.CookiesNamed("x-the-grid-uid")

	if len(userId) == 0 {
		w.Header().Add("Location", "/login")
		w.WriteHeader(301)
		return
	}

	template.App().Render(r.Context(), w)
}
