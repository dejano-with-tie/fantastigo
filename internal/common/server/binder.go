package server

import "github.com/labstack/echo/v4"

// Binder is a custom echo binder
type Binder struct {
	delegate echo.Binder
}

func NewBinder() *Binder {
	return &Binder{delegate: &echo.DefaultBinder{}}
}

// Bind tries to bind request into interface, and if it completed successfully apply validation
// https://echo.labstack.com/guide/request/#custom-binder
func (b *Binder) Bind(i interface{}, c echo.Context) error {
	if err := b.delegate.Bind(i, c); err != nil && err != echo.ErrUnsupportedMediaType {
		return err
	}
	return c.Validate(i)
}
