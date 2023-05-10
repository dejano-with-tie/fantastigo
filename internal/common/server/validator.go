package server

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// Validator is a custom echo validator
type Validator struct {
	d *validator.Validate
}

// NewValidator creates echo request Validator
func NewValidator(delegate *validator.Validate) *Validator {
	return &Validator{d: delegate}
}

func (v *Validator) Validate(i interface{}) error {
	return v.d.Struct(i)
}

func (v *Validator) Register(tag string, fn func(fl validator.FieldLevel) bool) {
	err := v.d.RegisterValidation(tag, fn)
	if err != nil {
		panic(fmt.Errorf("failed to register validator <tag=%s>; error: %w", tag, err))
	}
}
