package server

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type DriverEndpoints interface {
	CreateDriver(ctx echo.Context) error
	GetDriver(ctx echo.Context) error
}

type CreateDriver struct {
	FirstName         string                 `json:"firstName" validate:"required,min=2,max=32"`
	LastName          string                 `json:"lastName" validate:"required,min=2,max=32"`
	LicenceCategory   DrivingLicenceCategory `json:"licenceCategory" validate:"required,max=1,oneof=a b c d"`              // validated using 'oneof' validator
	LicenceCategoryV1 string                 `json:"licenceCategoryV1" validate:"required,max=1,driving_licence_category"` // validated using custom validator
	LicenceCategoryV2 string                 `json:"licenceCategoryV2" validate:"required"`                                // validated in service layer
}

type DrivingLicenceCategory string // type doesn't do anything in the context of validation because underlying type is a string.

const (
	A DrivingLicenceCategory = "a"
	B DrivingLicenceCategory = "b"
	C DrivingLicenceCategory = "c"
	D DrivingLicenceCategory = "d"
)

func drivingLicenceCategoryEnumValidator(fl validator.FieldLevel) bool {
	raw := fl.Field().Interface().(string)
	val := map[string]DrivingLicenceCategory{
		"a": A,
		"b": B,
		"c": C,
		"d": D,
	}[raw]
	return len(val) > 0
}

func RegisterValidations(validator *validator.Validate) {
	tag := "driving_licence_category"
	err := validator.RegisterValidation(tag, drivingLicenceCategoryEnumValidator)
	if err != nil {
		panic(fmt.Errorf("failed to register validator <tag=%s>; error: %w", tag, err))
	}
}

func RegisterDriverRoutes(echo *echo.Group, endpoints DriverEndpoints) {
	echo.GET("/driver/:id", endpoints.GetDriver)
	echo.POST("/driver", endpoints.CreateDriver)
}
