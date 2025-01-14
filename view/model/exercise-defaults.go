package model

import "github.com/a-h/templ"

type ExerciseDefaultView struct {
	ExerciseId     string
	Name           string
	Weight         string
	Reps           string
	EditButtonText string
	EditButtonUrl  templ.SafeURL
}

type ExerciseDefaultEditView struct {
	ExerciseId     string
	Name           string
	Weight         string
	Reps           string
	SaveButtonText string
	SaveUrl        templ.SafeURL
}
