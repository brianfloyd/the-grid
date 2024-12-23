package rest

const UUID_REGEX string = `[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`

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
