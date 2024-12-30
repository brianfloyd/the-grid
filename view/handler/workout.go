package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi"
)

type WorkoutViewService interface {
	GetExercises() templ.Component
}

type WorkoutViewHandler struct {
	svc WorkoutViewService
}

func NewWorkoutViewHandler(svc WorkoutViewService) *WorkoutViewHandler {
	return &WorkoutViewHandler{
		svc: svc,
	}
}

func (w *WorkoutViewHandler) Register(r *chi.Mux) {
	r.Get("/_t/exercises", w.getExercises)
}

// func (u *UserHandler) create(w http.ResponseWriter, r *http.Request) {

func (wv *WorkoutViewHandler) getExercises(w http.ResponseWriter, r *http.Request) {
	wv.svc.GetExercises().Render(r.Context(), w)
}
