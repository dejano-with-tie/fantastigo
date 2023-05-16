package app

import (
	"errors"
	"fmt"
	"github.com/dejano-with-tie/fantastigo/internal/ector/metrics"
)

var ErrorNoIdentity = errors.New("failed to fetch identity")

type Vehicle struct {
	Vin string
}

type Ector struct {
	metrics *metrics.Metrics
}

func NewEctor(storer metrics.Storer) *Ector {
	m := metrics.NewMetrics(storer)
	_ = metrics.NewPrometheusSourcer(storer)
	return &Ector{metrics: m}
}

type Status struct {
	internetConnection bool
}

type Identity struct {
	Vin string
}

func (e *Ector) Status() (*Vehicle, error) {
	return nil, nil
}

func (e *Ector) Tmp() {
	r, _ := e.metrics.Storer.Get(metrics.Query{})
	fmt.Println(r)
}

func (e *Ector) GetIdentity() (*Vehicle, error) {
	q := metrics.Query{
		MeasurementName: "system_uptime",
	}
	m, _ := e.metrics.Storer.Get(q)

	// TODO: Try to implement retry 4fun
	if len(m) == 0 {
		return nil, ErrorNoIdentity
	}

	i, ok := m[0].Tags["thingName"]
	if !ok {
		return nil, ErrorNoIdentity
	}

	return &Vehicle{Vin: i}, nil
}
