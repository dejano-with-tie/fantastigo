package app

import (
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/common/apperr"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

const (
	Truck VehicleType = "truck"
	Bus   VehicleType = "bus"
	Van   VehicleType = "van"
	Car   VehicleType = "car"
)
const (
	A DrivingLicenseCategory = "a"
	B DrivingLicenseCategory = "b"
	C DrivingLicenseCategory = "c"
	D DrivingLicenseCategory = "d"
)

var (
	vehicleTypes = map[string]VehicleType{
		"truck": Truck,
		"bus":   Bus,
		"van":   Van,
		"car":   Car,
	}
	drivingLicenseCategory = map[string]DrivingLicenseCategory{
		"a": A,
		"b": B,
		"c": C,
		"d": D,
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
	Driver      struct {
		firstName       string
		lastName        string
		licenseCategory DrivingLicenseCategory
	}
	DrivingLicenseCategory string
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

func GetDrivingLicenceCategory(val string) (*DrivingLicenseCategory, error) {
	d := drivingLicenseCategory[val]
	if len(d) == 0 {
		return nil, fmt.Errorf("driver license type <value=%s>: error: %w", val, apperr.InvalidValueError)
	}
	return &d, nil
}
