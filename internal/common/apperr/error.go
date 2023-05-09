package apperr

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

type AppErr struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (h *AppErr) Error() string {
	return fmt.Sprintf("%s: %s", h.Code, h.Message)
}

func New(code string, msg string) *AppErr {
	return &AppErr{
		Code:    code,
		Message: msg,
	}
}

func Wrap(code string, error error) *AppErr {
	if error == nil {
		return New(code, "")
	}

	return &AppErr{
		Code:    code,
		Message: error.Error(),
		Err:     error,
	}
}
