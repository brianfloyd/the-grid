package pg

import (
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

func (u *User) Create(params m.User) (m.User, error) {
	user, err := u.q.InsertUser(InsertUserParams{
		id:   uuid.NewString(),
		name: params.Name,
	})
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (u *User) ById(id string) (m.User, error) {
	user, err := u.q.ById(id)
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (u *User) ByName(name string) (m.User, error) {
	user, err := u.q.ByName(name)
	if err != nil {
		return m.User{}, err
	}
	return user, nil
}

func (u *User) List() ([]m.User, error) {
	users, err := u.q.List()
	if err != nil {
		return nil, err
	}
	return users, nil
}
