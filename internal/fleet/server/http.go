package server

import (
	"errors"
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/common/server/httperr"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"github.com/labstack/echo/v4"
	"net/http"
)

// HttpHandler implements http server handlers
type HttpHandler struct {
	app app.App
}

func NewHttpHandler(app app.App) *HttpHandler {
	return &HttpHandler{app: app}
}

func (h HttpHandler) CreateFleet(c echo.Context) error {
	c.Logger().Debug("test debug log level")
	r := CreateFleet{}

	if err := c.Bind(&r); err != nil {
		return err
	}

	fleet := app.Fleet{
		Name:     r.Name,
		Capacity: r.Capacity,
	}

	id, err := h.app.FleetSvc.Create(fleet)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, IdResponse{Id: id})
}

func (h HttpHandler) GetFleet(c echo.Context) error {
	return errors.New("should respond with status=500, code=unknown")
}

func (h HttpHandler) CreateVehicle(c echo.Context) error {
	return httperr.New(httperr.ErrCodeNotImplemented, "Should respond with status=501 and code=not-implemented")
}

func (h HttpHandler) GetVehicle(c echo.Context, id string) error {
	if e := h.app.VehicleSvc.Create(); e != nil {
		return httperr.Wrap("business-error-code", fmt.Errorf("should respond with 422 and wrapped error: %w", e))
	}
	return echo.NewHTTPError(http.StatusNotImplemented)
}
