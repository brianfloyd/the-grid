package service

import (
	"context"

	m "github.com/brianfloyd/the-grid/view/model"

	"github.com/a-h/templ"
	is "github.com/brianfloyd/the-grid/internal/service"
	"github.com/brianfloyd/the-grid/view/template/page"
)

type IExerciseDefaultViewService interface {
	GetExerciseDefaultsPage(ctx context.Context) templ.Component
}

type ExerciseDefaultViewService struct {
	exerciseService     is.IExercisesService
	exerciseViewService IExerciseViewService
}

func NewExerciseDefaultViewService(exerciseService is.IExercisesService, exerciseViewService IExerciseViewService) *ExerciseDefaultViewService {
	return &ExerciseDefaultViewService{
		exerciseService:     exerciseService,
		exerciseViewService: exerciseViewService,
	}
}

func (e *ExerciseDefaultViewService) GetExerciseDefaultsPage(ctx context.Context) templ.Component {
	exerciseDefaultViews := []m.ExerciseDefaultView{
		{
			ExerciseId:     "foo",
			Name:           "Bicep Curls",
			Weight:         "10",
			Reps:           "10",
			EditButtonText: "Edit",
		},
		{
			ExerciseId:     "foo",
			Name:           "Preacher Curls",
			Weight:         "10",
			Reps:           "10",
			EditButtonText: "Edit",
		},
		{
			ExerciseId:     "foo",
			Name:           "Preacher Curls",
			Weight:         "10",
			Reps:           "10",
			EditButtonText: "Edit",
		},
		{
			ExerciseId:     "foo",
			Name:           "Preacher Curls",
			Weight:         "10",
			Reps:           "10",
			EditButtonText: "Edit",
		},
		{
			ExerciseId:     "foo",
			Name:           "Preacher Curls",
			Weight:         "10",
			Reps:           "10",
			EditButtonText: "Edit",
		},
		{
			ExerciseId:     "foo",
			Name:           "Preacher Curls",
			Weight:         "10",
			Reps:           "10",
			EditButtonText: "Edit",
		},
	}

	exerciseGroups := e.exerciseViewService.GetGroupViews("01-11-2025")
	return page.ExerciseDefaultsPage(exerciseDefaultViews, exerciseGroups)
}
