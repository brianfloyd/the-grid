package service

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/brianfloyd/the-grid/internal/logger"
	im "github.com/brianfloyd/the-grid/internal/model"
	is "github.com/brianfloyd/the-grid/internal/service"
	"github.com/brianfloyd/the-grid/util"
	m "github.com/brianfloyd/the-grid/view/model"
	"github.com/brianfloyd/the-grid/view/template/page"
)

type IExerciseViewService interface {
	GetExercisesForGroupPage(ctx context.Context, group, date, uid string) templ.Component
	GetGroupViews(date string) []m.ExerciseGroupView
}

type ExerciseViewService struct {
	exerciseService is.IExercisesService
	workoutService  is.IWorkoutService
}

func NewExerciseViewService(exerciseService is.IExercisesService, workoutService is.IWorkoutService) *ExerciseViewService {
	return &ExerciseViewService{
		exerciseService: exerciseService,
		workoutService:  workoutService,
	}
}

func (e *ExerciseViewService) GetExercisesForGroupPage(ctx context.Context, group, date, uid string) templ.Component {
	if _, err := util.SanitizeDate(date); err != nil {
		date = util.MakeDateStringFromTime(time.Now())
	}

	groups := UpdateSelectedGroupProperties(e.GetGroupViews(date), group, date)
	viewExercises := e.GetExerciseViews(ctx, group)

	workout, err := e.workoutService.ByDate(ctx, uid, date)
	if err != nil {
		logger.ErrorArgs(ctx, "Could not get workout by date (%v)! %v\n", date, err)
	} else {
		viewExercises = UpdateExercisesForWorkoutState(viewExercises, workout)
	}

	return page.ExercisePage(groups, viewExercises)
}

func (e *ExerciseViewService) GetExerciseViews(ctx context.Context, group string) []m.ExerciseView {
	exercises, err := e.exerciseService.ListForGroup(ctx, group)
	if err != nil {
		panic("HANDLE ME")
	}

	viewExercises := make([]m.ExerciseView, len(exercises))
	for i, e := range exercises {
		viewExercises[i] = m.ExerciseView{
			Id:        e.Id,
			Name:      e.Name,
			Group:     string(e.Group),
			InWorkout: false,
		}
	}

	return viewExercises
}

func (e *ExerciseViewService) GetGroupViews(date string) []m.ExerciseGroupView {
	groupViews := make([]m.ExerciseGroupView, len(im.ExerciseGroupAll))
	for i, e := range im.ExerciseGroupAll {
		groupViews[i] = buildGroupView(e, date)
	}
	return groupViews
}

func UpdateExercisesForWorkoutState(viewExercises []m.ExerciseView, workout im.Workout) []m.ExerciseView {
	exerciseIdsInWorkout := []string{}
	for _, set := range workout.Sets {
		exerciseIdsInWorkout = append(exerciseIdsInWorkout, set.ExerciseId)
	}

	for i := 0; i < len(viewExercises); i++ {
		viewExercises[i].InWorkout = slices.Contains(exerciseIdsInWorkout, viewExercises[i].Id)
	}
	return viewExercises
}

func UpdateSelectedGroupProperties(groups []m.ExerciseGroupView, group string, date string) []m.ExerciseGroupView {
	for i := 0; i < len(groups); i++ {
		if strings.EqualFold(groups[i].Name, group) {
			groups[i].Selected = true
			groups[i].NavigationLink = templ.SafeURL(fmt.Sprintf("/%s", date))
			break
		}
	}
	return groups
}

func buildGroupView(group im.ExerciseGroup, date string) m.ExerciseGroupView {
	lower := strings.ToLower(string(group))
	return m.ExerciseGroupView{
		Name:           string(group),
		ImageUrl:       fmt.Sprintf("/static/images/icons/%s.png", lower),
		Selected:       false,
		NavigationLink: templ.SafeURL(fmt.Sprintf("/exercises/%s?date=%s", lower, date)),
	}
}
