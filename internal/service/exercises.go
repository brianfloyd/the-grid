package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/brianfloyd/the-grid/internal/db"
	"github.com/brianfloyd/the-grid/internal/logger"
	m "github.com/brianfloyd/the-grid/internal/model"
)

type IExercisesService interface {
	List(ctx context.Context) ([]m.Exercise, error)
	ListForGroup(ctx context.Context, group string) ([]m.Exercise, error)
	Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error)

	GetExerciseDefault(ctx context.Context, userId string, exerciseId string) (m.ExerciseDefault, error)
	CreateExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error)
	UpdateExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error)
	ListExerciseDefaultsForUser(ctx context.Context, userId string) ([]m.ExerciseDefault, error)
	ListAllExerciseDefaults(ctx context.Context) ([]m.ExerciseDefault, error)
}

type IExerciseRepository interface {
	List(ctx context.Context) ([]m.Exercise, error)
	ListForGroup(ctx context.Context, group string) ([]m.Exercise, error)
	Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error)

	SelectExerciseDefaultById(ctx context.Context, id string) (m.ExerciseDefault, error)
	SelectExerciseDefault(ctx context.Context, userId string, exerciseId string) (m.ExerciseDefault, error)
	InsertExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error)
	UpdateExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error)
	SelectAllExerciseDefaultsForUser(ctx context.Context, userId string) ([]m.ExerciseDefault, error)
	SelectAllExerciseDefaults(ctx context.Context) ([]m.ExerciseDefault, error)
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
		return m.Exercise{}, &m.ExerciseValidationError{Message: strings.Join(issues, " ")}
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

