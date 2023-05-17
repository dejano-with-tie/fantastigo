package metrics

import (
	"fmt"
	"strings"
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
	MeasurementName string // matched by value.StartsWith. use '*' to match all
	MeasurementTags map[string]string
}

func (q Query) IsZero() bool {
	return len(q.MeasurementName) == 0 && len(q.MeasurementTags) == 0
}

type Measurement struct {
	Name      string
	Tags      MeasurementTag
	Value     any
	Timestamp time.Time
}

type MeasurementTag map[string]string

func (t MeasurementTag) flatten() string {
	sb := strings.Builder{}
	for k, v := range t {
		sb.WriteString(fmt.Sprintf("%s=%s;", k, v))
	}
	return sb.String()
}
