package rest

import "time"

type Exercise struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Group string `json:"group"`
}

type Set struct {
	Id         string `json:"id"`
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
	Count      uint64 `json:"count"`
}

type Workout struct {
	Id         string    `json:"id"`
	UserId     string    `json:"userId"`
	Date       string    `json:"date"`
	Sets       []Set     `json:"sets"`
	CreatedAt  time.Time `json:"createdAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

// Begin request/response
type CreateWorkoutSet struct {
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
	Count      uint64 `json:"count"`
}

type CreateWorkoutRequest struct {
	UserId string             `json:"userId"`
	Date   string             `json:"date"`
	Sets   []CreateWorkoutSet `json:"sets"`
}

type GetWorkoutByIdResponse struct {
	Workout Workout `json:"workout"`
}

type CreateWorkoutResponse struct {
	Workout Workout `json:"workout"`
}

type GetWorkoutByDateRequest struct {
	Date   string `json:"date"`
	UserId string `json:"userId"`
}

type GetWorkoutByDateResponse struct {
	Workout Workout `json:"workout"`
}

type CreateSetRequest struct {
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
}

type CreateSetResponse struct {
	Set Set `json:"set"`
}

type UpdateSetRequest struct {
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
}

type UpdateSetResponse struct {
	Set Set `json:"set"`
}
