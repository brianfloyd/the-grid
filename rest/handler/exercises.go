package rest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	m "github.com/brianfloyd/the-grid/internal/model"
	rm "github.com/brianfloyd/the-grid/rest/model"
	"github.com/go-chi/chi"
)

type IExercisesService interface {
	List(ctx context.Context) ([]m.Exercise, error)
	Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error)
}

type ExercisesHandler struct {
	svc IExercisesService
}

func NewExerciseHandler(svc IExercisesService) *ExercisesHandler {
	return &ExercisesHandler{
		svc: svc,
	}
}

func (e *ExercisesHandler) Register(r *chi.Mux) {
	r.Get("/exercises", e.list)
	r.Post("/exercises", e.create)
}

func (e *ExercisesHandler) create(w http.ResponseWriter, r *http.Request) {
	var request rm.CreateExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		renderErrorResponse(w, r, rm.BadRequest, "Could not convert given body to a create exercise request.", err)
		return
	}

	defer r.Body.Close()

	exercise, err := e.svc.Create(r.Context(), m.Exercise{
		Group: m.ExerciseGroup(request.Group),
		Name:  request.Name,
	})

	if err != nil {
		var exerciseExistsError = &m.ExerciseExistsError{}
		var exerciseValidationError = &m.ExerciseValidationError{}
		if errors.As(err, &exerciseExistsError) {
			renderErrorResponse(w, r, rm.BadRequest, err.Error(), err)
		} else if errors.As(err, &exerciseValidationError) {
			renderErrorResponse(w, r, rm.BadRequest, err.Error(), err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while creating an exercise.", err)
		}
		return
	}

	renderResponse(w, r,
		&rm.CreateExerciseResponse{
			Exercise: convertModelExerciseToResponseExercise(exercise),
		},
		http.StatusCreated)
}

func (e *ExercisesHandler) list(w http.ResponseWriter, r *http.Request) {
	modelExercises, err := e.svc.List(r.Context())
	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while listing exercises.", err)
		return
	}

	exercises := make([]rm.ExerciseResponse, len(modelExercises))
	for i, exercise := range modelExercises {
		exercises[i] = convertModelExerciseToResponseExercise(exercise)
	}

	renderResponse(w, r, &rm.ListExercisesResponse{
		Exercises: exercises,
	}, http.StatusOK)
}

func convertModelExerciseToResponseExercise(exercise m.Exercise) rm.ExerciseResponse {
	return rm.ExerciseResponse{
		Id:    exercise.Id,
		Group: string(exercise.Group),
		Name:  exercise.Name,
	}
}
