package handler

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/brianfloyd/the-grid/util"
	"github.com/go-chi/chi"
)

type IExerciseDefaultViewService interface {
	GetExerciseDefaultsPage(ctx context.Context) templ.Component
}

type ExerciseDefaultViewHandler struct {
	svc IExerciseDefaultViewService
}

func NewExerciseDefaultViewHandler(svc IExerciseDefaultViewService) *ExerciseDefaultViewHandler {
	return &ExerciseDefaultViewHandler{
		svc: svc,
	}
}

func (e *ExerciseDefaultViewHandler) Register(r *chi.Mux) {
	r.Get("/exercise-defaults", e.getExerciseDefaults)
}

func (e *ExerciseDefaultViewHandler) getExerciseDefaults(w http.ResponseWriter, r *http.Request) {
	if _, ok := util.GetUid(w, r); ok {
		ctx := r.Context()
		e.svc.GetExerciseDefaultsPage(ctx).Render(ctx, w)
	}
}
