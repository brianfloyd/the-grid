package rest

import "time"

type SetResponse struct {
	Id         string `json:"id"`
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
	Count      uint64 `json:"count"`
}

type WorkoutResponse struct {
	Id         string        `json:"id"`
	UserId     string        `json:"userId"`
	Date       string        `json:"date"`
	Sets       []SetResponse `json:"sets"`
	CreatedAt  time.Time     `json:"createdAt"`
	ModifiedAt time.Time     `json:"modifiedAt"`
}

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
	Workout WorkoutResponse `json:"workout"`
}

type CreateWorkoutResponse struct {
	Workout WorkoutResponse `json:"workout"`
}

type GetWorkoutByDateRequest struct {
	Date   string `json:"date"`
	UserId string `json:"userId"`
}

type GetWorkoutByDateResponse struct {
	Workout WorkoutResponse `json:"workout"`
}

type CreateSetRequest struct {
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
}

type CreateSetResponse struct {
	Set SetResponse `json:"set"`
}

type UpdateSetRequest struct {
	ExerciseId string `json:"exerciseId,omitempty"`
	Weight     uint64 `json:"weight,omitempty"`
	Reps       uint64 `json:"reps,omitempty"`
}

type UpdateSetResponse struct {
	Set SetResponse `json:"set"`
}

type DeleteSetResponse struct {
	WorkoutId string `json:"workoutId"`
	SetId     string `json:"setId"`
}
