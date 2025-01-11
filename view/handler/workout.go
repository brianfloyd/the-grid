package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi"

	"github.com/brianfloyd/the-grid/util"
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
	r.Get("/", h.getTodayWorkout)
	r.Get(fmt.Sprintf("/{date:%s}", m.DATE_REGEX), h.getWorkout)
}

func (h *WorkoutViewHandler) getTodayWorkout(w http.ResponseWriter, r *http.Request) {
	h.getWorkoutWithDate(w, r, util.MakeDateStringFromTime(time.Now()))

}

func (h *WorkoutViewHandler) getWorkout(w http.ResponseWriter, r *http.Request) {
	date := chi.URLParam(r, "date")
	h.getWorkoutWithDate(w, r, date)
}

func (h *WorkoutViewHandler) getWorkoutWithDate(w http.ResponseWriter, r *http.Request, date string) {
	if uid, ok := util.GetUid(w, r); ok {
		ctx := r.Context()
		h.svc.GetWorkout(ctx, date, uid).Render(ctx, w)
	}
}
