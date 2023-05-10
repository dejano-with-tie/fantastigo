package main

import (
	"github.com/dejano-with-tie/fantastigo/config"
	"github.com/dejano-with-tie/fantastigo/internal/common/server"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/adapters/db"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	inner "github.com/dejano-with-tie/fantastigo/internal/fleet/server"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config")
	}

	// prolly pass validator to the app.App itself
	v := validator.New()
	inner.RegisterValidations(v)

	fleetSvc := app.NewFleetService(db.NewFleetInMemoryRepository())
	vehicleSvc := app.VehicleSvc{}
	driverSvc := app.DriverSvc{}

	application := app.App{
		FleetSvc:   fleetSvc,
		VehicleSvc: vehicleSvc,
		DriverSvc:  driverSvc,
	}

	server.Run(cfg, v, func(e *echo.Echo) {
		g := e.Group("/api")
		inner.RegisterHandlers(g, inner.NewFleetHttpHandler(application))
		inner.RegisterDriverRoutes(g, inner.NewDriverHttpHandler(application))
	})
}
