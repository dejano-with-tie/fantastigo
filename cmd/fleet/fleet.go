package main

import (
	"github.com/dejano-with-tie/fantastigo/config"
	"github.com/dejano-with-tie/fantastigo/internal/common/server"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/adapters/db"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	fleet "github.com/dejano-with-tie/fantastigo/internal/fleet/server"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config")
	}

	fs := app.NewFleetService(db.NewFleetInMemoryRepository())
	vs := app.VehicleSvc{}
	application := app.App{
		FleetSvc:   fs,
		VehicleSvc: vs,
	}

	server.Run(cfg, func(e *echo.Echo) {
		fleet.RegisterHandlers(e.Group("/api"), fleet.NewHttpHandler(application))
	})
}
