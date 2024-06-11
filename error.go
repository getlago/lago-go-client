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

type ErrorDetail struct {
	Multiple bool
	Errors   map[int]map[string][]string
}

func (ed *ErrorDetail) DetailsForRow(row int) (map[string][]string, error) {
	if !ed.Multiple {
		return nil, errors.New("error contains a single error, use Error()")
	}

	rowDetail, ok := ed.Errors[row]
	if !ok {
		return nil, nil
	}

	return rowDetail, nil
}

func (ed *ErrorDetail) Details() (map[string][]string, error) {
	if ed.Multiple {
		return nil, errors.New("error contains multiple errors, use ErrorForRow()")
	}

	return ed.Errors[0], nil
}

func (ed *ErrorDetail) UnmarshalJSON(data []byte) error {
	// First attempt to unmarshal singular.
	var singularErr map[string][]string
	err := json.Unmarshal(data, &singularErr)
	if err == nil {
		ed.Errors = map[int]map[string][]string{
			0: singularErr,
		}
		return nil
	}

	// Then attempt to unmarshal multiple.
	var multipleErr map[int]map[string][]string
	err = json.Unmarshal(data, &multipleErr)
	if err == nil {
		ed.Errors = multipleErr
		ed.Multiple = true
		return nil
	}

	return err
}

func (ed *ErrorDetail) MarshalJSON() ([]byte, error) {
	if ed.Multiple {
		return json.Marshal(ed.Errors)
	}
	return json.Marshal(ed.Errors[0])
}

type Error struct {
	Err error `json:"-"`

	HTTPStatusCode int    `json:"status"`
	Message        string `json:"error"`
	ErrorCode      string `json:"code"`

	ErrorDetail *ErrorDetail `json:"error_details,omitempty"`
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
