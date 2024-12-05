package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brianfloyd/the-grid/internal"
	"github.com/go-chi/chi"
)

type UserService interface {
	ById(id string) (internal.User, error)
	List() ([]internal.User, error)
	Create(user internal.User) (internal.User, error)
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
	r.Get(fmt.Sprintf("/users/{id:%s}", uuidRegEx), u.byId)
	r.Get("/users", u.list)
	r.Post("/users", u.create)
}

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	User User `json:"user"`
}

func (u *UserHandler) create(w http.ResponseWriter, r *http.Request) {
	var request CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, BadRequest, "Could not convert given body to a create user request.", err)
		return
	}

	defer r.Body.Close()

	user, err := u.svc.Create(internal.User{
		Name: request.Name,
	})

	if err != nil {
		var userExistsError = &internal.UserExistsError{}
		if errors.As(err, &userExistsError) {
			renderErrorResponse(w, r, BadRequest, err.Error(), err)
		} else {
			renderErrorResponse(w, r, GenericError, "Creating a user failed.", err)
		}
		return
	}

	renderResponse(w, r,
		&CreateUserResponse{
			User: User{
				Id:        user.Id,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
			},
		},
		http.StatusCreated)
}

type GetUserByIdResponse struct {
	User User `json:"user"`
}

func (u *UserHandler) byId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	user, err := u.svc.ById(id)
	if err != nil {
		userNotFoundError := &internal.UserNotFoundError{}
		if errors.As(err, &userNotFoundError) {
			renderErrorResponse(w, r, NotFound, "User was not found.", err)
		} else {
			renderErrorResponse(w, r, GenericError, "Generic user exception.", err)
		}
		return
	}

	renderResponse(w, r,
		&GetUserByIdResponse{
			User: convertInternalUserToResponse(user),
		},
		http.StatusOK)
}

type ListUsersResponse struct {
	Users []User `json:"users"`
}

func (u *UserHandler) list(w http.ResponseWriter, r *http.Request) {
	internalUsers, err := u.svc.List()
	if err != nil {
		renderErrorResponse(w, r, GenericError, "Could not list users.", err)
		return
	}

	users := make([]User, len(internalUsers))
	for i, iu := range internalUsers {
		users[i] = convertInternalUserToResponse(iu)
	}

	renderResponse(w, r, &ListUsersResponse{
		Users: users,
	}, http.StatusOK)
}

func convertInternalUserToResponse(user internal.User) User {
	return User{
		Id:        user.Id,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}
