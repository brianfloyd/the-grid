package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brianfloyd/the-grid/internal"
	"github.com/brianfloyd/the-grid/internal/db"
	"github.com/jackc/pgx/v5"
)

const InsertUser = `
insert into the_grid.user (usr_id, usr_name) values ($1, $2)
returning usr_id, usr_name, usr_created_at
`

type InsertUserParams struct {
	id   string
	name string
}

func (q *Queries) InsertUser(args InsertUserParams) (internal.User, error) {
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

const ById = `
select usr_id, usr_name, usr_created_at from the_grid.user where usr_id = $1
`

func (q *Queries) ById(id string) (internal.User, error) {
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

const ByName = `
select usr_id, usr_name, usr_created_at from the_grid.user where usr_name = $1
`

func (q *Queries) ByName(name string) (internal.User, error) {
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
