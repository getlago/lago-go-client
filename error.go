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
	Msg:            "type assertion failed",
}

type ErrorDetail struct {
	ErrorCode []ErrorCode `json:"code,omitempty"`
}

type Error struct {
	Err error `json:"-"`

	HTTPStatusCode int    `json:"status"`
	Msg            string `json:"message"`

	ErrorDetail ErrorDetail `json:"error_details"`
}

func (e ErrorCode) Error() string {
	return string(e)
}
