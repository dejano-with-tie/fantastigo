package metrics

import (
	"sync"
)

type Store struct {
	items     map[string]Measurement
	itemsLock sync.RWMutex
}

func NewStore() *Store {
	s := &Store{
		items:     make(map[string]Measurement, 1000),
		itemsLock: sync.RWMutex{},
	}

	return s
}

func (s *Store) Get(query Query) ([]Measurement, error) {
	//TODO implement me
	s.itemsLock.RLock()
	defer s.itemsLock.RUnlock()
	var result []Measurement
	return append(result, s.items["gps_lat"]), nil
}

func (s *Store) Put(m Measurement) error {
	s.itemsLock.Lock()
	defer s.itemsLock.Unlock()
	if _, ok := s.items[m.Name]; !ok {
		s.items[m.Name] = m
	}
	return nil
}
