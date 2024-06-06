package lago

import (
	"encoding/json"
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

func (e Error) Error() string {
	type alias struct {
		Error
		Err string `json:"err,omitempty"`
	}
	err := alias{Error: e}
	if e.Err != nil {
		err.Err = e.Err.Error()
	}
	msg, _ := json.Marshal(&err)
	return string(msg)
}

func (e ErrorCode) Error() string {
	return string(e)
}
