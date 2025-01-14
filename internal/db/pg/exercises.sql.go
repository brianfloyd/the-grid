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

type InsertDefaultExerciseParams struct {
	id         string
	userId     string
	exerciseId string
	weight     uint64
	reps       uint64
}

type UpdateDefaultExerciseParams struct {
	id     string
	weight uint64
	reps   uint64
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
		select exr_id, exr_group, exr_name from the_grid_go.exercise where exr_group = upper($1)
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

func (q *ExercisesQueries) SelectExerciseDefaultById(ctx context.Context, id string) (m.ExerciseDefault, error) {
	const SelectExerciseDefaultById = `
		select edf_id, edf_usr_id, edf_exr_id, edf_weight, edf_reps from the_grid_go.exercise_default where edf_id = $1
	`
	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return m.ExerciseDefault{}, err
	}

	exerciseDefault, err := scanExerciseDefault(tx.QueryRow(ctx, SelectExerciseDefaultById, id))
	if err != nil {
		tx.Rollback(ctx)
		return m.ExerciseDefault{}, err
	}

	tx.Commit(ctx)
	return exerciseDefault, nil
}

func (q *ExercisesQueries) SelectExerciseDefault(ctx context.Context, userId string, exerciseId string) (m.ExerciseDefault, error) {
	const SelectExerciseDefaultByUserIdAndExerciseId = `
		select edf_id, edf_usr_id, edf_exr_id, edf_weight, edf_reps from the_grid_go.exercise_default where edf_usr_id = $1 and edf_exr_id = $2
	`
	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return m.ExerciseDefault{}, err
	}

	exerciseDefault, err := scanExerciseDefault(tx.QueryRow(ctx, SelectExerciseDefaultByUserIdAndExerciseId, userId, exerciseId))
	if err != nil {
		tx.Rollback(ctx)
		return m.ExerciseDefault{}, err
	}

	tx.Commit(ctx)
	return exerciseDefault, nil
}

func (q *ExercisesQueries) SelectAllExerciseDefaultsForUser(ctx context.Context, userId string) ([]m.ExerciseDefault, error) {
	const SelectAllDefaultExercisesForUser = `
		select edf_id, edf_usr_id, edf_exr_id, edf_weight, edf_reps from the_grid_go.exercise_default where edf_usr_id = $1
	`
	return q.selectExerciseDefaultRows(ctx, SelectAllDefaultExercisesForUser, userId)
}

func (q *ExercisesQueries) SelectAllExerciseDefaults(ctx context.Context) ([]m.ExerciseDefault, error) {
	const SelectAllDefaultExercises = `
		select edf_id, edf_usr_id, edf_exr_id, edf_weight, edf_reps from the_grid_go.exercise_default
	`

	return q.selectExerciseDefaultRows(ctx, SelectAllDefaultExercises)
}

func (q *ExercisesQueries) InsertExerciseDefault(ctx context.Context, params InsertDefaultExerciseParams) (m.ExerciseDefault, error) {
	const InsertDefaultExercise = `
		insert into the_grid_go.exercise_default (edf_id, edf_usr_id, edf_exr_id, edf_weight, edf_reps) values ($1, $2, $3, $4, $5)
		returning edf_id, edf_usr_id, edf_exr_id, edf_weight, edf_reps
	`

	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})

	if err != nil {
		return m.ExerciseDefault{}, err
	}

	exerciseDefault, err := scanExerciseDefault(tx.QueryRow(ctx, InsertDefaultExercise, params.id, params.userId, params.exerciseId, params.weight, params.reps))

	if err != nil {
		tx.Rollback(ctx)
		return m.ExerciseDefault{}, err
	}

	tx.Commit(ctx)
	return exerciseDefault, nil
}

func (q *ExercisesQueries) UpdateExerciseDefault(ctx context.Context, params UpdateDefaultExerciseParams) (m.ExerciseDefault, error) {
	const UpdateDefaultExercise = `
		update the_grid_go.exercise_default set edf_weight = $1, edf_reps = $2 where edf_id = $3
		returning edf_id, edf_usr_id, edf_exr_id, edf_weight, edf_reps
	`

	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})

	if err != nil {
		return m.ExerciseDefault{}, err
	}
	exerciseDefault, err := scanExerciseDefault(tx.QueryRow(ctx, UpdateDefaultExercise, params.weight, params.reps, params.id))

	if err != nil {
		tx.Rollback(ctx)
		return m.ExerciseDefault{}, err
	}

	tx.Commit(ctx)
	return exerciseDefault, nil
}

func (q *ExercisesQueries) selectExerciseDefaultRows(ctx context.Context, sql string, args ...any) ([]m.ExerciseDefault, error) {
	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})

	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, sql, args...)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	exerciseDefaults, err := scanExerciseDefaults(rows)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	tx.Commit(ctx)
	return exerciseDefaults, nil
}

func scanExerciseDefaults(rows pgx.Rows) ([]m.ExerciseDefault, error) {
	exerciseDefaults := []m.ExerciseDefault{}
	for rows.Next() {
		exerciseDefault := m.ExerciseDefault{}
		err := rows.Scan(
			&exerciseDefault.Id,
			&exerciseDefault.UserId,
			&exerciseDefault.ExerciseId,
			&exerciseDefault.Weight,
			&exerciseDefault.Reps,
		)
		if err != nil {
			return nil, err
		}
		exerciseDefaults = append(exerciseDefaults, exerciseDefault)
	}
	return exerciseDefaults, nil
}

func scanExerciseDefault(row pgx.Row) (m.ExerciseDefault, error) {
	exerciseDefault := m.ExerciseDefault{}
	err := row.Scan(
		&exerciseDefault.Id,
		&exerciseDefault.UserId,
		&exerciseDefault.ExerciseId,
		&exerciseDefault.Weight,
		&exerciseDefault.Reps,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.ExerciseDefault{}, db.ErrDbNotFound
		} else {
			return m.ExerciseDefault{}, errors.Join(db.ErrDbGeneric, err)
		}
	}

	return exerciseDefault, nil
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
			return m.Exercise{}, errors.Join(db.ErrDbGeneric)
		}
	}

	return exercise, nil
}
