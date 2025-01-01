package rest

type CreateExerciseRequest struct {
	Group string `json:"group"`
	Name  string `json:"name"`
}

type ExericseResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type CreateExerciseResponse struct {
	Exercise ExerciseResponse `json:"exercise"`
}

type ListExercisesResponse struct {
	Exercises []ExerciseResponse `json:"exercises"`
}
