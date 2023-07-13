package lago

import (
	"errors"
	"net/http"
)

type ErrorCode string

const (
	ErrorCodeAlreadyExist ErrorCode = "value_already_exist"
	ErrorCodeInvalidValue
)

var ErrorTypeAssert = Error{
	Err:            errors.New("type assertion failed"),
	HTTPStatusCode: http.StatusUnprocessableEntity,
	Message:        "Type assertion failed",
}

type Error struct {
	Err error `json:"-"`

	HTTPStatusCode int    `json:"status"`
	Message        string `json:"error"`
	ErrorCode      string `json:"code"`

	ErrorDetail map[string][]string `json:"error_details,omitempty"`
}

func (e ErrorCode) Error() string {
	return string(e)
}
