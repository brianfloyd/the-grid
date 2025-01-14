package rest

type CreateExerciseRequest struct {
	Group string `json:"group"`
	Name  string `json:"name"`
}

type ExerciseResponse struct {
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

type CreateExerciseDefaultRequest struct {
	UserId     string `json:"userId"`
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
}

type UpdateExerciseDefaultRequest struct {
	Weight uint64 `json:"weight"`
	Reps   uint64 `json:"reps"`
}

type ExerciseDefaultResponse struct {
	Id         string `json:"id"`
	UserId     string `json:"userId"`
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
}

type CreateExerciseDefaultResponse struct {
	ExerciseDefault ExerciseDefaultResponse `json:"exerciseDefault"`
}

type UpdateExerciseDefaultResponse struct {
	ExerciseDefault ExerciseDefaultResponse `json:"exerciseDefault"`
}

type ListExeciseDefaultsResposne struct {
	ExerciseDefaults []ExerciseDefaultResponse `json:"exerciseDefaults"`
}

type GetExerciseDefaultResponse struct {
	ExerciseDefault ExerciseDefaultResponse `json:"exerciseDefault"`
}
