package handler

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/go-chi/chi"
)

type ILoginViewService interface {
	GetLoginPage() templ.Component
}

type LoginViewHandler struct {
	svc ILoginViewService
}

func NewLoginViewHandler(svc ILoginViewService) *LoginViewHandler {
	return &LoginViewHandler{
		svc: svc,
	}
}

func (l *LoginViewHandler) Register(r *chi.Mux) {
	r.Get("/login", l.getLoginPage)
	r.Post(fmt.Sprintf("/login/{uid:%s}", m.UUID_REGEX), l.login)
}

func (l *LoginViewHandler) getLoginPage(w http.ResponseWriter, r *http.Request) {
	l.svc.GetLoginPage().Render(r.Context(), w)
}

func (l *LoginViewHandler) login(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s;Path=/;SameSite=Strict;Max-Age=31536000", m.COOKIE_UID, uid))
	w.Header().Add("HX-Location", "/")
	w.WriteHeader(200)
}
