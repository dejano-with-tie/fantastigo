package app

import (
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/common/apperr"
	"log"
)

type (
	FleetSvc struct {
		repo FleetRepo
	}
	VehicleSvc struct {
	}
	DriverSvc struct {
	}
)

func (s DriverSvc) Create(firstName, lastName string, category DrivingLicenseCategory) error {
	log.Default().Println(fmt.Sprintf("Create driver fname:%s, lname:%s, category: %s", firstName, lastName, category))
	return nil
}

func NewFleetService(fleetRepo FleetRepo) FleetSvc {
	if fleetRepo == nil {
		panic("nil fleetRepo")
	}

	return FleetSvc{repo: fleetRepo}
}

func (s FleetSvc) Create(name string, capacity int, allowedVehicleTypes []VehicleType) (id *string, err error) {
	log.Default().Println("test default logging")

	fleet, err := NewFleet(name, capacity, allowedVehicleTypes)
	if err != nil {
		return nil, apperr.Wrap("fleet:capacity-overflow", err)
	}

	if err := s.repo.Save(*fleet); err != nil {
		return nil, err
	}

	return &fleet.Id, err
}

func (s VehicleSvc) Create() error {
	return fmt.Errorf("fail message from a service")
}

func (s FleetSvc) GetFleet(id string) (f Fleet, e error) {

	fleet, err := s.repo.GetById(id)

	if err != nil {
		return Fleet{}, apperr.Wrap("fleet:not-found", err)
	}

	return fleet, nil
}