func (e *ExerciseService) GetExerciseDefault(ctx context.Context, userId string, exerciseId string) (m.ExerciseDefault, error) {
	logger.InfoArgs(ctx, "Getting default exercise (%s) for user (%s)", exerciseId, userId)

	exerciseDefault, err := e.repo.SelectExerciseDefault(ctx, userId, exerciseId)
	if err != nil {
		if errors.Is(err, db.ErrDbNotFound) {
			return m.ExerciseDefault{}, &m.ExerciseDefaultNotFoundError{Message: fmt.Sprintf("Exercise default does not exist for user %s and exercise id %s.", userId, exerciseId)}
		} else {
			return m.ExerciseDefault{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected error occurred while getting default exercise for user."}, err)
		}
	}

	return exerciseDefault, nil
}

func (e *ExerciseService) ListExerciseDefaultsForUser(ctx context.Context, userId string) ([]m.ExerciseDefault, error) {
	logger.InfoArgs(ctx, "Listing default exercises for user (%s).", userId)

	exerciseDefaults, err := e.repo.SelectAllExerciseDefaultsForUser(ctx, userId)
	if err != nil {
		return []m.ExerciseDefault{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected error occurred while listing all exercise defaults for user."}, err)
	}

	return exerciseDefaults, nil
}

func (e *ExerciseService) ListAllExerciseDefaults(ctx context.Context) ([]m.ExerciseDefault, error) {
	logger.Info(ctx, "Listing all default exercises.")

	exerciseDefaults, err := e.repo.SelectAllExerciseDefaults(ctx)
	if err != nil {
		return []m.ExerciseDefault{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected error occurred while listing all default exercises."}, err)
	}
	return exerciseDefaults, nil
}

func (e *ExerciseService) CreateExerciseDefault(ctx context.Context, exerciseDefault m.ExerciseDefault) (m.ExerciseDefault, error) {
	logger.InfoArgs(ctx, "Creating exercise default (%v)", exerciseDefault)

	if issues, ok := e.validateExericseDefault(ctx, exerciseDefault, false); !ok {
		return m.ExerciseDefault{}, &m.ExerciseDefaultValidationError{Message: strings.Join(issues, " ")}
	}

	exists, err := e.doesExerciseDefaultExist(ctx, exerciseDefault)
	if err != nil {
		return m.ExerciseDefault{}, errors.Join(&m.GenericExerciseError{Message: "Could not verify if the exercise default already exists."}, err)
	}

	if exists {
		return m.ExerciseDefault{}, &m.ExerciseDefaultExistsError{Message: fmt.Sprintf("Exercise default for user '%s' and exercise id '%s' already exists.", exerciseDefault.UserId, exerciseDefault.ExerciseId)}
	}

	exerciseDefault, err = e.repo.InsertExerciseDefault(ctx, exerciseDefault)
	if err != nil {
		return m.ExerciseDefault{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected error occurred while creating an exercise default."}, err)
	}

	return exerciseDefault, nil
}

func (e *ExerciseService) UpdateExerciseDefault(ctx context.Context, exerciseDefault m.ExerciseDefault) (m.ExerciseDefault, error) {
	logger.InfoArgs(ctx, "Updating exercise default (%v).", exerciseDefault)

	if issues, ok := e.validateExericseDefault(ctx, exerciseDefault, true); !ok {
		return m.ExerciseDefault{}, &m.ExerciseDefaultValidationError{Message: strings.Join(issues, " ")}
	}

	existing, err := e.repo.SelectExerciseDefaultById(ctx, exerciseDefault.Id)
	if err != nil {
		if errors.Is(err, db.ErrDbNotFound) {
			return m.ExerciseDefault{}, &m.ExerciseDefaultNotFoundError{Message: fmt.Sprintf("Cannot update exercise defaults with id '%s' becuase it does not exist.", exerciseDefault.Id)}
		} else {
			return m.ExerciseDefault{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected error occurred while updating an exercise default."}, err)
		}
	}

	// Make sure these are immutable.
	exerciseDefault.Id = existing.Id
	exerciseDefault.ExerciseId = existing.ExerciseId
	exerciseDefault.UserId = existing.UserId

	exerciseDefault, err = e.repo.UpdateExerciseDefault(ctx, exerciseDefault)
	if err != nil {
		return m.ExerciseDefault{}, errors.Join(&m.GenericExerciseError{Message: "An unexpected error occurred while updating an exercise default."}, err)
	}

	return exerciseDefault, nil
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

func (e *ExerciseService) validateExericseDefault(ctx context.Context, exerciseDefault m.ExerciseDefault, needsId bool) ([]string, bool) {
	issues := []string{}

	if exerciseDefault.Reps <= 0 {
		logger.InfoArgs(ctx, "Exercise default (%v) is invalid, reps must be a non-zero positive number.", exerciseDefault)
		issues = append(issues, "Reps must be a non-zero positive number.")
	}

	if exerciseDefault.Weight <= 0 {
		logger.InfoArgs(ctx, "Exercise default (%v) is invalid, weight must be a non-zero positive number.", exerciseDefault)
		issues = append(issues, "Weight must be a non-zero positive number.")
	}

	if needsId && len(exerciseDefault.Id) < m.UUID_LEN {
		logger.InfoArgs(ctx, "Exercise default (%v) is ivnalid, id must be specified.", exerciseDefault)
		issues = append(issues, "Id must be specified.")
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

func (e *ExerciseService) doesExerciseDefaultExist(ctx context.Context, exerciseDefault m.ExerciseDefault) (bool, error) {
	if exerciseDefault.Id != "" {
		_, err := e.repo.SelectExerciseDefaultById(ctx, exerciseDefault.Id)
		if err != nil {
			if !errors.Is(err, db.ErrDbNotFound) {
				// An error occurred that wasn't not found. Report that. If it wasn't found, try another method.
				return true, err
			}
		}
	}

	if exerciseDefault.UserId != "" && exerciseDefault.ExerciseId != "" {
		_, err := e.repo.SelectExerciseDefault(ctx, exerciseDefault.UserId, exerciseDefault.ExerciseId)
		if err != nil {
			if errors.Is(err, db.ErrDbNotFound) {
				return false, nil
			} else {
				return true, err
			}
		}
		return true, nil
	}

	return false, nil
}
