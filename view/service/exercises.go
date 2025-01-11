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
	"github.com/brianfloyd/the-grid/view/template/component"
	"github.com/brianfloyd/the-grid/view/template/page"
)

type IExerciseViewService interface {
	GetExercisesForGroupPage(ctx context.Context, group, date, uid string) templ.Component
	AddExerciseToWorkout(ctx context.Context, exerciseId, date, uid string) templ.Component
	RemoveExerciseFromWorkout(ctx context.Context, exerciseId, date, uid string) templ.Component
	GetGroupViews(date string) []m.ExerciseGroupViewResponse
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

	return page.ExercisePage(groups, buildExerciseViewResponse(viewExercises, date))
}

func (e *ExerciseViewService) RemoveExerciseFromWorkout(ctx context.Context, exerciseId, date, uid string) templ.Component {
	if _, err := util.SanitizeDate(date); err != nil {
		date = util.MakeDateStringFromTime(time.Now())
	}

	workout, err := e.workoutService.ByDate(ctx, uid, date)
	if err != nil {
		logger.ErrorArgs(ctx, "Could not get workout by date (%s). %v\n", date, err)
		return component.GlobalErrorComponent("Could not find a workout for the current date.")
	}

	setIds := []string{}
	for _, set := range workout.Sets {
		if set.ExerciseId == exerciseId {
			setIds = append(setIds, set.Id)
		}
	}

	err = e.workoutService.DeleteSets(ctx, workout.Id, setIds)
	if err != nil {
		logger.ErrorArgs(ctx, "Could not remove exercises from workout.", exerciseId, workout.Id, err)
		return component.GlobalErrorComponent("An error occurred removing the exercise from the workout.")
	}

	return nil
}

func (e *ExerciseViewService) AddExerciseToWorkout(ctx context.Context, exerciseId, date, uid string) templ.Component {
	if _, err := util.SanitizeDate(date); err != nil {
		date = util.MakeDateStringFromTime(time.Now())
	}

	exercises, err := e.exerciseService.List(ctx)
	if err != nil {
		logger.ErrorArgs(ctx, "Could not get all exercises. %v\n", err)
		return component.GlobalErrorComponent("An error occurred getting all exercises.")
	}

	var exercise *im.Exercise
	for _, e := range exercises {
		if e.Id == exerciseId {
			exercise = &e
		}
	}

	if exercise == nil {
		logger.ErrorArgs(ctx, "Could not find exercise for id (%s).\n", exerciseId)
		return component.GlobalErrorComponent("Could not find an exercise matching the specified id.")
	}

	workout, err := e.workoutService.ByDate(ctx, uid, date)
	if err != nil {
		logger.ErrorArgs(ctx, "Could not get workout by date (%s). %v\n", date, err)
		return component.GlobalErrorComponent("Could not find a workout for the current date.")
	}

	for _, set := range workout.Sets {
		if set.ExerciseId == exerciseId {
			logger.WarnArgs(ctx, "Workout (%s) already had the exercise (%s).", workout.Id, exerciseId)
			return nil
		}
	}

	// TODO: Get default reps, weight, count, etc.
	_, err = e.workoutService.CreateSet(ctx, workout.Id, im.Set{
		ExerciseId: exerciseId,
	})

	if err != nil {
		logger.ErrorArgs(ctx, "Could not add exercise (%s) to workout (%s). %v\n", exerciseId, workout.Id, err)
		return component.GlobalErrorComponent("An error occurred adding the exercise to the workout.")
	}

	return nil
}

func buildExerciseViewResponse(views []m.ExerciseView, date string) m.ExerciseViewResponse {
	return m.ExerciseViewResponse{
		Views: views,
		Meta: m.ExerciseViewMeta{
			Date: date,
		},
	}
}

func (e *ExerciseViewService) GetExerciseViews(ctx context.Context, group string) []m.ExerciseView {
	exercises, err := e.exerciseService.ListForGroup(ctx, group)
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

	return viewExercises
}

func (e *ExerciseViewService) GetGroupViews(date string) []m.ExerciseGroupViewResponse {
	groupViews := make([]m.ExerciseGroupViewResponse, len(im.ExerciseGroupAll))
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
		if !slices.Contains(exerciseIdsInWorkout, viewExercises[i].Id) {
			viewExercises[i].Form = m.ExerciseViewForm{
				SubmitButtonText:   "Add",
				SubmitURL:          templ.SafeURL("/exercises/form/add"),
				DisableInputTarget: "[add-remove-btn]",
				IndicatorId:        fmt.Sprintf("submit-indicator-%d", i),
			}
		} else {
			viewExercises[i].Form = m.ExerciseViewForm{
				SubmitButtonText:   "Remove",
				SubmitURL:          templ.SafeURL("/exercises/form/remove"),
				DisableInputTarget: "[add-remove-btn]",
				IndicatorId:        fmt.Sprintf("submit-indicator-%d", i),
			}
		}
	}
	return viewExercises
}

func UpdateSelectedGroupProperties(groups []m.ExerciseGroupViewResponse, group string, date string) []m.ExerciseGroupViewResponse {
	for i := 0; i < len(groups); i++ {
		if strings.EqualFold(groups[i].View.Name, group) {
			groups[i].View.Selected = true
			groups[i].Meta.NavigationLink = templ.SafeURL(fmt.Sprintf("/%s", date))
			break
		}
	}
	return groups
}

func buildGroupView(group im.ExerciseGroup, date string) m.ExerciseGroupViewResponse {
	lower := strings.ToLower(string(group))
	return m.ExerciseGroupViewResponse{
		View: m.ExerciseGroupView{
			Name:     string(group),
			ImageUrl: fmt.Sprintf("/static/images/icons/%s.png", lower),
			Selected: false,
		},
		Meta: m.ExerciseGroupViewMeta{
			NavigationLink: templ.SafeURL(fmt.Sprintf("/exercises/%s?date=%s", lower, date)),
		},
	}
}
