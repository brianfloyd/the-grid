package service

import (
	"github.com/a-h/templ"
	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/brianfloyd/the-grid/view/template"
)

type IWorkoutViewService interface {
	GetExercises() templ.Component
}

type WorkoutViewService struct {
}

func NewWorkoutViewService() *WorkoutViewService {
	return &WorkoutViewService{}
}

func (w *WorkoutViewService) GetExercises() templ.Component {
	exercises := []m.ExerciseView{
		{
			Id:       1,
			Name:     "Test",
			ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/91/Octicons-mark-github.svg/640px-Octicons-mark-github.svg.png",
		},
		{
			Id:       2,
			Name:     "Test2",
			ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/thumb/9/91/Octicons-mark-github.svg/640px-Octicons-mark-github.svg.png",
		},
	}
	return template.ExercisesTemplate(exercises)
}
