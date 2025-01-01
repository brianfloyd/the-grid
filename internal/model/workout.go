package model

import "time"

type Workout struct {
	Id         string
	UserId     string
	Date       string
	Sets       []Set
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type Set struct {
	Id         string
	WorkoutId  string
	ExerciseId string
	Reps       uint64
	Weight     uint64
	Count      uint64
}

type ExerciseDefaults struct {
	Id         string
	UserId     string
	ExerciseId string
	Weight     uint64
	Reps       uint64
}

type GenericWorkoutError struct {
	Message string
}

type WorkoutNotFoundError struct {
	Message string
}

type WorkoutInvalidInputError struct {
	Message string
}

type WorkoutAlreadyExsitsError struct {
	Message string
}

const (
	ExerciseGroupBiceps   ExerciseGroup = "BICEP"
	ExerciseGroupBack     ExerciseGroup = "BACK"
	ExerciseGroupTricep   ExerciseGroup = "TRICEP"
	ExerciseGroupChest    ExerciseGroup = "CHEST"
	ExerciseGroupShoulder ExerciseGroup = "SHOULDER"
	ExerciseGroupLegs     ExerciseGroup = "LEGS"
	ExerciseGroupCardio   ExerciseGroup = "CARDIO"
	ExerciseGroupAbs      ExerciseGroup = "ABS"
	ExerciseGroupMisc     ExerciseGroup = "MISC"
)

var ExerciseGroupAll []ExerciseGroup = []ExerciseGroup{
	ExerciseGroupBiceps, ExerciseGroupBack, ExerciseGroupTricep, ExerciseGroupChest, ExerciseGroupShoulder,
	ExerciseGroupLegs, ExerciseGroupAbs, ExerciseGroupCardio, ExerciseGroupMisc,
}

func (e *GenericWorkoutError) Error() string {
	return e.Message
}

func (e *WorkoutNotFoundError) Error() string {
	return e.Message
}

func (e *WorkoutInvalidInputError) Error() string {
	return e.Message
}

func (e *WorkoutAlreadyExsitsError) Error() string {
	return e.Message
}
