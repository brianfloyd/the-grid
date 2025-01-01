package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/brianfloyd/the-grid/internal/db"
	"github.com/brianfloyd/the-grid/internal/logger"
	m "github.com/brianfloyd/the-grid/internal/model"
)

type IUserService interface {
	List(ctx context.Context) ([]m.User, error)
	ById(ctx context.Context, id string) (m.User, error)
	Create(ctx context.Context, user m.User) (m.User, error)
}

type UserRepository interface {
	List(ctx context.Context) ([]m.User, error)
	ById(ctx context.Context, id string) (m.User, error)
	ByName(ctx context.Context, name string) (m.User, error)
	Create(ctx context.Context, user m.User) (m.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) ById(ctx context.Context, id string) (m.User, error) {
	logger.TraceArgs(ctx, "Looking up user by id (%s).", id)
	user, err := u.repo.ById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrDbNotFound) {
			return m.User{}, &m.UserNotFoundError{Message: "Could not find user by the given id."}
		} else {
			return m.User{}, &m.GenericUserError{Message: "An unexpected exception occurred while performing the user operation."}
		}
	}
	return user, nil
}

func (u *UserService) Create(ctx context.Context, user m.User) (m.User, error) {
	logger.InfoArgs(ctx, "Creating user (%v).", user)
	exists, err := u.doesUserExistByName(ctx, user.Name)
	if err != nil {
		return m.User{}, &m.GenericUserError{Message: "Could not verify if the user already exists."}
	}

	if exists {
		return m.User{}, &m.UserExistsError{Message: fmt.Sprintf("User with the name %s already exists.", user.Name)}
	}

	createdUser, err := u.repo.Create(ctx, user)
	if err != nil {
		return m.User{}, errors.Join(&m.GenericUserError{Message: "An unexpected exception occurred while performing the user operation."}, err)
	}
	return createdUser, nil
}

func (u *UserService) List(ctx context.Context) ([]m.User, error) {
	logger.Trace(ctx, "Listing users.")
	users, err := u.repo.List(ctx)
	if err != nil {
		return nil, errors.Join(&m.GenericUserError{Message: "An error occurred while listing users"}, err)
	}
	return users, nil
}

func (u *UserService) doesUserExistByName(ctx context.Context, name string) (bool, error) {
	_, err := u.repo.ByName(ctx, name)

	if err == nil {
		// A user was successfully returned by their name.
		return true, nil
	}

	if errors.Is(err, db.ErrDbNotFound) {
		return false, nil
	}

	return false, err
}
