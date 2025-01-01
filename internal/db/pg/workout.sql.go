package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/brianfloyd/the-grid/internal/db"
	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/jackc/pgx/v5"
)

const InsertWorkout = `
insert into the_grid_go.workout (wrk_id, wrk_usr_id, wrk_date, wrk_created_at, wrk_modified_at) values ($1, $2, $3, current_timestamp, current_timestamp)
returning wrk_id, wrk_usr_id, wrk_date, wrk_created_at, wrk_modified_at
`

const InsertSet = `
insert into the_grid_go.set (set_id, set_wrk_id, set_exr_id, set_reps, set_weight, set_count) values ($1, $2, $3, $4, $5, $6)
returning set_id, set_wrk_id, set_exr_id, set_reps, set_weight, set_count
`

type InsertWorkoutParams struct {
	id     string
	userId string
	date   string
	sets   []InsertSetParams
}

type InsertSetParams struct {
	id         string
	workoutId  string
	exerciseId string
	reps       uint64
	weight     uint64
	count      uint64
}

func (q *WorkoutQueries) InsertWorkout(ctx context.Context, args InsertWorkoutParams) (m.Workout, error) {
	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return m.Workout{}, err
	}

	workout, err := scanWorkout(tx.QueryRow(ctx, InsertWorkout, args.id, args.userId, args.date))
	if err != nil {
		tx.Rollback(ctx)
		return m.Workout{}, err
	}

	batch := &pgx.Batch{}
	for _, set := range args.sets {
		batch.Queue(InsertSet, set.id, args.id, set.exerciseId, set.reps, set.weight, set.count)
	}
	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	sets := []m.Set{}
	for i := 0; i < len(args.sets); i++ {
		rows, err := results.Query()
		if err != nil {
			tx.Rollback(ctx)
			return m.Workout{}, err
		}
		scannedSets, err := scanSets(rows)
		if err != nil {
			tx.Rollback(ctx)
			return m.Workout{}, err
		}
		sets = append(sets, scannedSets...)
	}
	workout.Sets = sets

	err = results.Close()
	if err != nil {
		tx.Rollback(ctx)
		return m.Workout{}, err
	}

	tx.Commit(ctx)
	return workout, nil
}

const SelectWorkoutByDate = `
select * from the_grid_go.workout where wrk_usr_id = $1 and wrk_date = $2
`

const SelectSetsByWorkoutId = `
select * from the_grid_go.set where set_wrk_id = $1
`

type ByDateParams struct {
	userId string
	date   string
}

func (q *WorkoutQueries) ByDate(ctx context.Context, args ByDateParams) (m.Workout, error) {
	tx, err := q.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return m.Workout{}, err
	}

	workout, err := scanWorkout(tx.QueryRow(ctx, SelectWorkoutByDate, args.userId, args.date))
	if err != nil {
		tx.Rollback(ctx)
		return m.Workout{}, err
	}

	rows, err := tx.Query(ctx, SelectSetsByWorkoutId, workout.Id)
	if err != nil {
		tx.Rollback(ctx)
		return m.Workout{}, err
	}
	sets, err := scanSets(rows)
	if err != nil {
		tx.Rollback(ctx)
		return m.Workout{}, err
	}
	workout.Sets = sets

	tx.Commit(ctx)
	return workout, nil
}

// Date is represented as a time in the database but is handled in the domain as a string.
type DatabaseWorkout struct {
	id         string
	userId     string
	date       time.Time
	createdAt  time.Time
	modifiedAt time.Time
}

func scanWorkout(row pgx.Row) (m.Workout, error) {
	dto := DatabaseWorkout{}
	err := row.Scan(
		&dto.id,
		&dto.userId,
		&dto.date,
		&dto.createdAt,
		&dto.modifiedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return m.Workout{}, db.ErrDbNotFound
		} else {
			return m.Workout{}, errors.Join(db.ErrDbGeneric, err)
		}
	}

	return m.Workout{
		Id:         dto.id,
		UserId:     dto.userId,
		Date:       dto.date.Format("01/02/2006"),
		CreatedAt:  dto.createdAt,
		ModifiedAt: dto.modifiedAt,
	}, nil
}

func scanSets(rows pgx.Rows) ([]m.Set, error) {
	var sets []m.Set
	for rows.Next() {
		set := m.Set{}
		err := rows.Scan(
			&set.Id,
			&set.WorkoutId,
			&set.ExerciseId,
			&set.Reps,
			&set.Weight,
			&set.Count,
		)

		if err != nil {
			return sets, err
		}

		sets = append(sets, set)
	}

	err := rows.Err()
	if err != nil {
		return sets, err
	}

	return sets, nil
}
