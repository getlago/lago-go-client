package lago

import (
	"errors"
	"testing"
)

func TestErrorErr(t *testing.T) {
	var hasErr error = Error{
		Err:            errors.New("type assertion failed"),
		HTTPStatusCode: 422,
		Message:        "Type assertion failed",
	}
	t.Logf("%s", hasErr.Error())
}

func TestErrorNoErr(t *testing.T) {
	var noErr error = Error{
		HTTPStatusCode: 500,
		Message:        "500",
	}
	t.Logf("%s", noErr.Error())
}
