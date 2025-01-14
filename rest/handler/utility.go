package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	rm "github.com/brianfloyd/the-grid/rest/model"
	"github.com/go-chi/render"
)

func renderErrorResponse(w http.ResponseWriter, r *http.Request, errorCode rm.ErrorCode, message string, e error) {
	fmt.Printf("ERROR [%s] %s - %v\n", errorCode.Name, message, e)

	response := rm.ErrorResponse{
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

func decodeBody(w http.ResponseWriter, r *http.Request, body interface{}, e string) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		renderErrorResponse(w, r, rm.BadRequest, fmt.Sprintf("Could not convert given body to a %s request.", e), err)
		return false
	}
	return true
}
