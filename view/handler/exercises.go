package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/brianfloyd/the-grid/util"
	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/brianfloyd/the-grid/view/template/component"
	"github.com/go-chi/chi"
)

type IExerciseViewService interface {
	GetExercisesForGroupPage(ctx context.Context, group, date, uid string) templ.Component
	AddExerciseToWorkout(ctx context.Context, exerciseId, date, uid string) templ.Component
	RemoveExerciseFromWorkout(ctx context.Context, exerciseId, date, uid string) templ.Component
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
	r.Post("/exercises/form/add", e.addExerciseToWorkout)
	r.Post("/exercises/form/remove", e.removeExerciseFromWorkout)
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

func (e *ExerciseViewHandler) removeExerciseFromWorkout(w http.ResponseWriter, r *http.Request) {
	if uid, ok := util.GetUid(w, r); ok {
		ctx := r.Context()

		fd, err := readExerciseFormData(r)
		if err != nil {
			component.GlobalErrorComponent("Invalid data submission.").Render(ctx, w)
			return
		}

		e.svc.RemoveExerciseFromWorkout(ctx, fd.ExerciseId, fd.Date, uid).Render(ctx, w)
	}
}

func (e *ExerciseViewHandler) addExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	if uid, ok := util.GetUid(w, r); ok {
		ctx := r.Context()

		fd, err := readExerciseFormData(r)
		if err != nil {
			component.GlobalErrorComponent("Invalid data submission.").Render(ctx, w)
			return
		}

		e.svc.AddExerciseToWorkout(ctx, fd.ExerciseId, fd.Date, uid).Render(ctx, w)
	}
}

func readExerciseFormData(r *http.Request) (m.ExerciseFormData, error) {
	err := r.ParseForm()
	if err != nil {
		return m.ExerciseFormData{}, err
	}

	return m.ExerciseFormData{
		Date:       r.Form.Get("date"),
		Group:      r.Form.Get("group"),
		ExerciseId: r.Form.Get("exerciseId"),
		TargetId:   r.Form.Get("targetId"),
	}, nil
}
