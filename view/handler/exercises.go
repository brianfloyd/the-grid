package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/brianfloyd/the-grid/util"
	"github.com/go-chi/chi"
)

type IExerciseViewService interface {
	GetExercisesForGroupPage(ctx context.Context, group, date, uid string) templ.Component
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
	r.Get(fmt.Sprintf("/exercises/{group:%s}", "[A-Za-z]+"), e.getExercisesForGroupPage)
}

func (e *ExerciseViewHandler) getExercisesForGroupPage(w http.ResponseWriter, r *http.Request) {
	if uid, ok := util.GetUid(w, r); ok {
		ctx := r.Context()
		q := r.URL.Query()

		date := q.Get("date")
		group := chi.URLParam(r, "group")

		e.svc.GetExercisesForGroupPage(ctx, group, date, uid).Render(ctx, w)
	}
}
