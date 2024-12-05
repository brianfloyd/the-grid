package pg

import (
	"github.com/brianfloyd/the-grid/internal"
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

func (u *User) Create(params internal.User) (internal.User, error) {
	user, err := u.q.InsertUser(InsertUserParams{
		id:   uuid.NewString(),
		name: params.Name,
	})
	if err != nil {
		return internal.User{}, err
	}
	return user, nil
}

func (u *User) ById(id string) (internal.User, error) {
	user, err := u.q.ById(id)
	if err != nil {
		return internal.User{}, err
	}
	return user, nil
}

func (u *User) ByName(name string) (internal.User, error) {
	user, err := u.q.ByName(name)
	if err != nil {
		return internal.User{}, err
	}
	return user, nil
}

func (u *User) List() ([]internal.User, error) {
	users, err := u.q.List()
	if err != nil {
		return nil, err
	}
	return users, nil
}
