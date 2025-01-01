package service

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	m "github.com/brianfloyd/the-grid/internal/model"
)

type IExercisesService interface {
	List() ([]m.Exercise, error)
	ListForGroup(group string) ([]m.Exercise, error)
	Create(exercise m.Exercise) (m.Exercise, error)
}

type IExerciseRepository interface {
	List() ([]m.Exercise, error)
	ListForGroup(group string) ([]m.Exercise, error)
	Create(exercise m.Exercise) (m.Exercise, error)
}

type ExerciseService struct {
	repo IExerciseRepository
}

func NewExerciseService(repo IExerciseRepository) *ExerciseService {
	return &ExerciseService{
		repo: repo,
	}
}

func (e *ExerciseService) List() ([]m.Exercise, error) {
	exercises, err := e.repo.List()
	if err != nil {
		return nil, errors.Join(&m.GenericExerciseError{Message: "An error occurred while listing exercises."}, err)
	}
	return exercises, nil
}

func (e *ExerciseService) ListForGroup(group string) ([]m.Exercise, error) {
	exercises, err := e.repo.ListForGroup(group)
	if err != nil {
		return nil, errors.Join(&m.GenericExerciseError{Message: fmt.Sprintf("An error occurred while listing exercises for group %s.", group)}, err)
	}
	return exercises, nil
}

func (e *ExerciseService) Create(exercise m.Exercise) (m.Exercise, error) {
	if issues, ok := e.validate(exercise); !ok {
		return m.Exercise{}, &m.ExerciseValidationError{Message: strings.Join(issues, "\n")}
	}

	exists, err := e.doesExerciseExist(exercise)
	if err != nil {
		return m.Exercise{}, errors.Join(&m.GenericExerciseError{Message: "Could not verify if the exercise already exists."}, err)
	}

	if exists {
		return m.Exercise{}, &m.ExerciseExistsError{Message: fmt.Sprintf("Exercise in the group (%s) with name (%s) already exists.", exercise.Group, exercise.Name)}
	}

	createdExercise, err := e.repo.Create(exercise)
	if err != nil {
		return m.Exercise{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected exception occurred while creating an exercise."}, err)
	}

	return createdExercise, nil
}

func (e *ExerciseService) validate(exercise m.Exercise) ([]string, bool) {
	issues := []string{}
	if len(exercise.Name) < 5 {
		issues = append(issues, "Exercise name must be at least five characters.")
	}

	if !slices.Contains(m.ExerciseGroupAll, exercise.Group) {
		issues = append(issues, "Exercise group does not exist.")
	}

	return issues, len(issues) == 0
}

func (e *ExerciseService) doesExerciseExist(exercise m.Exercise) (bool, error) {
	exercises, err := e.List()
	if err != nil {
		return false, err
	}

	for _, e := range exercises {
		if e.Group == exercise.Group && e.Name == e.Name {
			return true, nil
		}
	}

	return false, nil
}
