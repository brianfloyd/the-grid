package rest

import (
	"time"
)

type User struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	User User `json:"user"`
}

type GetUserByIdResponse struct {
	User User `json:"user"`
}
