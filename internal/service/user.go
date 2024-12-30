package service

import (
	"errors"
	"fmt"

	"github.com/brianfloyd/the-grid/internal/db"
	m "github.com/brianfloyd/the-grid/internal/model"
)

type IUserService interface {
	List() ([]m.User, error)
	ById(id string) (m.User, error)
	Create(user m.User) (m.User, error)
}

type UserRepository interface {
	List() ([]m.User, error)
	ById(id string) (m.User, error)
	ByName(name string) (m.User, error)
	Create(user m.User) (m.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) ById(id string) (m.User, error) {
	user, err := u.repo.ById(id)
	if err != nil {
		if errors.Is(err, db.ErrDbNotFound) {
			return m.User{}, &m.UserNotFoundError{Message: "Could not find user by the given id."}
		} else {
			return m.User{}, &m.GenericUserError{Message: "An unexpected exception occurred while performing the user operation."}
		}
	}
	return user, nil
}

func (u *UserService) Create(user m.User) (m.User, error) {
	exists, err := u.doesUserExistByName(user.Name)
	if err != nil {
		return m.User{}, &m.GenericUserError{Message: "Could not verify if the user already exists."}
	}

	if exists {
		return m.User{}, &m.UserExistsError{Message: fmt.Sprintf("User with the name %s already exists.", user.Name)}
	}

	createdUser, err := u.repo.Create(user)
	if err != nil {
		return m.User{}, errors.Join(&m.GenericUserError{Message: "An unexpected exception occurred while performing the user operation."}, err)
	}
	return createdUser, nil
}

func (u *UserService) List() ([]m.User, error) {
	users, err := u.repo.List()
	if err != nil {
		return nil, errors.Join(&m.GenericUserError{Message: "An error occurred while listing users"}, err)
	}
	return users, nil
}

func (u *UserService) doesUserExistByName(name string) (bool, error) {
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
