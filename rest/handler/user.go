package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	m "github.com/brianfloyd/the-grid/internal/model"
	rm "github.com/brianfloyd/the-grid/rest/model"
	"github.com/go-chi/chi"
)

type UserService interface {
	ById(ctx context.Context, id string) (m.User, error)
	List(ctx context.Context) ([]m.User, error)
	Create(ctx context.Context, user m.User) (m.User, error)
}

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (u *UserHandler) Register(r *chi.Mux) {
	r.Get(fmt.Sprintf("/users/{id:%s}", rm.UUID_REGEX), u.byId)
	r.Get("/users", u.list)
	r.Post("/users", u.create)
}

func (u *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var request rm.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, rm.BadRequest, "Could not convert given body to a create user request.", err)
		return
	}

	defer r.Body.Close()

	user, err := u.svc.Create(r.Context(), m.User{
		Name: request.Name,
	})

	if err != nil {
		var userExistsError = &m.UserExistsError{}
		if errors.As(err, &userExistsError) {
			renderErrorResponse(w, r, rm.BadRequest, err.Error(), err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "Creating a user failed.", err)
		}
		return
	}

	renderResponse(w, r,
		&rm.CreateUserResponse{
			User: rm.UserResponse{
				Id:        user.Id,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
			},
		},
		http.StatusCreated)
}

func (u *UserHandler) byId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := u.svc.ById(r.Context(), id)
	if err != nil {
		userNotFoundError := &m.UserNotFoundError{}
		if errors.As(err, &userNotFoundError) {
			renderErrorResponse(w, r, rm.NotFound, "User was not found.", err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "Generic user exception.", err)
		}
		return
	}

	renderResponse(w, r,
		&rm.GetUserByIdResponse{
			User: convertModelUserToResponseUser(user),
		},
		http.StatusOK)
}

func (u *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	mUsers, err := u.svc.List(r.Context())
	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "Could not list users.", err)
		return
	}

	users := make([]rm.UserResponse, len(mUsers))
	for i, iu := range mUsers {
		users[i] = convertModelUserToResponseUser(iu)
	}

	renderResponse(w, r, &rm.ListUsersResponse{
		Users: users,
	}, http.StatusOK)
}

func convertModelUserToResponseUser(user m.User) rm.UserResponse {
	return rm.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}
