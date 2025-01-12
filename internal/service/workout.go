package service

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/brianfloyd/the-grid/internal/db"
	"github.com/brianfloyd/the-grid/internal/logger"
	m "github.com/brianfloyd/the-grid/internal/model"
	"github.com/brianfloyd/the-grid/util"
)

type IWorkoutService interface {
	ById(ctx context.Context, id string) (m.Workout, error)
	ByDate(ctx context.Context, userId string, date string) (m.Workout, error)
	Create(ctx context.Context, userId string, workout m.Workout) (m.Workout, error)
	CreateSet(ctx context.Context, workoutId string, set m.Set) (m.Set, error)
	UpdateSet(ctx context.Context, workoutId string, setId string, set m.Set) (m.Set, error)
	DeleteSet(ctx context.Context, workoutId string, setId string) error
	DeleteSets(ctx context.Context, workoutId string, setIds []string) error
}

type WorkoutRespository interface {
	ById(ctx context.Context, id string) (m.Workout, error)
	ByDate(ctx context.Context, userId string, date string) (m.Workout, error)
	Create(ctx context.Context, userId string, workout m.Workout) (m.Workout, error)
	CreateSet(ctx context.Context, workoutId string, set m.Set) (m.Set, error)
	UpdateSet(ctx context.Context, set m.Set) (m.Set, error)
	DeleteSets(ctx context.Context, setIds []string) error
}

type WorkoutService struct {
	repo    WorkoutRespository
	userSvc IUserService
}

func NewWorkoutService(repo WorkoutRespository, userSvc IUserService) *WorkoutService {
	return &WorkoutService{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (w *WorkoutService) Create(ctx context.Context, userId string, workout m.Workout) (m.Workout, error) {
	logger.InfoArgs(ctx, "Creating workout (%v) for user (%s).", workout, userId)

	date, err := util.SanitizeDate(workout.Date)
	if err != nil {
		return m.Workout{}, &m.WorkoutInvalidInputError{Message: "The given date is not in the MM/DD/YYYY format."}
	}

	exists, err := w.doesWorkoutExistForDate(ctx, userId, date)
	if err != nil {
		return m.Workout{}, errors.Join(&m.GenericWorkoutError{
			Message: "Could not verify if the workout already exists for the given date.",
		}, err)
	}

	if exists {
		return m.Workout{}, &m.WorkoutAlreadyExsitsError{Message: fmt.Sprintf("Workout with the date %s already exists.", workout.Date)}
	}

	_, err = w.userSvc.ById(ctx, userId)
	if err != nil {
		var userNotFoundError = &m.UserNotFoundError{}

		if errors.As(err, &userNotFoundError) {
			return m.Workout{}, &m.WorkoutInvalidInputError{Message: fmt.Sprintf("The user with id %s does not exist.", userId)}
		} else {
			return m.Workout{}, &m.GenericWorkoutError{Message: "Could not verify if the user exists."}
		}
	}

	// TODO: For each set, populate the defaults if there is no given value.

	createdWorkout, err := w.repo.Create(ctx, userId, workout)
	if err != nil {
		return m.Workout{}, errors.Join(&m.GenericWorkoutError{Message: "An unexpected exception occurred while creating a workout."}, err)
	}

	return createdWorkout, nil
}

func (w *WorkoutService) ById(ctx context.Context, id string) (m.Workout, error) {
	logger.TraceArgs(ctx, "Getting workout by id (%s).", id)

	workout, err := w.repo.ById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrDbNotFound) {
			return m.Workout{}, &m.WorkoutNotFoundError{Message: fmt.Sprintf("Could not find workout by id '%s'.", id)}
		} else {
			message := fmt.Sprintf("An unexpected error occurred while getting workout by id (%s).", id)
			return m.Workout{}, errors.Join(&m.GenericWorkoutError{Message: message}, err)
		}
	}

	return workout, nil
}

