package server

import "github.com/go-playground/validator/v10"

// Validator is a custom echo validator
type Validator struct {
	delegate *validator.Validate
}

// NewValidator creates echo request Validator
func NewValidator() *Validator {
	return &Validator{delegate: validator.New()}
}

func (v Validator) Validate(i interface{}) error {
	return v.delegate.Struct(i)
}
