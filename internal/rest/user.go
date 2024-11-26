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

const uuidRegEx string = `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`

type UserService interface {
	ById(id string) (internal.User, error)
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
		renderErrorResponse(w, r, BadRequest, "Could not convert given body to a create user request.")
		return
	}

	defer r.Body.Close()

	user, err := u.svc.Create(internal.User{
		Name: request.Name,
	})

	if err != nil {
		var userExistsError = &internal.UserExistsError{}
		if errors.As(err, &userExistsError) {
			renderErrorResponse(w, r, BadRequest, err.Error())
		} else {
			renderErrorResponse(w, r, GenericError, "Creating a user failed.")
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
			renderErrorResponse(w, r, NotFound, "User was not found.")
		} else {
			renderErrorResponse(w, r, GenericError, "Generic user exception.")
		}
		return
	}

	renderResponse(w, r,
		&GetUserByIdResponse{
			User: User{
				Id:        user.Id,
				Name:      user.Name,
				CreatedAt: user.CreatedAt,
			},
		},
		http.StatusOK)
}
