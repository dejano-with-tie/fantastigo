package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/dejano-with-tie/fantastigo/config"
	"github.com/dejano-with-tie/fantastigo/internal/common/server"
	repo "github.com/dejano-with-tie/fantastigo/internal/fleet/adapters/db"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	inner "github.com/dejano-with-tie/fantastigo/internal/fleet/server"
	"github.com/go-playground/validator/v10"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	var activeEnv string
	flag.StringVar(&activeEnv, "env", "", "active environment: dev, qa, prod")
	flag.Parse()
	cfg := config.MustLoad(config.NewLoadConfig(activeEnv))

	dbConn := executeMigrations(cfg.DB)

	// prolly pass validator to the app.App itself
	v := validator.New()
	inner.RegisterValidations(v)

	fleetSvc := app.NewFleetService(repo.NewFleetRepository(dbConn))
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

func executeMigrations(cfg config.DB) *sql.DB {
	pgConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)

	db, err := sql.Open("postgres", pgConnection)
	must("Unable to open db connection", err)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	must("Unable to create driver instance", err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	must("Unable to create migration", err)

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			must("Error executing migrations", err)
		}
		log.Println("db migrations: no change")
	}

	return db
}

func must(desc string, err error) {
	if err != nil {
		fmt.Println(desc)
		panic(err)
	}
}
