package service

import (
	"context"
	"strconv"
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

type IWorkoutViewService interface {
	GetWorkout(ctx context.Context, date string, uid string) templ.Component
}

type WorkoutViewService struct {
	workoutService      is.IWorkoutService
	exercisesService    is.IExercisesService
	exerciseViewService IExerciseViewService
}

func NewWorkoutViewService(workoutService is.IWorkoutService, exercisesService is.IExercisesService, exerciseViewService IExerciseViewService) *WorkoutViewService {
	return &WorkoutViewService{
		workoutService:      workoutService,
		exercisesService:    exercisesService,
		exerciseViewService: exerciseViewService,
	}
}

func (s *WorkoutViewService) GetWorkout(ctx context.Context, date string, uid string) templ.Component {
	if _, err := util.SanitizeDate(date); err != nil {
		date = util.MakeDateStringFromTime(time.Now())
	}

	groups := s.exerciseViewService.GetGroupViews(date)

	workoutView := m.WorkoutView{
		Date:                        date,
		PreferredExerciseGroupOrder: getPreferredExerciseGroupOrder(),
	}

	workout, err := s.workoutService.ByDate(ctx, uid, date)

	if err == nil {
		workoutView.Date = workout.Date
		workoutView.Groups = s.buildWorkoutGroups(ctx, workout)
	} else {
		logger.ErrorArgs(ctx, "Could not load workout for date '%s'. %s\n", date, err)
	}

	return page.WorkoutPage(workoutView, groups)
}

func (s *WorkoutViewService) buildWorkoutGroups(ctx context.Context, workout im.Workout) map[string]m.WorkoutGroupView {
	mv := make(map[string]m.WorkoutGroupView)

	exercisesById := s.getExercisesById(ctx)
	setsByGroup := s.getSetsByGroup(workout.Sets, exercisesById)

	for group, sets := range setsByGroup {
		groupView := m.WorkoutGroupView{
			Group: m.GroupDescriptorView{
				Value:         group,
				FriendlyValue: getFriendlyGroupValue(group),
			},
			Sets: buildSetViews(sets, exercisesById),
		}
		mv[group] = groupView
	}

	return mv
}

func buildSetViews(sets []im.Set, exercisesById map[string]im.Exercise) []m.SetView {
	setViews := make([]m.SetView, len(sets))
	for i, set := range sets {
		setViews[i] = buildSetView(set, exercisesById)
	}
	return setViews
}

func buildSetView(set im.Set, exercisesById map[string]im.Exercise) m.SetView {
	exerciseName := "NA"
	if exercise, ok := exercisesById[set.ExerciseId]; ok {
		exerciseName = exercise.Name
	}

	return m.SetView{
		Exercise: m.WorkoutExerciseView{
			Id:   set.ExerciseId,
			Name: exerciseName,
		},
		Weight:     strconv.FormatUint(set.Weight, 10),
		WeightType: "lbs",
		Reps:       strconv.FormatUint(set.Reps, 10),
	}
}

func getFriendlyGroupValue(group string) string {
	return string(group[0]) + strings.ToLower(group[1:])
}

func (s *WorkoutViewService) getSetsByGroup(sets []im.Set, exercisesById map[string]im.Exercise) map[string][]im.Set {
	m := make(map[string][]im.Set)
	for _, set := range sets {
		if exercise, ok := exercisesById[set.ExerciseId]; ok {
			grp := string(exercise.Group)
			m[grp] = append(m[grp], set)
		}
	}
	return m
}

func (s *WorkoutViewService) getExercisesById(ctx context.Context) map[string]im.Exercise {
	m := make(map[string]im.Exercise)

	exercises, err := s.exercisesService.List(ctx)
	if err != nil {
		return m
	}

	for _, exercise := range exercises {
		m[exercise.Id] = exercise
	}

	return m
}

func getPreferredExerciseGroupOrder() []string {
	return []string{
		string(im.ExerciseGroupBiceps),
		string(im.ExerciseGroupBack),
		string(im.ExerciseGroupTricep),
		string(im.ExerciseGroupChest),
		string(im.ExerciseGroupShoulder),
		string(im.ExerciseGroupLegs),
		string(im.ExerciseGroupAbs),
		string(im.ExerciseGroupCardio),
	}
}
