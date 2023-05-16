package server

import (
	"github.com/dejano-with-tie/fantastigo/internal/ector/app"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EctorHttpHandler struct {
	app *app.Ector
}

var _ ServerInterface = &EctorHttpHandler{}

func NewEctorHttpHandler(app *app.Ector) *EctorHttpHandler {
	return &EctorHttpHandler{app: app}
}

func (e EctorHttpHandler) GetIdentity(ctx echo.Context) error {
	//TODO implement me
	i, err := e.app.GetIdentity()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, i)
}

func (e EctorHttpHandler) GetStatus(ctx echo.Context) error {
	info, err := e.app.Status()
	if err != nil {
		return err
	}

	r := VehicleStatusResponse{Vin: info.Vin}

	return ctx.JSON(http.StatusOK, &r)
}

func (e EctorHttpHandler) GetMetrics(c echo.Context) error {
	//TODO implement me
	return nil
}
