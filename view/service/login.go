package service

import (
	"context"

	"github.com/a-h/templ"
	im "github.com/brianfloyd/the-grid/internal/model"
	is "github.com/brianfloyd/the-grid/internal/service"
	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/brianfloyd/the-grid/view/template/page"
)

type ILoginViewService interface {
	GetLoginPage(ctx context.Context) templ.Component
}

type LoginViewService struct {
	userService is.IUserService
}

func NewLoginViewService(userService is.IUserService) *LoginViewService {
	return &LoginViewService{
		userService: userService,
	}
}

func (l *LoginViewService) GetLoginPage(ctx context.Context) templ.Component {
	users, err := l.userService.List(ctx)
	if err != nil {
		panic("crash")
	}
	userViews := makeUserViews(users)
	return page.LoginPage(userViews)
}

func makeUserViews(users []im.User) []m.UserView {
	viewUsers := make([]m.UserView, len(users))
	for i, u := range users {
		viewUsers[i] = m.UserView{
			Id:   u.Id,
			Name: u.Name,
		}
	}
	return viewUsers
}
