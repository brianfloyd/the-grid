package rest

import (
	"time"
)

type UserResponse struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	User UserResponse `json:"user"`
}

type GetUserByIdResponse struct {
	User UserResponse `json:"user"`
}
