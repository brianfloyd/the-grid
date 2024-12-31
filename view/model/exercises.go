package model

type ExerciseGroupView struct {
	Name     string
	ImageUrl string
	Selected bool
}

type ExerciseView struct {
	Id    string
	Group string
	Name  string
}
