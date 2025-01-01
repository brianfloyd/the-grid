package pg

import (
	"context"

	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/google/uuid"
)

type WorkoutQueries struct {
	db DbConnection
}

func NewWorkoutQueries(conn DbConnection) *WorkoutQueries {
	return &WorkoutQueries{
		db: conn,
	}
}

type Workout struct {
	q *WorkoutQueries
}

func NewWorkout(conn DbConnection) *Workout {
	return &Workout{
		q: NewWorkoutQueries(conn),
	}
}

func (w *Workout) Create(ctx context.Context, userId string, params m.Workout) (m.Workout, error) {
	workoutId := uuid.NewString()

	insertSetParams := make([]InsertSetParams, len(params.Sets))
	for idx, set := range params.Sets {
		insertSetParams[idx] = InsertSetParams{
			id:         uuid.NewString(),
			workoutId:  workoutId,
			exerciseId: set.ExerciseId,
			reps:       set.Reps,
			weight:     set.Weight,
			count:      set.Count,
		}
	}

	workout, err := w.q.InsertWorkout(ctx, InsertWorkoutParams{
		id:     workoutId,
		userId: userId,
		date:   params.Date,
		sets:   insertSetParams,
	})

	if err != nil {
		return m.Workout{}, err
	}
	return workout, nil
}

func (w *Workout) ByDate(ctx context.Context, userId string, date string) (m.Workout, error) {
	workout, err := w.q.ByDate(ctx, ByDateParams{
		userId: userId,
		date:   date,
	})

	if err != nil {
		return m.Workout{}, err
	}
	return workout, nil
}
