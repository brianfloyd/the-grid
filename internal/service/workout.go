package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/brianfloyd/the-grid/internal/db"
	"github.com/brianfloyd/the-grid/internal/logger"
	m "github.com/brianfloyd/the-grid/internal/model"
)

type IWorkoutService interface {
	ById(ctx context.Context, id string) (m.Workout, error)
	ByDate(ctx context.Context, userId string, date string) (m.Workout, error)
	Create(ctx context.Context, userId string, workout m.Workout) (m.Workout, error)
	CreateSet(ctx context.Context, workoutId string, set m.Set) (m.Set, error)
	UpdateSet(ctx context.Context, workoutId string, setId string, set m.Set) (m.Set, error)
	DeleteSet(ctx context.Context, workoutId string, setId string) error
}

type WorkoutRespository interface {
	ByDate(ctx context.Context, userId string, date string) (m.Workout, error)
	Create(ctx context.Context, userId string, workout m.Workout) (m.Workout, error)
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

	date, err := sanitizeDate(workout.Date)
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
	return m.Workout{}, nil
}

func (w *WorkoutService) ByDate(ctx context.Context, userId string, dateString string) (m.Workout, error) {
	date, err := sanitizeDate(dateString)
	if err != nil {
		return m.Workout{}, &m.WorkoutInvalidInputError{Message: "The given date is not in the MM/DD/YYYY format."}
	}

	workout, err := w.repo.ByDate(ctx, userId, makeDateStringFromTime(date))
	if err != nil {
		workoutNotFoundError := &m.WorkoutNotFoundError{}

		if errors.As(err, &workoutNotFoundError) {
			return m.Workout{}, &m.WorkoutNotFoundError{Message: "Workout not found for the given user and date."}
		} else {
			return m.Workout{}, errors.Join(&m.GenericWorkoutError{Message: "Generic workout error."}, err)
		}
	}

	return workout, nil
}

func (w *WorkoutService) CreateSet(ctx context.Context, workoutId string, set m.Set) (m.Set, error) {
	return m.Set{}, nil
}

func (w *WorkoutService) UpdateSet(ctx context.Context, workoutId string, setId string, set m.Set) (m.Set, error) {
	return m.Set{}, nil
}

func (w *WorkoutService) DeleteSet(ctx context.Context, workoutId string, setId string) error {
	return nil
}

func (w *WorkoutService) doesWorkoutExistForDate(ctx context.Context, userId string, date time.Time) (bool, error) {
	dateStr := makeDateStringFromTime(date)

	_, err := w.repo.ByDate(ctx, userId, dateStr)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, db.ErrDbNotFound) {
		return false, nil
	}

	return false, err
}

func makeDateStringFromTime(date time.Time) string {
	return date.Format("1/2/2006")
}

func sanitizeDate(date string) (time.Time, error) {
	t, err := time.Parse("01/02/2006", date)
	if err != nil {
		return time.Time{}, errors.New("could not convert given string to the MM/DD/YYYY time format")
	}
	return t, nil
}
