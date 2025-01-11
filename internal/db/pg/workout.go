package pg

import (
	"context"

	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/google/uuid"
)

type WorkoutQueries struct {
	db IDbConnection
}

func NewWorkoutQueries(conn IDbConnection) *WorkoutQueries {
	return &WorkoutQueries{
		db: conn,
	}
}

type WorkoutRepository struct {
	q *WorkoutQueries
}

func NewWorkoutRepository(conn IDbConnection) *WorkoutRepository {
	return &WorkoutRepository{
		q: NewWorkoutQueries(conn),
	}
}

func (w *WorkoutRepository) Create(ctx context.Context, userId string, params m.Workout) (m.Workout, error) {
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

	return w.q.InsertWorkout(ctx, InsertWorkoutParams{
		id:     workoutId,
		userId: userId,
		date:   params.Date,
		sets:   insertSetParams,
	})
}

func (w *WorkoutRepository) ByDate(ctx context.Context, userId string, date string) (m.Workout, error) {
	return w.q.ByDate(ctx, ByDateParams{
		userId: userId,
		date:   date,
	})
}

func (w *WorkoutRepository) ById(ctx context.Context, id string) (m.Workout, error) {
	return w.q.ById(ctx, id)
}

func (w *WorkoutRepository) CreateSet(ctx context.Context, workoutId string, set m.Set) (m.Set, error) {
	return w.q.InsertSet(ctx, InsertSetParams{
		id:         uuid.NewString(),
		workoutId:  workoutId,
		exerciseId: set.ExerciseId,
		reps:       set.Reps,
		weight:     set.Weight,
		count:      set.Count,
	})
}

func (w *WorkoutRepository) UpdateSet(ctx context.Context, set m.Set) (m.Set, error) {
	return w.q.UpdateSet(ctx, InsertSetParams{
		id:         set.Id,
		workoutId:  set.WorkoutId,
		exerciseId: set.ExerciseId,
		reps:       set.Reps,
		weight:     set.Weight,
		count:      set.Count,
	})
}

func (w *WorkoutRepository) DeleteSets(ctx context.Context, setIds []string) error {
	return w.q.DeleteSets(ctx, setIds)
}
