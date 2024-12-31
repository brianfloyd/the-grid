package service

import (
	"github.com/a-h/templ"
	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/brianfloyd/the-grid/view/template"
)

type IExerciseViewService interface {
	GetExerciseGroups() templ.Component
	GetExercisesForGroup(string) templ.Component
}

type ExerciseViewService struct {
}

func NewExerciseViewService() *ExerciseViewService {
	return &ExerciseViewService{}
}

func (e *ExerciseViewService) GetExerciseGroups() templ.Component {
	return template.GetExerciseGroups(getGroups())
}

func (e *ExerciseViewService) GetExercisesForGroup(group string) templ.Component {
	exercises := []m.ExerciseView{}
	if group == "BICEP" {
		exercises = []m.ExerciseView{
			{
				Id:    "1",
				Group: "BICEP",
				Name:  "Dumbbell Curls",
			},
			{
				Id:    "2",
				Group: "BICEP",
				Name:  "Hammer Curls",
			},
			{
				Id:    "3",
				Group: "BICEP",
				Name:  "Preacher Curls",
			},
		}
	}

	groups := getGroups()
	for i := 0; i < len(groups); i++ {
		if groups[i].Name == group {
			groups[i].Selected = true
			break
		}
	}

	return template.SelectedExerciseGroup(groups, exercises)
}

func getGroups() []m.ExerciseGroupView {
	return []m.ExerciseGroupView{
		{
			Name:     "BICEP",
			ImageUrl: "/static/images/icons/bicep.png",
			Selected: false,
		},
		{
			Name:     "BACK",
			ImageUrl: "/static/images/icons/back.png",
			Selected: false,
		},
		{
			Name:     "TRICEP",
			ImageUrl: "/static/images/icons/tricep.png",
			Selected: false,
		},
		{
			Name:     "CHEST",
			ImageUrl: "/static/images/icons/chest.png",
			Selected: false,
		},
		{
			Name:     "SHOULDER",
			ImageUrl: "/static/images/icons/shoulder.png",
			Selected: false,
		},
		{
			Name:     "LEGS",
			ImageUrl: "/static/images/icons/legs.png",
			Selected: false,
		},
		{
			Name:     "ABS",
			ImageUrl: "/static/images/icons/abs.png",
			Selected: false,
		},
		{
			Name:     "CARDIO",
			ImageUrl: "/static/images/icons/misc.png",
			Selected: false,
		},
	}
}
