package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/brianfloyd/the-grid/internal/logger"
	m "github.com/brianfloyd/the-grid/internal/model"
)

type IExercisesService interface {
	List(ctx context.Context) ([]m.Exercise, error)
	ListForGroup(ctx context.Context, group string) ([]m.Exercise, error)
	Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error)
}

type IExerciseRepository interface {
	List(ctx context.Context) ([]m.Exercise, error)
	ListForGroup(ctx context.Context, group string) ([]m.Exercise, error)
	Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error)
}

type ExerciseService struct {
	repo IExerciseRepository
}

func NewExerciseService(repo IExerciseRepository) *ExerciseService {
	return &ExerciseService{
		repo: repo,
	}
}

func (e *ExerciseService) List(ctx context.Context) ([]m.Exercise, error) {
	logger.Trace(ctx, "Listing exercises.")
	exercises, err := e.repo.List(ctx)
	if err != nil {
		return nil, errors.Join(&m.GenericExerciseError{Message: "An error occurred while listing exercises."}, err)
	}
	return exercises, nil
}

func (e *ExerciseService) ListForGroup(ctx context.Context, group string) ([]m.Exercise, error) {
	logger.TraceArgs(ctx, "Listing exercises for group (%s).", group)
	exercises, err := e.repo.ListForGroup(ctx, group)
	if err != nil {
		return nil, errors.Join(&m.GenericExerciseError{Message: fmt.Sprintf("An error occurred while listing exercises for group %s.", group)}, err)
	}
	return exercises, nil
}

func (e *ExerciseService) Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error) {
	logger.InfoArgs(ctx, "Creating exercise (%v)", exercise)
	if issues, ok := e.validate(ctx, exercise); !ok {
		return m.Exercise{}, &m.ExerciseValidationError{Message: strings.Join(issues, "\n")}
	}

	exists, err := e.doesExerciseExist(ctx, exercise)
	if err != nil {
		return m.Exercise{}, errors.Join(&m.GenericExerciseError{Message: "Could not verify if the exercise already exists."}, err)
	}

	if exists {
		return m.Exercise{}, &m.ExerciseExistsError{Message: fmt.Sprintf("Exercise in the group %s with name %s already exists.", exercise.Group, exercise.Name)}
	}

	createdExercise, err := e.repo.Create(ctx, exercise)
	if err != nil {
		return m.Exercise{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected exception occurred while creating an exercise."}, err)
	}

	return createdExercise, nil
}

func (e *ExerciseService) validate(ctx context.Context, exercise m.Exercise) ([]string, bool) {
	issues := []string{}
	if len(exercise.Name) < 5 {
		logger.InfoArgs(ctx, "Exercise (%v) is invalid, the name must be at least five characters long.", exercise)
		issues = append(issues, "Exercise name must be at least five characters.")
	}

	if !slices.Contains(m.ExerciseGroupAll, exercise.Group) {
		logger.InfoArgs(ctx, "Exercise (%v) is invalid, the exercise group does not exist.", exercise)
		issues = append(issues, "Exercise group does not exist.")
	}

	return issues, len(issues) == 0
}

func (e *ExerciseService) doesExerciseExist(ctx context.Context, exercise m.Exercise) (bool, error) {
	exercises, err := e.List(ctx)
	if err != nil {
		return false, err
	}

	for _, e := range exercises {
		if e.Group == exercise.Group && e.Name == exercise.Name {
			return true, nil
		}
	}

	return false, nil
}
