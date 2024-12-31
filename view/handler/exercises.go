package handler

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi"
)

type IExerciseViewService interface {
	GetExerciseGroups() templ.Component
	GetExercisesForGroup(string) templ.Component
}

type ExerciseViewHandler struct {
	svc IExerciseViewService
}

func NewExerciseViewHandler(svc IExerciseViewService) *ExerciseViewHandler {
	return &ExerciseViewHandler{
		svc: svc,
	}
}

func (e *ExerciseViewHandler) Register(r *chi.Mux) {
	r.Get("/_t/exercises", e.getExerciseGroups)
	r.Get(fmt.Sprintf("/_t/exercises/{group:%s}", "[A-Za-z]+"), e.getExercisesForGroup)
}

func (e *ExerciseViewHandler) getExerciseGroups(w http.ResponseWriter, r *http.Request) {
	e.svc.GetExerciseGroups().Render(r.Context(), w)
}

func (e *ExerciseViewHandler) getExercisesForGroup(w http.ResponseWriter, r *http.Request) {
	group := chi.URLParam(r, "group")
	e.svc.GetExercisesForGroup(group).Render(r.Context(), w)
}
