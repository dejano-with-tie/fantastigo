package server

import (
	"errors"
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/common/apperr"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"runtime/debug"
)

var httpErrorStatuses = map[string]int{
	apperr.ErrCodeUnknown:        http.StatusInternalServerError,
	apperr.ErrCodeInternal:       http.StatusInternalServerError,
	apperr.ErrCodeBadRequest:     http.StatusBadRequest,
	apperr.ErrCodeNotfound:       http.StatusNotFound,
	apperr.ErrCodeConflict:       http.StatusConflict,
	apperr.ErrCodeNotImplemented: http.StatusNotImplemented,
	apperr.ErrCodeUnauthorized:   http.StatusUnauthorized,
}

type (
	HttpErrResponse struct {
		Status      int                       `json:"-"`
		Code        string                    `json:"code"`
		Validations []ValidationFieldResponse `json:"validations,omitempty"`
		Message     string                    `json:"message,omitempty"`
		Trace       string                    `json:"trace,omitempty"`
	}
	ValidationFieldResponse struct {
		Property      string `json:"property"`
		Key           string `json:"error"`
		RejectedValue string `json:"rejectedValue"`
		Message       string `json:"message,omitempty"`
	}
)

func getHttpStatus(code string) int {
	status := httpErrorStatuses[code]
	if status == 0 {
		return http.StatusUnprocessableEntity
	}

	return status
}

func errHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	// default
	status := getHttpStatus(apperr.ErrCodeUnknown)
	code := apperr.ErrCodeUnknown
	message := err.Error()
	var validations []ValidationFieldResponse

	// clean up & try to rewrite this with errors.Is
	var appErr *apperr.AppErr
	var echoHttpErr *echo.HTTPError
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, v := range validationErrs {
			vfe := ValidationFieldResponse{
				Property:      v.Field(),
				Key:           v.ActualTag(),
				Message:       fmt.Sprintf("Property validation for '%s' failed validation tag '%s' tag", v.Field(), v.ActualTag()),
				RejectedValue: fmt.Sprintf("%v", v.Value()),
			}
			validations = append(validations, vfe)
		}
		status = getHttpStatus(apperr.ErrCodeBadRequest)
		code = apperr.ErrCodeBadRequest
		message = "Request validation failed"
	} else if errors.As(err, &appErr) {
		status = getHttpStatus(appErr.Code)
		code = appErr.Code
		message = appErr.Message
	} else if errors.As(err, &echoHttpErr) {
		status = echoHttpErr.Code
		message = echoHttpErr.Error()
	}

	trace := ""
	if c.Echo().Debug {
		debug.PrintStack()
		trace = string(debug.Stack())
	}

	he := &HttpErrResponse{
		Status:      status,
		Code:        code,
		Message:     message,
		Validations: validations,
		Trace:       trace,
	}

	// Send response
	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Status)
	} else {
		err = c.JSON(he.Status, he)
	}

	if err != nil {
		c.Logger().Error(err)
	}
}
