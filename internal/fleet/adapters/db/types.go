// Package db Implements repository pattern for interacting with the database
package db

type VehicleTye int

const (
	Truck VehicleTye = iota
	Bus
	Van
	Car
)

type Fleet struct {
	Id                  string
	Name                string
	Capacity            int
	AllowedVehicleTypes []VehicleTye
	Vehicles            []Vehicle
}

type Vehicle struct {
	Id   string
	Vin  string
	Type VehicleTye
}

type Driver struct {
	Id        string
	FirstName string
	LastName  string
}
