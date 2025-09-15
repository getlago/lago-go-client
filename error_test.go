package lago

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestError_Err(t *testing.T) {
	var hasErr error = Error{
		Err:            errors.New("type assertion failed"),
		HTTPStatusCode: 422,
		Message:        "Type assertion failed",
	}
	t.Logf("%s", hasErr.Error())
}

func TestError_NoErr(t *testing.T) {
	var noErr error = Error{
		HTTPStatusCode: 500,
		Message:        "500",
	}
	t.Logf("%s", noErr.Error())
}

func TestError_Details(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  *Error
	}{
		{
			name: "Single detail",
			input: `{
  "status": 422,
  "error": "Unprocessable Entity",
  "code": "validation_errors",
  "error_details": {
    "transaction_id": [
      "value_already_exist"
    ]
  }
}`,
			want: &Error{
				HTTPStatusCode: 422,
				Message:        "Unprocessable Entity",
				ErrorCode:      "validation_errors",
				ErrorDetail: &ErrorDetail{
					Multiple: false,
					Errors: map[int]map[string][]string{
						0: {
							"transaction_id": {
								"value_already_exist",
							},
						},
					},
				},
			},
		},
		{
			name: "Multiple details",
			input: `{
  "status": 422,
  "error": "Unprocessable Entity",
  "code": "validation_errors",
  "error_details": {
    "0": {
      "transaction_id": [
        "value_already_exist"
      ]
    },
    "1": {
      "transaction_id": [
        "value_already_exist"
      ]
    },
    "2": {
      "transaction_id": [
        "value_already_exist"
      ]
    }
  }
}`,
			want: &Error{
				HTTPStatusCode: 422,
				Message:        "Unprocessable Entity",
				ErrorCode:      "validation_errors",
				ErrorDetail: &ErrorDetail{
					Multiple: true,
					Errors: map[int]map[string][]string{
						0: {
							"transaction_id": {
								"value_already_exist",
							},
						},
						1: {
							"transaction_id": {
								"value_already_exist",
							},
						},
						2: {
							"transaction_id": {
								"value_already_exist",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errObj := &Error{}
			err := json.Unmarshal([]byte(tt.input), errObj)
			if err != nil {
				t.Errorf("got error %s", err.Error())
				return
			}

			expectErr, err := json.Marshal(tt.want)
			if err != nil {
				t.Errorf("got error %s", err.Error())
				return
			}

			gotErr, err := json.Marshal(errObj)
			if err != nil {
				t.Errorf("got error %s", err.Error())
				return
			}

			if string(expectErr) != string(gotErr) {
				t.Errorf("got error %s, but expected error %s", string(gotErr), string(expectErr))
				return
			}
		})
	}
}
