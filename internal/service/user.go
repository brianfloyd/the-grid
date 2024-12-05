package service

import (
	"errors"
	"fmt"

	"github.com/brianfloyd/the-grid/internal"
	"github.com/brianfloyd/the-grid/internal/db"
)

type UserService interface {
	List() ([]internal.User, error)
	ById(id string) (internal.User, error)
	Create(user internal.User) (internal.User, error)
}

type UserRepository interface {
	List() ([]internal.User, error)
	ById(id string) (internal.User, error)
	ByName(name string) (internal.User, error)
	Create(user internal.User) (internal.User, error)
}

type User struct {
	repo UserRepository
}

func NewUser(repo UserRepository) *User {
	return &User{
		repo: repo,
	}
}

func (u *User) ById(id string) (internal.User, error) {
	user, err := u.repo.ById(id)
	if err != nil {
		if errors.Is(err, db.ErrDbNotFound) {
			return internal.User{}, &internal.UserNotFoundError{Message: "Could not find user by the given id."}
		} else {
			return internal.User{}, &internal.GenericUserError{Message: "An unexpected exception occurred while performing the user operation."}
		}
	}
	return user, nil
}

func (u *User) Create(user internal.User) (internal.User, error) {
	exists, err := u.doesUserExistByName(user.Name)
	if err != nil {
		return internal.User{}, &internal.GenericUserError{Message: "Could not verify if the user already exists."}
	}

	if exists {
		return internal.User{}, &internal.UserExistsError{Message: fmt.Sprintf("User with the name %s already exists.", user.Name)}
	}

	createdUser, err := u.repo.Create(user)
	if err != nil {
		return internal.User{}, errors.Join(&internal.GenericUserError{Message: "An unexpected exception occurred while performing the user operation."}, err)
	}
	return createdUser, nil
}

func (u *User) List() ([]internal.User, error) {
	users, err := u.repo.List()
	if err != nil {
		return nil, errors.Join(&internal.GenericUserError{Message: "An error occurred while listing users"}, err)
	}
	return users, nil
}

func (u *User) doesUserExistByName(name string) (bool, error) {
	_, err := u.repo.ByName(name)

	if err == nil {
		// A user was successfully returned by their name.
		return true, nil
	}

	if errors.Is(err, db.ErrDbNotFound) {
		return false, nil
	}

	return false, err
}
