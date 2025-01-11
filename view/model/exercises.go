package model

import "github.com/a-h/templ"

type ExerciseGroupView struct {
	Name           string
	ImageUrl       string
	Selected       bool
	NavigationLink templ.SafeURL
}

type ExerciseView struct {
	Id        string
	Group     string
	Name      string
	InWorkout bool
}
