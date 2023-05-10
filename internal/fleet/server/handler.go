package server

import (
	"errors"
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/common/apperr"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	// FleetHttpHandler implements fleet http server handlers
	FleetHttpHandler struct {
		app app.App
	}
	// DriverHttpHandler implements driver http server handlers
	DriverHttpHandler struct {
		app app.App
	}
)

func NewDriverHttpHandler(app app.App) *DriverHttpHandler {
	return &DriverHttpHandler{app: app}
}

func (h FleetHttpHandler) CreateFleet(c echo.Context) error {
	c.Logger().Debug("test debug log level")
	r := CreateFleet{}

	if err := c.Bind(&r); err != nil {
		return err
	}

	id, err := h.app.FleetSvc.Create(r.Name, r.Capacity, mapVehicleTypes(r.VehicleTypes))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, IdResponse{Id: *id})
}

func mapVehicleTypes(types []CreateFleetVehicleTypes) []app.VehicleType {
	r := make([]app.VehicleType, len(types))
	for _, v := range types {
		r = append(r, app.GetVehicleType(string(v)))
	}
	return r
}

func (h FleetHttpHandler) GetFleet(c echo.Context) error {
	return errors.New("should respond with status=500, code=unknown")
}

func (h FleetHttpHandler) CreateVehicle(c echo.Context) error {
	return apperr.New(apperr.ErrCodeNotImplemented, "Should respond with status=501 and code=not-implemented")
}

func (h FleetHttpHandler) GetVehicle(c echo.Context, id string) error {
	if e := h.app.VehicleSvc.Create(); e != nil {
		return apperr.Wrap("business-error-code", fmt.Errorf("should respond with 422 and wrapped error: %w", e))
	}
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func NewFleetHttpHandler(app app.App) *FleetHttpHandler {
	return &FleetHttpHandler{app: app}
}

func (d DriverHttpHandler) CreateDriver(ctx echo.Context) error {
	var r = &CreateDriver{}
	if err := ctx.Bind(r); err != nil {
		return err
	}

	dlc, err := app.GetDrivingLicenceCategory(r.LicenceCategoryV2)
	if err != nil {
		return err
	}

	return d.app.DriverSvc.Create(r.FirstName, r.LastName, *dlc)
}

func (d DriverHttpHandler) GetDriver(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
