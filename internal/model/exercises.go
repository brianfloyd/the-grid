package model

type ExerciseGroup string

type Exercise struct {
	Id    string
	Name  string
	Group ExerciseGroup
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

func (e *GenericExerciseError) Error() string {
	return e.Message
}

func (e *ExerciseValidationError) Error() string {
	return e.Message
}

func (e *ExerciseExistsError) Error() string {
	return e.Message
}
