package metrics

import (
	"time"
)

type Metrics struct {
	Storer Storer
}

func NewMetrics(storer Storer) *Metrics {
	return &Metrics{Storer: storer}
}

type Storer interface {
	Get(query Query) ([]Measurement, error)
	Put(m Measurement) error
}

type NopStorer struct {
}

func (n *NopStorer) Put(m Measurement) error {
	return nil
}

func (n *NopStorer) Get(query Query) ([]Measurement, error) {
	//fmt.Println("noop storer get(query)")
	return nil, nil
}

type Query struct {
	MeasurementName string
	MeasurementTags map[string]string
}

type Measurement struct {
	Name      string
	Tags      map[string]string
	Value     any
	Timestamp time.Time
}
