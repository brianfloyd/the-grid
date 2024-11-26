package pg

import (
	"github.com/brianfloyd/the-grid/internal"
	"github.com/google/uuid"
)

type User struct {
	q *Queries
}

func NewUser(conn DbConnection) *User {
	return &User{
		q: New(conn),
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
