package db

import (
	"database/sql"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"log"
)

type FleetRepository struct {
	db *sql.DB
}

func NewFleetRepository(dbConn *sql.DB) *FleetRepository {
	return &FleetRepository{
		db: dbConn,
	}
}

func (m *FleetRepository) GetById(id string) (app.Fleet, error) {

	fleet := app.Fleet{}

	err := m.db.QueryRow("SELECT id, name, capacity FROM fleet WHERE id = $1", id).
		Scan(&fleet.Id, &fleet.Name, &fleet.Capacity)

	if err != nil {
		log.Fatal(err)
	}

	return fleet, nil
}

func (m *FleetRepository) Save(f app.Fleet) error {

	_, err := m.db.Exec("insert into fleet (id, name, capacity, vehicle_id) VALUES ($1, $2, $3, $4)",
		f.Id, f.Name, f.Capacity, nil)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
