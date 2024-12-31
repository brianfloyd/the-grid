package handler

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/brianfloyd/the-grid/view/template"
	"github.com/go-chi/chi"
)

type AppViewHandler struct {
	loginViewService ILoginViewService
	// appViewService IAppViewService
}

func NewAppViewHandler(loginViewService ILoginViewService) *AppViewHandler {
	return &AppViewHandler{
		loginViewService: loginViewService,
	}
}

func (a *AppViewHandler) Register(r *chi.Mux) {
	r.Get("/", a.getApp)
}

func (a *AppViewHandler) getApp(w http.ResponseWriter, r *http.Request) {
	userId := r.CookiesNamed("x-the-grid-uid")
	fmt.Printf("Cookie %v\n", userId)

	var component templ.Component
	if len(userId) != 0 {
		component = template.App()
	} else {
		component = a.loginViewService.GetLoginPage()
	}
	component.Render(r.Context(), w)
}
