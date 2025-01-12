package model

import "github.com/a-h/templ"

type ExerciseGroupView struct {
	Name     string
	ImageUrl string
	Selected bool
}

type ExerciseGroupViewMeta struct {
	NavigationLink templ.SafeURL
}

type ExerciseGroupViewResponse struct {
	View ExerciseGroupView
	Meta ExerciseGroupViewMeta
}

type ExerciseView struct {
	Id    string
	Group string
	Name  string
	Form  ExerciseViewForm
}

type ExerciseViewForm struct {
	SubmitURL          templ.SafeURL
	DisableInputTarget string
	SubmitButtonText   string
	IndicatorId        string
	TargetId           string
}

type ExerciseViewMeta struct {
	Date string
}

type ExerciseViewResponse struct {
	Views []ExerciseView
	Meta  ExerciseViewMeta
}

type ExerciseFormData struct {
	Date       string
	Group      string
	ExerciseId string
	TargetId   string
}
