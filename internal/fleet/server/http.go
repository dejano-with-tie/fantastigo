package server

import (
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
	return c.JSON(http.StatusOK, struct {
		ID string `json:"id"`
	}{ID: "hello"})
}

func (h HttpHandler) CreateVehicle(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h HttpHandler) GetVehicle(c echo.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
