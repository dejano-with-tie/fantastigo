package db

import (
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"sync"
	"time"
)

// MemoryFleetRepository is an in memory implementation of fleet repository
type MemoryFleetRepository struct {
	storage     map[string]fleet
	storageLock *sync.RWMutex
}

// ensure interface is implemented at compile time
var _ app.FleetRepo = &MemoryFleetRepository{}

// NewFleetInMemoryRepository constructor
func NewFleetInMemoryRepository() *MemoryFleetRepository {
	return &MemoryFleetRepository{
		storage:     make(map[string]fleet),
		storageLock: &sync.RWMutex{},
	}
}

// GetById fleet from backing storage by id.
//
// It's kinda confusing to have a pointer receiver here. In case of value receiver, linter screams the following:
// "Struct MemoryFleetRepository has methods on both value and pointer receivers. Such usage is not recommended by the Go Documentation."
// Unfortunately, I wasn't able to find reasoning behind this in golang docs (mutable vs. immutable structs?). This should be further investigated.
func (m *MemoryFleetRepository) GetById(id string) (app.Fleet, error) {
	m.storageLock.RLock()
	defer m.storageLock.RUnlock()

	f, ok := m.storage[id]

	if !ok {
		return app.Fleet{}, fmt.Errorf("fleet with <id=%s> not found", id)
	}

	return app.Fleet{Id: f.Id}, nil
}

// Save fleet to a backing storage.
func (m *MemoryFleetRepository) Save(f app.Fleet) error {
	m.storageLock.Lock()
	defer m.storageLock.Unlock()

	created := fleet{
		Id:                  f.Id,
		Name:                f.Name,
		Capacity:            f.Capacity,
		AllowedVehicleTypes: f.VehicleTypes,
		Vehicles:            nil,
		CreatedAt:           time.Now().UTC(),
	}

	m.storage[created.Id] = created

	return nil
}
