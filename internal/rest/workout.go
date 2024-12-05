package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brianfloyd/the-grid/internal"
	"github.com/go-chi/chi"
)

type WorkoutService interface {
	ById(id string) (internal.Workout, error)
	ByDate(userId string, date string) (internal.Workout, error)
	Create(userId string, workout internal.Workout) (internal.Workout, error)
	CreateSet(workoutId string, set internal.Set) (internal.Set, error)
	UpdateSet(workoutId string, setId string, set internal.Set) (internal.Set, error)
	DeleteSet(workoutId string, setId string) error
}

type WorkoutHandler struct {
	svc WorkoutService
}

func NewWorkoutHandler(svc WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		svc: svc,
	}
}

func (h *WorkoutHandler) Register(r *chi.Mux) {
	r.Get(fmt.Sprintf("/workouts/{id:%s}", uuidRegEx), h.byId)
	r.Post("/workouts", h.create)
	r.Post("/workouts/date", h.byDate)
	r.Post(fmt.Sprintf("/workouts/{id:%s}/sets", uuidRegEx), h.createSet)
	r.Patch(fmt.Sprintf("/workouts/{id:%s}/sets/{setId:%s}", uuidRegEx, uuidRegEx), h.updateSet)
}

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

type CreateWorkoutResponse struct {
	Workout Workout `json:"workout"`
}

func (h *WorkoutHandler) create(w http.ResponseWriter, r *http.Request) {
	var request CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, BadRequest, "Could not convert given body to a create workout request.", err)
		return
	}
	defer r.Body.Close()

	internalWorkout := convertCreateWorkoutRequestToInternalWorkout(request)

	workout, err := h.svc.Create(request.UserId, internalWorkout)

	if err != nil {
		renderErrorResponse(w, r, GenericError, "Creating a workout failed.", err)
		return
	}

	renderResponse(w, r,
		&CreateWorkoutResponse{
			Workout: convertInternalWorkoutToResponseWorkout(workout),
		},
		http.StatusCreated)
}

type GetWorkoutByIdResponse struct {
	Workout Workout `json:"workout"`
}

func (h *WorkoutHandler) byId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	workout, err := h.svc.ById(id)
	if err != nil {
		workoutNotFoundError := &internal.WorkoutNotFoundError{}
		if errors.As(err, &workoutNotFoundError) {
			renderErrorResponse(w, r, NotFound, "Workout was not found.", err)
		} else {
			renderErrorResponse(w, r, GenericError, "Generic workout exception.", err)
		}
		return
	}

	renderResponse(w, r,
		&GetWorkoutByIdResponse{
			Workout: convertInternalWorkoutToResponseWorkout(workout),
		},
		http.StatusOK)
}

type GetWorkoutByDateRequest struct {
	Date   string `json:"date"`
	UserId string `json:"userId"`
}

type GetWorkoutByDateResponse struct {
	Workout Workout `json:"workout"`
}

func (h *WorkoutHandler) byDate(w http.ResponseWriter, r *http.Request) {
	var request GetWorkoutByDateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, BadRequest, "Could not convert given body to a get workout by date request.", err)
		return
	}
	defer r.Body.Close()

	workout, err := h.svc.ByDate(request.UserId, request.Date)
	if err != nil {
		workoutNotFoundError := &internal.WorkoutNotFoundError{}
		if errors.As(err, &workoutNotFoundError) {
			renderErrorResponse(w, r, NotFound, "Workout was not found.", err)
		} else {
			renderErrorResponse(w, r, GenericError, "Generic workout exception.", err)
		}
		return
	}

	if err != nil {
		renderErrorResponse(w, r, GenericError, "Getting a workout by date failed.", err)
		return
	}

	renderResponse(w, r,
		&GetWorkoutByDateResponse{
			Workout: convertInternalWorkoutToResponseWorkout(workout),
		},
		http.StatusCreated)
}

type CreateSetRequest struct {
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
}

type CreateSetResponse struct {
	Set Set `json:"set"`
}

func (h *WorkoutHandler) createSet(w http.ResponseWriter, r *http.Request) {
	workoutId := chi.URLParam(r, "id")

	var request CreateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, BadRequest, "Could not convert given body to a create set request.", err)
		return
	}
	defer r.Body.Close()

	set, err := h.svc.CreateSet(workoutId, internal.Set{
		ExerciseId: request.ExerciseId,
		Reps:       request.Reps,
		Weight:     request.Weight,
	})

	if err != nil {
		renderErrorResponse(w, r, GenericError, "Creating a set failed.", err)
		return
	}

	renderResponse(w, r,
		&CreateSetResponse{
			Set: convertInternalSetToSet(set),
		},
		http.StatusCreated)
}

type UpdateSetRequest struct {
	ExerciseId string `json:"exerciseId"`
	Weight     uint64 `json:"weight"`
	Reps       uint64 `json:"reps"`
}

type UpdateSetResponse struct {
	Set Set `json:"set"`
}

func (h *WorkoutHandler) updateSet(w http.ResponseWriter, r *http.Request) {
	workoutId := chi.URLParam(r, "id")
	setId := chi.URLParam(r, "setId")

	var request UpdateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, BadRequest, "Could not convert given body to an update set request.", err)
		return
	}
	defer r.Body.Close()

	set, err := h.svc.UpdateSet(workoutId, setId, internal.Set{
		ExerciseId: request.ExerciseId,
		Reps:       request.Reps,
		Weight:     request.Weight,
	})

	if err != nil {
		renderErrorResponse(w, r, GenericError, "Updating a set failed.", err)
		return
	}

	renderResponse(w, r,
		&UpdateSetResponse{
			Set: convertInternalSetToSet(set),
		},
		http.StatusOK)
}

func (h *WorkoutHandler) deleteSet(w http.ResponseWriter, r *http.Request) {
	workoutId := chi.URLParam(r, "id")
	setId := chi.URLParam(r, "setId")

	err := h.svc.DeleteSet(workoutId, setId)
	if err != nil {
		renderErrorResponse(w, r, GenericError, "Deleting a set failed.", err)
		return
	}

	renderResponse(w, r, "ok", http.StatusOK)
}

func convertInternalSetToSet(set internal.Set) Set {
	return Set{
		Id:         set.Id,
		ExerciseId: set.ExerciseId,
		Weight:     set.Weight,
		Reps:       set.Reps,
		Count:      set.Count,
	}
}

func convertCreateWorkoutRequestToInternalWorkout(request CreateWorkoutRequest) internal.Workout {
	internalSets := make([]internal.Set, len(request.Sets))
	for idx, set := range request.Sets {
		internalSets[idx] = internal.Set{
			ExerciseId: set.ExerciseId,
			Weight:     set.Weight,
			Reps:       set.Reps,
			Count:      set.Count,
		}
	}

	return internal.Workout{
		Date: request.Date,
		Sets: internalSets,
	}
}

func convertInternalWorkoutToResponseWorkout(workout internal.Workout) Workout {
	return Workout{
		Id:         workout.Id,
		UserId:     workout.UserId,
		Date:       workout.Date,
		Sets:       convertSetToResponseSet(workout.Sets),
		ModifiedAt: workout.ModifiedAt,
		CreatedAt:  workout.CreatedAt,
	}
}

func convertSetToResponseSet(internalSets []internal.Set) []Set {
	sets := make([]Set, len(internalSets))
	for idx, set := range internalSets {
		sets[idx] = convertInternalSetToSet(set)
	}
	return sets
}
