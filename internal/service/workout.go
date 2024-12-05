package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/brianfloyd/the-grid/internal"
	"github.com/brianfloyd/the-grid/internal/db"
)

type WorkoutService interface {
	ById(id string) (internal.Workout, error)
	ByDate(userId string, date string) (internal.Workout, error)
	Create(userId string, workout internal.Workout) (internal.Workout, error)
	CreateSet(workoutId string, set internal.Set) (internal.Set, error)
	UpdateSet(workoutId string, setId string, set internal.Set) (internal.Set, error)
	DeleteSet(workoutId string, setId string) error
}

type WorkoutRespository interface {
	ByDate(userId string, date string) (internal.Workout, error)
	Create(userId string, workout internal.Workout) (internal.Workout, error)
}

type Workout struct {
	repo    WorkoutRespository
	userSvc UserService
}

func NewWorkout(repo WorkoutRespository, userSvc UserService) *Workout {
	return &Workout{
		repo:    repo,
		userSvc: userSvc,
	}
}

func (w *Workout) Create(userId string, workout internal.Workout) (internal.Workout, error) {
	date, err := sanitizeDate(workout.Date)
	if err != nil {
		return internal.Workout{}, &internal.WorkoutInvalidInputError{Message: "The given date is not in the MM/DD/YYYY format."}
	}

	exists, err := w.doesWorkoutExistForDate(userId, date)
	if err != nil {
		return internal.Workout{}, errors.Join(&internal.GenericWorkoutError{
			Message: "Could not verify if the workout already exists for the given date.",
		}, err)
	}

	if exists {
		return internal.Workout{}, &internal.WorkoutAlreadyExsitsError{Message: fmt.Sprintf("Workout with the date %s already exists.", workout.Date)}
	}

	_, err = w.userSvc.ById(userId)
	if err != nil {
		var userNotFoundError = &internal.UserNotFoundError{}

		if errors.As(err, &userNotFoundError) {
			return internal.Workout{}, &internal.WorkoutInvalidInputError{Message: fmt.Sprintf("The user with id %s does not exist.", userId)}
		} else {
			return internal.Workout{}, &internal.GenericWorkoutError{Message: "Could not verify if the user exists."}
		}
	}

	// TODO: For each set, populate the defaults if there is no given value.

	createdWorkout, err := w.repo.Create(userId, workout)
	if err != nil {
		return internal.Workout{}, errors.Join(&internal.GenericWorkoutError{Message: "An unexpected exception occurred while creating a workout."}, err)
	}

	return createdWorkout, nil
}

func (w *Workout) ById(id string) (internal.Workout, error) {
	return internal.Workout{}, nil
}

func (w *Workout) ByDate(userId string, dateString string) (internal.Workout, error) {
	date, err := sanitizeDate(dateString)
	if err != nil {
		return internal.Workout{}, &internal.WorkoutInvalidInputError{Message: "The given date is not in the MM/DD/YYYY format."}
	}

	workout, err := w.repo.ByDate(userId, makeDateStringFromTime(date))
	if err != nil {
		workoutNotFoundError := &internal.WorkoutNotFoundError{}

		if errors.As(err, &workoutNotFoundError) {
			return internal.Workout{}, &internal.WorkoutNotFoundError{Message: "Workout not found for the given user and date."}
		} else {
			return internal.Workout{}, errors.Join(&internal.GenericWorkoutError{Message: "Generic workout error."}, err)
		}
	}

	return workout, nil
}

func (w *Workout) CreateSet(workoutId string, set internal.Set) (internal.Set, error) {
	return internal.Set{}, nil
}

func (w *Workout) UpdateSet(workoutId string, setId string, set internal.Set) (internal.Set, error) {
	return internal.Set{}, nil
}

func (w *Workout) DeleteSet(workoutId string, setId string) error {
	return nil
}

func (w *Workout) doesWorkoutExistForDate(userId string, date time.Time) (bool, error) {
	dateStr := makeDateStringFromTime(date)

	_, err := w.repo.ByDate(userId, dateStr)

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
