package app

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

const (
	Truck VehicleType = "truck"
	Bus   VehicleType = "bus"
	Van   VehicleType = "van"
	Car   VehicleType = "car"
)

var (
	vehicleTypes = map[string]VehicleType{
		"truck": Truck,
		"bus":   Bus,
		"van":   Van,
		"car":   Car,
	}
)

type (
	Fleet struct {
		Id           string
		Name         string
		Capacity     int
		VehicleTypes []VehicleType
	}
	FleetRepo interface {
		GetById(id string) (Fleet, error)
		Save(f Fleet) error
	}
	VehicleType string
)

func NewFleet(name string, capacity int, vehicleTypes []VehicleType) (*Fleet, error) {
	if slices.Contains(vehicleTypes, Truck) && capacity > 30 {
		return nil, fmt.Errorf("max capacity for a fleet of trucks is 30")
	}

	return &Fleet{
		Id:           uuid.New().String(),
		Name:         name,
		Capacity:     capacity,
		VehicleTypes: vehicleTypes,
	}, nil
}

func GetVehicleType(val string) VehicleType {
	return vehicleTypes[val]
}
