package pg

import (
	"context"

	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/google/uuid"
)

type UserQueries struct {
	db DbConnection
}

func NewUserQueries(conn DbConnection) *UserQueries {
	return &UserQueries{
		db: conn,
	}
}

type User struct {
	q *UserQueries
}

func NewUser(conn DbConnection) *User {
	return &User{
		q: NewUserQueries(conn),
	}
}

func (u *User) Create(ctx context.Context, params m.User) (m.User, error) {
	user, err := u.q.InsertUser(ctx, InsertUserParams{
		id:   uuid.NewString(),
		name: params.Name,
	})
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (u *User) ById(ctx context.Context, id string) (m.User, error) {
	user, err := u.q.ById(ctx, id)
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (u *User) ByName(ctx context.Context, name string) (m.User, error) {
	user, err := u.q.ByName(ctx, name)
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (u *User) List(ctx context.Context) ([]m.User, error) {
	users, err := u.q.List(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}
