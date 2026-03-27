package errno

import (
	"errors"
)

type Errno interface {
	Error() string
	GetHTTPStatus() int
	GetCode() int
	GetMessage() string
	GetData() any
}

type errno struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code" example:"0"`
	Message    string `json:"message" example:"OK"`
	Data       any    `json:"data"`
}

// Error implements error.
func (e errno) Error() string {
	return e.Message
}

func (e errno) GetHTTPStatus() int {
	return e.HTTPStatus
}

func (e errno) GetCode() int {
	return e.Code
}

func (e errno) GetMessage() string {
	return e.Message
}

func (e errno) GetData() any {
	return e.Data
}

func (e errno) WithData(data any) *errno {
	next := e
	next.Data = data
	return &next
}

func decodeData(data any) Errno {
	return &errno{
		HTTPStatus: OK.HTTPStatus,
		Code:       OK.Code,
		Message:    OK.Message,
		Data:       data,
	}
}

func decodeError(err error) Errno {
	if err == nil {
		return &OK
	}

	var ptr *errno
	if errors.As(err, &ptr) && ptr != nil {
		return ptr
	}

	var val errno
	if errors.As(err, &val) {
		return &val
	}

	return InternalServerError.WithData(err.Error())
}

// Decode decodes data and an error into an Errno struct.
func Decode(data any, err error) Errno {
	if err != nil {
		return decodeError(err)
	}

	return decodeData(data)
}
