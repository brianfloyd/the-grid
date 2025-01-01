package rest

type CreateExerciseRequest struct {
	Group string
	Name  string
}

type ExericseResponse struct {
	Id    string
	Name  string
	Group string
}

type CreateExerciseResponse struct {
	Exercise ExerciseResponse
}

type ListExercisesResponse struct {
	Exercises []ExerciseResponse
}