func (w *WorkoutService) ByDate(ctx context.Context, userId string, dateString string) (m.Workout, error) {
	date, err := util.SanitizeDate(dateString)
	if err != nil {
		return m.Workout{}, &m.WorkoutInvalidInputError{Message: "The given date is not in the correct format."}
	}

	workout, err := w.repo.ByDate(ctx, userId, util.MakeDateStringFromTime(date))
	if err != nil {
		if errors.Is(err, db.ErrDbNotFound) {
			return m.Workout{}, &m.WorkoutNotFoundError{Message: "Workout not found for the given user and date."}
		} else {
			return m.Workout{}, errors.Join(&m.GenericWorkoutError{Message: "Generic workout error."}, err)
		}
	}

	return workout, nil
}

func (w *WorkoutService) CreateSet(ctx context.Context, workoutId string, set m.Set) (m.Set, error) {
	_, err := w.ById(ctx, workoutId)
	if err != nil {
		return m.Set{}, err
	}

	set, err = w.repo.CreateSet(ctx, workoutId, set)
	if err != nil {
		return m.Set{}, errors.Join(&m.GenericWorkoutError{Message: "An unexpected error occurred while creating a set."}, err)
	}

	return set, nil
}

func (w *WorkoutService) UpdateSet(ctx context.Context, workoutId string, setId string, set m.Set) (m.Set, error) {
	workout, err := w.ById(ctx, workoutId)
	if err != nil {
		return m.Set{}, err
	}

	var targetSet *m.Set
	for i := 0; i < len(workout.Sets); i++ {
		if workout.Sets[i].Id == setId {
			targetSet = &workout.Sets[i]
			break
		}
	}

	if targetSet == nil {
		return m.Set{}, &m.WorkoutInvalidInputError{Message: fmt.Sprintf("Workout '%s' does not contain set '%s'.", workoutId, setId)}
	}

	set.Id = targetSet.Id
	set.WorkoutId = targetSet.WorkoutId
	if set.ExerciseId == "" {
		set.ExerciseId = targetSet.ExerciseId
	}
	if set.Reps == 0 {
		set.Reps = targetSet.Reps
	}
	if set.Weight == 0 {
		set.Weight = targetSet.Weight
	}
	if set.Count == 0 {
		set.Count = targetSet.Count
	}

	set, err = w.repo.UpdateSet(ctx, set)
	if err != nil {
		return m.Set{}, errors.Join(&m.GenericWorkoutError{Message: "An unexpected error occurred while updating a set."}, err)
	}

	return set, nil
}

func (w *WorkoutService) DeleteSets(ctx context.Context, workoutId string, setIds []string) error {
	logger.InfoArgs(ctx, "Deleting sets '%v' from workout '%s'.", setIds, workoutId)

	workout, err := w.ById(ctx, workoutId)
	if err != nil {
		return err
	}
	workoutSetIds := make([]string, len(workout.Sets))
	for i, set := range workout.Sets {
		workoutSetIds[i] = set.Id
	}

	// TODO: Create a reponse object that details which ones were successfully deleted.
	targetSets := []string{}
	for _, setId := range setIds {
		if slices.Contains(workoutSetIds, setId) {
			targetSets = append(targetSets, setId)
		} else {
			logger.WarnArgs(ctx, "Requested set to be deleted (%s) was not a part of workout (%s).", setId, workoutId)
		}
	}

	if len(targetSets) > 0 {
		err = w.repo.DeleteSets(ctx, targetSets)
		if err != nil {
			return errors.Join(&m.GenericWorkoutError{Message: "An unexpected error occurred while deleting a set."}, err)
		}
	}

	return nil
}

func (w *WorkoutService) DeleteSet(ctx context.Context, workoutId string, setId string) error {
	logger.InfoArgs(ctx, "Deleting set '%s' from workout '%s'.", setId, workoutId)
	return w.DeleteSets(ctx, workoutId, []string{setId})
}

func (w *WorkoutService) doesWorkoutExistForDate(ctx context.Context, userId string, date time.Time) (bool, error) {
	dateStr := util.MakeDateStringFromTime(date)

	_, err := w.repo.ByDate(ctx, userId, dateStr)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, db.ErrDbNotFound) {
		return false, nil
	}

	return false, err
}
