package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brianfloyd/the-grid/internal"
	"github.com/brianfloyd/the-grid/internal/db"
	"github.com/jackc/pgx/v5"
)

type InsertUserParams struct {
	id   string
	name string
}

func (q *UserQueries) InsertUser(args InsertUserParams) (internal.User, error) {
	const InsertUser = `
		insert into the_grid_go.user (usr_id, usr_name) values ($1, $2)
		returning usr_id, usr_name, usr_created_at
	`
	tx, err := q.db.BeginTx(context.TODO(), pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})

	if err != nil {
		return internal.User{}, err
	}

	user, err := scanUser(tx.QueryRow(context.TODO(), InsertUser, args.id, args.name))
	if err != nil {
		tx.Rollback(context.TODO())
		return internal.User{}, err
	}
	tx.Commit(context.TODO())
	return user, nil
}

func (q *UserQueries) ById(id string) (internal.User, error) {
	const ById = `
		select usr_id, usr_name, usr_created_at from the_grid_go.user where usr_id = $1
	`

	tx, err := q.db.BeginTx(context.TODO(), pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return internal.User{}, err
	}

	user, err := scanUser(tx.QueryRow(context.TODO(), ById, id))
	if err != nil {
		tx.Rollback(context.TODO())
		return internal.User{}, err
	}
	tx.Commit(context.TODO())
	return user, nil
}

func (q *UserQueries) ByName(name string) (internal.User, error) {
	const ByName = `
		select usr_id, usr_name, usr_created_at from the_grid_go.user where usr_name = $1
	`

	tx, err := q.db.BeginTx(context.TODO(), pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return internal.User{}, err
	}

	user, err := scanUser(tx.QueryRow(context.TODO(), ByName, name))
	if err != nil {
		tx.Rollback(context.TODO())
		return internal.User{}, err
	}
	tx.Commit(context.TODO())
	return user, nil
}

func (q *UserQueries) List() ([]internal.User, error) {
	const ListUsers = `
		select * from the_grid_go.user
	`

	tx, err := q.db.BeginTx(context.TODO(), pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(context.TODO(), ListUsers)
	if err != nil {
		tx.Rollback(context.TODO())
		return nil, err
	}

	users, err := scanUsers(rows)
	if err != nil {
		tx.Rollback(context.TODO())
		return nil, err
	}

	tx.Commit(context.TODO())
	return users, nil
}

func scanUsers(rows pgx.Rows) ([]internal.User, error) {
	users := []internal.User{}
	for rows.Next() {
		user := internal.User{}
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func scanUser(row pgx.Row) (internal.User, error) {
	user := internal.User{}
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return internal.User{}, db.ErrDbNotFound
		} else {
			return internal.User{}, db.ErrDbGeneric
		}
	}

	return user, nil
}
