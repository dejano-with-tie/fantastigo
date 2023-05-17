package metrics

import (
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/common/util/collection"
	"golang.org/x/exp/maps"
	"strings"
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
	if query.IsZero() {
		return []Measurement{}, nil
	}

	s.itemsLock.RLock()
	defer s.itemsLock.RUnlock()

	var results = make([]Measurement, 0, len(s.items))
	qmn := query.MeasurementName
	qmt := query.MeasurementTags

	// filter by measurement name
	if len(qmn) > 0 && qmn != "*" {
		results = s.find(qmn)
	} else {
		results = maps.Values(s.items)
	}

	// filter by tags
	if len(qmt) > 0 {
		results = collection.Filter(results, func(m Measurement) bool {
			for k, v := range qmt {
				if strings.Contains(m.Tags.flatten(), fmt.Sprintf("%s=%v", k, v)) {
					return true
				}
			}
			return false
		})
	}

	return results, nil
}

func (s *Store) Put(m Measurement) error {
	s.itemsLock.Lock()
	defer s.itemsLock.Unlock()
	if _, ok := s.items[m.Name]; !ok {
		s.items[m.Name] = m
	}
	return nil
}

func (s *Store) find(key string) []Measurement {
	s.itemsLock.RLock()
	defer s.itemsLock.RUnlock()

	values := maps.Values(s.items)
	return collection.Filter(values, func(m Measurement) bool {
		return strings.HasPrefix(m.Name, key)
	})
}
