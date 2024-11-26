package db

import (
	"errors"
)

var ErrDbNotFound = errors.New("could not find requested id")
var ErrDbGeneric = errors.New("generic database error")
