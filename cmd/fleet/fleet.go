package main

import (
	"github.com/dejano-with-tie/fantastigo/config"
	"github.com/dejano-with-tie/fantastigo/internal/common/server"
	repo "github.com/dejano-with-tie/fantastigo/internal/fleet/adapters/db"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	inner "github.com/dejano-with-tie/fantastigo/internal/fleet/server"
	"github.com/labstack/echo/v4"

	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config")
	}

	dbConn := executeMigrations()

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

func executeMigrations() *sql.DB {
	pgConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", pgConnection)
	if err != nil {
		panic("Unable to open db connection")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		panic("Unable to create driver instance")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		panic("Unable to create migration\n" + err.Error())
	}

	err = m.Up()
	if err != nil {
		log.Printf("Error executing migrations\n" + err.Error())
	}

	return db
}
