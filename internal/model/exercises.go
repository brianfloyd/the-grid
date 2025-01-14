package model

type ExerciseGroup string

type Exercise struct {
	Id    string
	Name  string
	Group ExerciseGroup
}

type ExerciseDefault struct {
	Id         string
	ExerciseId string
	UserId     string
	Weight     uint64
	Reps       uint64
}

type GenericExerciseError struct {
	Message string
}

type ExerciseValidationError struct {
	Message string
}

type ExerciseExistsError struct {
	Message string
}

type ExerciseDefaultValidationError struct {
	Message string
}

type ExerciseDefaultExistsError struct {
	Message string
}

type ExerciseDefaultNotFoundError struct {
	Message string
}

func (e *GenericExerciseError) Error() string {
	return e.Message
}

func (e *ExerciseValidationError) Error() string {
	return e.Message
}

func (e *ExerciseExistsError) Error() string {
	return e.Message
}

func (e *ExerciseDefaultValidationError) Error() string {
	return e.Message
}

func (e *ExerciseDefaultExistsError) Error() string {
	return e.Message
}

func (e *ExerciseDefaultNotFoundError) Error() string {
	return e.Message
}
