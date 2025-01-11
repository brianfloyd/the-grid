package service

import (
	"context"
	"slices"
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
		setViews := buildSetViews(sets, exercisesById)
		groupView := m.WorkoutGroupView{
			Group: m.GroupDescriptorView{
				Value:         group,
				FriendlyValue: getFriendlyGroupValue(group),
			},
			Count: strconv.Itoa(len(setViews)),
			Sets:  setViews,
		}
		mv[group] = groupView
	}

	return mv
}

func buildSetViews(sets []im.Set, exercisesById map[string]im.Exercise) []m.WorkoutSetsView {
	setViewsMap := make(map[m.WorkoutExerciseView][]m.SetView)

	for _, set := range sets {
		exerciseName := "NA"
		if exercise, ok := exercisesById[set.ExerciseId]; ok {
			exerciseName = exercise.Name
		}

		workoutExerciseView := m.WorkoutExerciseView{
			Id:   set.ExerciseId,
			Name: exerciseName,
		}

		setView := buildSetView(set)
		if slice, ok := setViewsMap[workoutExerciseView]; ok {
			slice = append(slice, setView)
			setViewsMap[workoutExerciseView] = slice
		} else {
			setViewsMap[workoutExerciseView] = []m.SetView{setView}
		}

	}

	workoutSets := make([]m.WorkoutSetsView, len(setViewsMap))
	i := 0
	for k, v := range setViewsMap {
		workoutSets[i] = m.WorkoutSetsView{
			WorkoutExerciseView: k,
			Sets:                v,
			Count:               strconv.Itoa(len(v)),
			CountValue:          len(v),
		}
		i += 1
	}

	slices.SortFunc(workoutSets, func(a, b m.WorkoutSetsView) int {
		if a.Count != b.Count {
			return b.CountValue - a.CountValue
		}
		return strings.Compare(a.WorkoutExerciseView.Name, b.WorkoutExerciseView.Name)
	})

	return workoutSets
}

func buildSetView(set im.Set) m.SetView {
	return m.SetView{
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
