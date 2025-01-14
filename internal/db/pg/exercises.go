package pg

import (
	"context"

	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/google/uuid"
)

type ExercisesQueries struct {
	db IDbConnection
}

func NewExercisesQueries(conn IDbConnection) *ExercisesQueries {
	return &ExercisesQueries{
		db: conn,
	}
}

type ExercisesRepository struct {
	q *ExercisesQueries
}

func NewExercisesRepository(conn IDbConnection) *ExercisesRepository {
	return &ExercisesRepository{
		q: NewExercisesQueries(conn),
	}
}

func (e *ExercisesRepository) Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error) {
	createdExercise, err := e.q.InsertExercise(ctx, InsertExerciseParams{
		id:    uuid.NewString(),
		group: string(exercise.Group),
		name:  exercise.Name,
	})

	if err != nil {
		return m.Exercise{}, err
	}

	return createdExercise, nil
}

func (e *ExercisesRepository) List(ctx context.Context) ([]m.Exercise, error) {
	exercises, err := e.q.List(ctx)
	if err != nil {
		return nil, err
	}
	return exercises, nil
}

func (e *ExercisesRepository) ListForGroup(ctx context.Context, group string) ([]m.Exercise, error) {
	exercises, err := e.q.ListForGroup(ctx, group)
	if err != nil {
		return nil, err
	}
	return exercises, nil
}

func (e *ExercisesRepository) SelectExerciseDefaultById(ctx context.Context, id string) (m.ExerciseDefault, error) {
	return e.q.SelectExerciseDefaultById(ctx, id)
}

func (e *ExercisesRepository) SelectExerciseDefault(ctx context.Context, userId string, exerciseId string) (m.ExerciseDefault, error) {
	return e.q.SelectExerciseDefault(ctx, userId, exerciseId)
}

func (e *ExercisesRepository) SelectAllExerciseDefaultsForUser(ctx context.Context, userId string) ([]m.ExerciseDefault, error) {
	return e.q.SelectAllExerciseDefaultsForUser(ctx, userId)
}

func (e *ExercisesRepository) SelectAllExerciseDefaults(ctx context.Context) ([]m.ExerciseDefault, error) {
	return e.q.SelectAllExerciseDefaults(ctx)
}

func (e *ExercisesRepository) InsertExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error) {
	return e.q.InsertExerciseDefault(ctx, InsertDefaultExerciseParams{
		id:         uuid.NewString(),
		userId:     defaultExercise.UserId,
		exerciseId: defaultExercise.ExerciseId,
		weight:     defaultExercise.Weight,
		reps:       defaultExercise.Reps,
	})
}

func (e *ExercisesRepository) UpdateExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error) {
	return e.q.UpdateExerciseDefault(ctx, UpdateDefaultExerciseParams{
		id:     defaultExercise.Id,
		weight: defaultExercise.Weight,
		reps:   defaultExercise.Reps,
	})
}
