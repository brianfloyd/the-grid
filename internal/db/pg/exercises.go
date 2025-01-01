package pg

import (
	"context"

	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/google/uuid"
)

type ExercisesQueries struct {
	db DbConnection
}

func NewExercisesQueries(conn DbConnection) *ExercisesQueries {
	return &ExercisesQueries{
		db: conn,
	}
}

type ExercisesRepository struct {
	q *ExercisesQueries
}

func NewExercisesRepository(conn DbConnection) *ExercisesRepository {
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
