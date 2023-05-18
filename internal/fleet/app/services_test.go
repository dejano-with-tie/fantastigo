package app

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type MockFleetRepository struct {
	mock.Mock
}

func (m *MockFleetRepository) GetById(id string) (Fleet, error) {
	args := m.Called() // loads arguments specified in a Return statement when mocking call
	f := args.Get(0)
	// Get function returns interface, so we need Type assertion here to verify that argument holds value of type Fleet
	// this expression (Type assertion) makes assertion and also assigns (in this case returns) concrete value of f
	return f.(Fleet), args.Error(1)
}

func (m *MockFleetRepository) Save(fleet Fleet) error {
	args := m.Called()
	return args.Error(0)
}

func TestFleetSvc_Create(t *testing.T) {

	mockRepo := new(MockFleetRepository)

	mockRepo.On("Save").Return(nil) // mocking call to a repo function

	fleetSvc := NewFleetService(mockRepo)
	id, err := fleetSvc.Create("Test fleet", 10, []VehicleType{GetVehicleType("truck"),
		GetVehicleType("car")})

	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
}

func TestFleetSvc_GetFleet(t *testing.T) {

	mockRepo := new(MockFleetRepository)
	mockId := "3545"
	mockFleet := Fleet{mockId, "Test fleet", 10, []VehicleType{GetVehicleType("truck"),
		GetVehicleType("car")}}

	mockRepo.On("GetById").Return(mockFleet, nil)

	fleetSvc := NewFleetService(mockRepo)
	fleet, err := fleetSvc.GetFleet(mockId)

	mockRepo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, mockId, fleet.Id)
	assert.Equal(t, "Test fleet", fleet.Name)
	assert.Equal(t, 10, fleet.Capacity)
	assert.Equal(t, 2, len(fleet.VehicleTypes))
}
