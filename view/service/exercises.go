package service

import (
	"context"

	"github.com/a-h/templ"
	im "github.com/brianfloyd/the-grid/internal/model"
	is "github.com/brianfloyd/the-grid/internal/service"
	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/brianfloyd/the-grid/view/template"
)

type IExerciseViewService interface {
	GetExerciseGroups() templ.Component
	GetExercisesForGroup(ctx context.Context, group string) templ.Component
}

type ExerciseViewService struct {
	svc is.IExercisesService
}

func NewExerciseViewService(svc is.IExercisesService) *ExerciseViewService {
	return &ExerciseViewService{
		svc: svc,
	}
}

func (e *ExerciseViewService) GetExerciseGroups() templ.Component {
	return template.GetExerciseGroups(getGroups())
}

func (e *ExerciseViewService) GetExercisesForGroup(ctx context.Context, group string) templ.Component {
	groups := getGroups()
	for i := 0; i < len(groups); i++ {
		if groups[i].Name == group {
			groups[i].Selected = true
			break
		}
	}

	exercises, err := e.svc.ListForGroup(ctx, group)
	if err != nil {
		panic("HANDLE ME")
	}

	viewExercises := make([]m.ExerciseView, len(exercises))
	for i, e := range exercises {
		viewExercises[i] = m.ExerciseView{
			Id:    e.Id,
			Name:  e.Name,
			Group: string(e.Group),
		}
	}

	return template.SelectedExerciseGroup(groups, viewExercises)
}

func getGroups() []m.ExerciseGroupView {
	return []m.ExerciseGroupView{
		{
			Name:     string(im.ExerciseGroupBiceps),
			ImageUrl: "/static/images/icons/bicep.png",
			Selected: false,
		},
		{
			Name:     string(im.ExerciseGroupBack),
			ImageUrl: "/static/images/icons/back.png",
			Selected: false,
		},
		{
			Name:     string(im.ExerciseGroupTricep),
			ImageUrl: "/static/images/icons/tricep.png",
			Selected: false,
		},
		{
			Name:     string(im.ExerciseGroupChest),
			ImageUrl: "/static/images/icons/chest.png",
			Selected: false,
		},
		{
			Name:     string(im.ExerciseGroupShoulder),
			ImageUrl: "/static/images/icons/shoulder.png",
			Selected: false,
		},
		{
			Name:     string(im.ExerciseGroupLegs),
			ImageUrl: "/static/images/icons/legs.png",
			Selected: false,
		},
		{
			Name:     string(im.ExerciseGroupAbs),
			ImageUrl: "/static/images/icons/abs.png",
			Selected: false,
		},
		{
			Name:     string(im.ExerciseGroupCardio),
			ImageUrl: "/static/images/icons/misc.png",
			Selected: false,
		},
	}
}
