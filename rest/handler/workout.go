package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	m "github.com/brianfloyd/the-grid/internal/model"
	rm "github.com/brianfloyd/the-grid/rest/model"
	"github.com/go-chi/chi"
)

type WorkoutService interface {
	ById(ctx context.Context, id string) (m.Workout, error)
	ByDate(ctx context.Context, userId string, date string) (m.Workout, error)
	Create(ctx context.Context, userId string, workout m.Workout) (m.Workout, error)
	CreateSet(ctx context.Context, workoutId string, set m.Set) (m.Set, error)
	UpdateSet(ctx context.Context, workoutId string, setId string, set m.Set) (m.Set, error)
	DeleteSet(ctx context.Context, workoutId string, setId string) error
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
	r.Get(fmt.Sprintf("/workouts/{id:%s}", rm.UUID_REGEX), h.byId)
	r.Post("/workouts", h.create)
	r.Post("/workouts/date", h.byDate)
	r.Post(fmt.Sprintf("/workouts/{id:%s}/sets", rm.UUID_REGEX), h.createSet)
	r.Patch(fmt.Sprintf("/workouts/{id:%s}/sets/{setId:%s}", rm.UUID_REGEX, rm.UUID_REGEX), h.updateSet)
}

func (h *WorkoutHandler) create(w http.ResponseWriter, r *http.Request) {
	var request rm.CreateWorkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, rm.BadRequest, "Could not convert given body to a create workout request.", err)
		return
	}
	defer r.Body.Close()

	mWorkout := convertCreateWorkoutRequestToModelWorkout(request)

	workout, err := h.svc.Create(r.Context(), request.UserId, mWorkout)

	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "Creating a workout failed.", err)
		return
	}

	renderResponse(w, r,
		&rm.CreateWorkoutResponse{
			Workout: convertModelWorkoutToResponseWorkout(workout),
		},
		http.StatusCreated)
}

func (h *WorkoutHandler) byId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	workout, err := h.svc.ById(r.Context(), id)
	if err != nil {
		workoutNotFoundError := &m.WorkoutNotFoundError{}
		if errors.As(err, &workoutNotFoundError) {
			renderErrorResponse(w, r, rm.NotFound, "Workout was not found.", err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while looking for workout by id.", err)
		}
		return
	}

	renderResponse(w, r,
		&rm.GetWorkoutByIdResponse{
			Workout: convertModelWorkoutToResponseWorkout(workout),
		},
		http.StatusOK)
}

func (h *WorkoutHandler) byDate(w http.ResponseWriter, r *http.Request) {
	var request rm.GetWorkoutByDateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, rm.BadRequest, "Could not convert given body to a get workout by date request.", err)
		return
	}
	defer r.Body.Close()

	workout, err := h.svc.ByDate(r.Context(), request.UserId, request.Date)
	if err != nil {
		workoutNotFoundError := &m.WorkoutNotFoundError{}
		if errors.As(err, &workoutNotFoundError) {
			renderErrorResponse(w, r, rm.NotFound, "Workout was not found.", err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "Generic workout exception.", err)
		}
		return
	}

	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "Getting a workout by date failed.", err)
		return
	}

	renderResponse(w, r,
		&rm.GetWorkoutByDateResponse{
			Workout: convertModelWorkoutToResponseWorkout(workout),
		},
		http.StatusCreated)
}

func (h *WorkoutHandler) createSet(w http.ResponseWriter, r *http.Request) {
	workoutId := chi.URLParam(r, "id")

	var request rm.CreateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, rm.BadRequest, "Could not convert given body to a create set request.", err)
		return
	}
	defer r.Body.Close()

	set, err := h.svc.CreateSet(r.Context(), workoutId, m.Set{
		ExerciseId: request.ExerciseId,
		Reps:       request.Reps,
		Weight:     request.Weight,
	})

	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "Creating a set failed.", err)
		return
	}

	renderResponse(w, r,
		&rm.CreateSetResponse{
			Set: convertModelSetToResponseSet(set),
		},
		http.StatusCreated)
}

func (h *WorkoutHandler) updateSet(w http.ResponseWriter, r *http.Request) {
	workoutId := chi.URLParam(r, "id")
	setId := chi.URLParam(r, "setId")

	var request rm.UpdateSetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, rm.BadRequest, "Could not convert given body to an update set request.", err)
		return
	}
	defer r.Body.Close()

	set, err := h.svc.UpdateSet(r.Context(), workoutId, setId, m.Set{
		ExerciseId: request.ExerciseId,
		Reps:       request.Reps,
		Weight:     request.Weight,
	})

	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "Updating a set failed.", err)
		return
	}

	renderResponse(w, r,
		&rm.UpdateSetResponse{
			Set: convertModelSetToResponseSet(set),
		},
		http.StatusOK)
}

func (h *WorkoutHandler) deleteSet(w http.ResponseWriter, r *http.Request) {
	workoutId := chi.URLParam(r, "id")
	setId := chi.URLParam(r, "setId")

	err := h.svc.DeleteSet(r.Context(), workoutId, setId)
	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "Deleting a set failed.", err)
		return
	}

	renderResponse(w, r, "ok", http.StatusOK)
}

func convertModelSetToResponseSet(set m.Set) rm.SetResponse {
	return rm.SetResponse{
		Id:         set.Id,
		ExerciseId: set.ExerciseId,
		Weight:     set.Weight,
		Reps:       set.Reps,
		Count:      set.Count,
	}
}

func convertCreateWorkoutRequestToModelWorkout(request rm.CreateWorkoutRequest) m.Workout {
	mSets := make([]m.Set, len(request.Sets))
	for idx, set := range request.Sets {
		mSets[idx] = m.Set{
			ExerciseId: set.ExerciseId,
			Weight:     set.Weight,
			Reps:       set.Reps,
			Count:      set.Count,
		}
	}

	return m.Workout{
		Date: request.Date,
		Sets: mSets,
	}
}

func convertModelWorkoutToResponseWorkout(workout m.Workout) rm.WorkoutResponse {
	return rm.WorkoutResponse{
		Id:         workout.Id,
		UserId:     workout.UserId,
		Date:       workout.Date,
		Sets:       convertSetToResponseSet(workout.Sets),
		ModifiedAt: workout.ModifiedAt,
		CreatedAt:  workout.CreatedAt,
	}
}

func convertSetToResponseSet(mSets []m.Set) []rm.SetResponse {
	sets := make([]rm.SetResponse, len(mSets))
	for idx, set := range mSets {
		sets[idx] = convertModelSetToResponseSet(set)
	}
	return sets
}
