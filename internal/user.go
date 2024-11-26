package internal

import "time"

type User struct {
	Id        string
	Name      string
	CreatedAt time.Time
}

type GenericUserError struct {
	Message string
}

func (e *GenericUserError) Error() string {
	return e.Message
}

type UserNotFoundError struct {
	Message string
}

func (e *UserNotFoundError) Error() string {
	return e.Message
}

type UserExistsError struct {
	Message string
}

func (e *UserExistsError) Error() string {
	return e.Message
}
