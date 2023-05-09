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
)

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
