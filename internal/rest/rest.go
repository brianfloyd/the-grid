package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

const uuidRegEx string = `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`

var (
	GenericError = ErrorCode{Name: "GENERIC_ERROR", Status: 500}
	NotFound     = ErrorCode{Name: "NOT_FOUND", Status: 404}
	BadRequest   = ErrorCode{Name: "BAD_REQUEST", Status: 400}
)

type ErrorCode struct {
	Name   string
	Status int
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func renderErrorResponse(w http.ResponseWriter, r *http.Request, errorCode ErrorCode, message string, e error) {
	fmt.Printf("ERROR [%s] %s - %v\n", errorCode.Name, message, e)

	response := ErrorResponse{
		Code:    errorCode.Name,
		Status:  errorCode.Status,
		Message: message,
	}

	render.Status(r, errorCode.Status)
	render.JSON(w, r, &response)
}

func renderResponse(w http.ResponseWriter, r *http.Request, response interface{}, status int) {
	render.Status(r, status)
	render.JSON(w, r, response)
}
