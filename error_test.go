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

func TestError_V2Details(t *testing.T) {
	t.Run("Decodes the v2 code and details", func(t *testing.T) {
		errObj := &Error{}
		input := `{"code":"validation_errors","details":[{"code":"value_is_invalid","field":"signature_algo","meta":{"allowed_values":["jwt","hmac","standard_webhooks"]}}]}`
		if err := json.Unmarshal([]byte(input), errObj); err != nil {
			t.Fatalf("got error %s", err.Error())
		}

		if errObj.ErrorCode != "validation_errors" {
			t.Errorf("got code %q, but expected validation_errors", errObj.ErrorCode)
		}
		if errObj.ErrorDetail != nil {
			t.Errorf("got error_details %v, but expected nil", errObj.ErrorDetail)
		}
		if len(errObj.Details) != 1 || errObj.Details[0].Code != "value_is_invalid" || errObj.Details[0].Field != "signature_algo" {
			t.Fatalf("got details %+v", errObj.Details)
		}
		allowed, ok := errObj.Details[0].Meta["allowed_values"].([]any)
		if !ok || len(allowed) != 3 || allowed[2] != "standard_webhooks" {
			t.Errorf("got allowed_values %v", errObj.Details[0].Meta["allowed_values"])
		}
	})

	t.Run("Leaves Details nil and the Error string unchanged for a v1 body", func(t *testing.T) {
		errObj := &Error{}
		input := `{"status":422,"error":"Unprocessable Entity","code":"validation_errors","error_details":{"transaction_id":["value_already_exist"]}}`
		if err := json.Unmarshal([]byte(input), errObj); err != nil {
			t.Fatalf("got error %s", err.Error())
		}

		if errObj.Details != nil {
			t.Errorf("got details %+v, but expected nil", errObj.Details)
		}
		want := `{"status":422,"error":"Unprocessable Entity","code":"validation_errors","error_details":{"transaction_id":["value_already_exist"]}}`
		if got := errObj.Error(); got != want {
			t.Errorf("got %s, but expected %s", got, want)
		}
	})
}
