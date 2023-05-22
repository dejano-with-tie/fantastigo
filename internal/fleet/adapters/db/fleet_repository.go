package db

import (
	"database/sql"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"log"
	sq "github.com/Masterminds/squirrel"
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

	err := sq.Select("id", "name", "capacity").From("fleet").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(m.db).
		QueryRow().
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
