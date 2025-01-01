package pg

import (
	"context"
	"database/sql"
	"errors"

	"github.com/brianfloyd/the-grid/internal/db"
	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/jackc/pgx/v5"
)

type InsertExerciseParams struct {
	id    string
	name  string
	group string
}

func (q *ExercisesQueries) InsertExercise(ctx context.Context, args InsertExerciseParams) (m.Exercise, error) {
	const InsertExercise = `
		insert into the_grid_go.exercise (exr_id, exr_group, exr_name) values ($1, $2, $3)
		returning exr_id, exr_group, exr_name
	`

	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})

	if err != nil {
		return m.Exercise{}, err
	}

	user, err := scanExercise(tx.QueryRow(ctx, InsertExercise, args.id, args.group, args.name))
	if err != nil {
		tx.Rollback(ctx)
		return m.Exercise{}, err
	}
	tx.Commit(ctx)
	return user, nil
}

func (q *ExercisesQueries) List(ctx context.Context) ([]m.Exercise, error) {
	const ListExercises = `
		select exr_id, exr_group, exr_name from the_grid_go.exercise
	`

	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, ListExercises)

	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	exercises, err := scanExercises(rows)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	tx.Commit(ctx)
	return exercises, nil
}

func (q *ExercisesQueries) ListForGroup(ctx context.Context, group string) ([]m.Exercise, error) {
	const ListExercisesForGroup = `
		select exr_id, exr_group, exr_name from the_grid_go.exercise where exr_group = $1
	`

	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, ListExercisesForGroup, group)

	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	exercises, err := scanExercises(rows)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	tx.Commit(ctx)
	return exercises, nil
}

func scanExercises(rows pgx.Rows) ([]m.Exercise, error) {
	exercises := []m.Exercise{}
	for rows.Next() {
		exercise := m.Exercise{}
		err := rows.Scan(
			&exercise.Id,
			&exercise.Group,
			&exercise.Name,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}

func scanExercise(row pgx.Row) (m.Exercise, error) {
	exercise := m.Exercise{}
	err := row.Scan(
		&exercise.Id,
		&exercise.Group,
		&exercise.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.Exercise{}, db.ErrDbNotFound
		} else {
			return m.Exercise{}, db.ErrDbGeneric
		}
	}

	return exercise, nil
}
