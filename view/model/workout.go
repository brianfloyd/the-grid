package model

type WorkoutView struct {
	Date                        string
	Groups                      map[string]WorkoutGroupView
	PreferredExerciseGroupOrder []string
}

type WorkoutGroupView struct {
	Group GroupDescriptorView
	Sets  []SetView
}

type SetView struct {
	Exercise   WorkoutExerciseView
	Weight     string
	WeightType string
	Reps       string
}

type GroupDescriptorView struct {
	Value         string
	FriendlyValue string
}

type WorkoutExerciseView struct {
	Id   string
	Name string
}

type GetWorkoutByDateRequest struct {
	Date string `json:"date"`
}
