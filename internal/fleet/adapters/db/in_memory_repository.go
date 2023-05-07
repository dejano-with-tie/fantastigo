package db

import (
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/fleet/app"
	"github.com/google/uuid"
	"sync"
)

// MemoryFleetRepository is an in memory implementation of fleet repository
type MemoryFleetRepository struct {
	storage     map[string]Fleet
	storageLock *sync.RWMutex
}

// NewFleetInMemoryRepository constructor
func NewFleetInMemoryRepository() *MemoryFleetRepository {
	return &MemoryFleetRepository{
		storage:     make(map[string]Fleet),
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
func (m *MemoryFleetRepository) Save(f app.Fleet) (app.Fleet, error) {
	m.storageLock.Lock()
	defer m.storageLock.Unlock()
	// TODO via factory
	created := Fleet{
		Id:                  uuid.New().String(),
		Name:                f.Name,
		Capacity:            f.Capacity,
		AllowedVehicleTypes: nil,
		Vehicles:            nil,
	}

	m.storage[created.Id] = created

	f.Id = created.Id

	return f, nil
}
