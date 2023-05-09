// Package db Implements repository pattern for interacting with the database
package db

import (
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"time"
)

type fleet struct {
	Id                  string
	Name                string
	Capacity            int
	AllowedVehicleTypes []app.VehicleType
	Vehicles            []vehicle
	CreatedAt           time.Time
}

type vehicle struct {
	Id   string
	Vin  string
	Type app.VehicleType
}

type driver struct {
	Id        string
	FirstName string
	LastName  string
}
