package util

import (
	"errors"
	"net/http"
	"time"

	vm "github.com/brianfloyd/the-grid/view/model"
)

func MakeDateStringFromTime(date time.Time) string {
	return date.Format("01-02-2006")
}

func SanitizeDate(date string) (time.Time, error) {
	t, err := time.Parse("01-02-2006", date)
	if err != nil {
		t, err = time.Parse("01/02/2006", date)
		if err != nil {
			return time.Time{}, errors.New("could not convert given string to a recognized time format")
		}
	}
	return t, nil
}

func GetUid(w http.ResponseWriter, r *http.Request) (string, bool) {
	userId := r.CookiesNamed(vm.COOKIE_UID)

	if len(userId) == 0 {
		Redirect(w, "/login")
		return "", false
	}

	return userId[0].Value, true
}

func Redirect(w http.ResponseWriter, location string) {
	w.Header().Add("Location", location)
	w.WriteHeader(301)
}

func HTMXRedirect(w http.ResponseWriter, location string) {
	w.Header().Add("HX-Redirect", location)
	w.WriteHeader(200)
}
