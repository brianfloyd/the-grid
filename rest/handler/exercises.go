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

type IExercisesService interface {
	List(ctx context.Context) ([]m.Exercise, error)
	Create(ctx context.Context, exercise m.Exercise) (m.Exercise, error)

	GetExerciseDefault(ctx context.Context, userId string, exerciseId string) (m.ExerciseDefault, error)
	CreateExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error)
	UpdateExerciseDefault(ctx context.Context, defaultExercise m.ExerciseDefault) (m.ExerciseDefault, error)
	ListExerciseDefaultsForUser(ctx context.Context, userId string) ([]m.ExerciseDefault, error)
	ListAllExerciseDefaults(ctx context.Context) ([]m.ExerciseDefault, error)
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

	r.Post("/exercise-defaults", e.createDefaultExercise)
	r.Patch(fmt.Sprintf("/exercise-defaults/{id:%s}", rm.UUID_REGEX), e.updateDefaultExercise)
	r.Get(fmt.Sprintf("/exercise-defaults/user/{userId:%s}", rm.UUID_REGEX), e.listDefaultExercisesForUser)
	r.Get("/exercise-defaults", e.listAllDefaultExercises)
	r.Get(fmt.Sprintf("/exercise-defaults/user/{userId:%s}/exercise/{exerciseId:%s}", rm.UUID_REGEX, rm.UUID_REGEX), e.getDefaultExerciseForUser)
}

func (e *ExercisesHandler) createDefaultExercise(w http.ResponseWriter, r *http.Request) {
	var request rm.CreateExerciseDefaultRequest
	if ok := decodeBody(w, r, &request, "create exercise default"); !ok {
		return
	}

	exerciseDefault, err := e.svc.CreateExerciseDefault(r.Context(), m.ExerciseDefault{
		UserId:     request.UserId,
		ExerciseId: request.ExerciseId,
		Weight:     request.Reps,
		Reps:       request.Weight,
	})

	if err != nil {
		var exerciseDefaultValidationError = &m.ExerciseDefaultValidationError{}
		var exerciseDefaultExistsError = &m.ExerciseDefaultExistsError{}
		if errors.As(err, &exerciseDefaultValidationError) {
			renderErrorResponse(w, r, rm.BadRequest, err.Error(), err)
		} else if errors.As(err, &exerciseDefaultExistsError) {
			renderErrorResponse(w, r, rm.BadRequest, err.Error(), err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while creating an exercise default.", err)
		}
		return
	}

	renderResponse(w, r,
		&rm.CreateExerciseDefaultResponse{
			ExerciseDefault: toResponseExerciseDefault(exerciseDefault),
		},
		http.StatusCreated,
	)
}

func (e *ExercisesHandler) updateDefaultExercise(w http.ResponseWriter, r *http.Request) {
	var request rm.UpdateExerciseDefaultRequest
	if ok := decodeBody(w, r, &request, "update exercise default"); !ok {
		return
	}
	id := chi.URLParam(r, "id")

	exerciseDefault, err := e.svc.UpdateExerciseDefault(r.Context(), m.ExerciseDefault{
		Id:     id,
		Weight: request.Weight,
		Reps:   request.Reps,
	})

	if err != nil {
		var exerciseDefaultValidationError = &m.ExerciseDefaultValidationError{}
		if errors.As(err, &exerciseDefaultValidationError) {
			renderErrorResponse(w, r, rm.BadRequest, err.Error(), err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while updating an exercise default.", err)
		}
		return
	}

	renderResponse(w, r,
		&rm.UpdateExerciseDefaultResponse{
			ExerciseDefault: toResponseExerciseDefault(exerciseDefault),
		},
		http.StatusOK,
	)
}

func (e *ExercisesHandler) listAllDefaultExercises(w http.ResponseWriter, r *http.Request) {
	exerciseDefaults, err := e.svc.ListAllExerciseDefaults(r.Context())
	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while listing all exercise defaults.", err)
		return
	}

	renderResponse(w, r,
		&rm.ListExeciseDefaultsResposne{
			ExerciseDefaults: toResponseExerciseDefaults(exerciseDefaults),
		},
		http.StatusOK,
	)
}

func (e *ExercisesHandler) listDefaultExercisesForUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")

	exerciseDefaults, err := e.svc.ListExerciseDefaultsForUser(r.Context(), userId)
	if err != nil {
		renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while listing all exercise defaults for the user.", err)
		return
	}

	renderResponse(w, r,
		&rm.ListExeciseDefaultsResposne{
			ExerciseDefaults: toResponseExerciseDefaults(exerciseDefaults),
		},
		http.StatusOK,
	)
}

func (e *ExercisesHandler) getDefaultExerciseForUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	exerciseId := chi.URLParam(r, "exerciseId")

	exerciseDefault, err := e.svc.GetExerciseDefault(r.Context(), userId, exerciseId)
	if err != nil {
		var exerciseDefaultNotFoundError = &m.ExerciseDefaultNotFoundError{}
		if errors.As(err, &exerciseDefaultNotFoundError) {
			renderErrorResponse(w, r, rm.NotFound, err.Error(), err)
		} else {
			renderErrorResponse(w, r, rm.GenericError, "An unexpected error occurred while getting an exercise default for the user.", err)
		}
		return
	}

	renderResponse(w, r,
		&rm.GetExerciseDefaultResponse{
			ExerciseDefault: toResponseExerciseDefault(exerciseDefault),
		},
		http.StatusOK,
	)
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
			Exercise: toResponseExercise(exercise),
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
		exercises[i] = toResponseExercise(exercise)
	}

	renderResponse(w, r, &rm.ListExercisesResponse{
		Exercises: exercises,
	}, http.StatusOK)
}

func toResponseExercise(exercise m.Exercise) rm.ExerciseResponse {
	return rm.ExerciseResponse{
		Id:    exercise.Id,
		Group: string(exercise.Group),
		Name:  exercise.Name,
	}
}

func toResponseExerciseDefault(exerceiseDefault m.ExerciseDefault) rm.ExerciseDefaultResponse {
	return rm.ExerciseDefaultResponse{
		Id:         exerceiseDefault.Id,
		UserId:     exerceiseDefault.UserId,
		ExerciseId: exerceiseDefault.ExerciseId,
		Weight:     exerceiseDefault.Weight,
		Reps:       exerceiseDefault.Reps,
	}
}

func toResponseExerciseDefaults(defaults []m.ExerciseDefault) []rm.ExerciseDefaultResponse {
	r := make([]rm.ExerciseDefaultResponse, len(defaults))
	for i, def := range defaults {
		r[i] = toResponseExerciseDefault(def)
	}
	return r
}
