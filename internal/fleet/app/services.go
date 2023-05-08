package app

import (
	"fmt"
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

func (s FleetSvc) Create(fleet Fleet) (id string, err error) {
	log.Default().Println("test default logging")
	created, err := s.repo.Save(fleet)
	return created.Id, err
}

func (s VehicleSvc) Create() error {
	return fmt.Errorf("fail message from a service")
}
