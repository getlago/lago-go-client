package lago

type ErrorCode string

const (
	ErrorCodeAlreadyExist ErrorCode = "value_already_exist"
)

type ErrorDetail struct {
	ErrorCode []ErrorCode `json:"code"`
}

type Error struct {
	Err error `json:"-"`

	HTTPStatusCode int    `json:"status"`
	Msg            string `json:"message"`

	ErrorDetail ErrorDetail `json:"error_details"`
}
