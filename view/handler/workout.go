package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi"

	"github.com/brianfloyd/the-grid/internal/logger"
	m "github.com/brianfloyd/the-grid/view/model"
)

type IWorkoutViewService interface {
	GetWorkout(ctx context.Context, date string, uid string) templ.Component
}

type WorkoutViewHandler struct {
	svc IWorkoutViewService
}

func NewWorkoutViewHandler(svc IWorkoutViewService) *WorkoutViewHandler {
	return &WorkoutViewHandler{
		svc: svc,
	}
}

func (h *WorkoutViewHandler) Register(r *chi.Mux) {
	r.Post("/_t/workout", h.getWorkout)
}

func (h *WorkoutViewHandler) getWorkout(w http.ResponseWriter, r *http.Request) {
	if uid, ok := GetUid(w, r); ok {
		var request m.GetWorkoutByDateRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			logger.ErrorArgs(r.Context(), "Could not decode json body. %s\n", err)
			panic("crash")
		}
		defer r.Body.Close()
		h.svc.GetWorkout(r.Context(), request.Date, uid).Render(r.Context(), w)
	}
}
