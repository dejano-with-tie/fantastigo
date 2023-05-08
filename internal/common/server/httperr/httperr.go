package httperr

import (
	"fmt"
)

const (
	ErrCodeUnknown        = "unknown"
	ErrCodeInternal       = "internal"
	ErrCodeBadRequest     = "bad-request"
	ErrCodeNotfound       = "not-found"
	ErrCodeConflict       = "conflict"
	ErrCodeNotImplemented = "not-implemented"
	ErrCodeUnauthorized   = "unauthorized"
)

type HttpErr struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (h *HttpErr) Error() string {
	return fmt.Sprintf("%s: %s", h.Code, h.Message)
}

func New(code string, msg string) *HttpErr {
	return &HttpErr{
		Code:    code,
		Message: msg,
	}
}

func Wrap(code string, error error) *HttpErr {
	if error == nil {
		return New(code, "")
	}

	return &HttpErr{
		Code:    code,
		Message: error.Error(),
		Err:     error,
	}
}
