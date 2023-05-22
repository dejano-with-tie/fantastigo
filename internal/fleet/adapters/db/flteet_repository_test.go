package db

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"regexp"
)

var MockFleet = app.Fleet{Id: "3545", Name: "Test fleet", Capacity: 10, VehicleTypes: []app.VehicleType{app.GetVehicleType("truck"),
	app.GetVehicleType("car")}}

func TestFleetSvc_CreateSQL(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("insert into fleet (id, name, capacity, vehicle_id) VALUES ($1, $2, $3, $4)")).
		WithArgs(MockFleet.Id, MockFleet.Name, MockFleet.Capacity, nil).
		WillReturnResult(sqlmock.NewResult(3545, 1))

	myDb := NewFleetRepository(db)

	err = myDb.Save(MockFleet)
	if err != nil {
		t.Errorf("something went wrong: %s", err.Error())
	}
}

func TestFleetSvc_GetSQL(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, name, capacity FROM fleet").
		WithArgs(MockFleet.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "capacity"}).
			AddRow(MockFleet.Id, MockFleet.Name, MockFleet.Capacity))

	myDb := NewFleetRepository(db)

	_, err = myDb.GetById(MockFleet.Id)
	if err != nil {
		t.Errorf("something went wrong: %s", err.Error())
	}
}
