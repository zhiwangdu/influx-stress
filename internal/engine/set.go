package engine

import (
	"fmt"
	"strings"

	influxclient "github.com/influxdata/influx-stress/internal/influx"
)

// SetStatement set state variables for the test
type SetStatement struct {
	Var   string
	Value string

	StatementID string

	Tracer *influxclient.Tracer
}

// SetID statisfies the Statement Interface
func (i *SetStatement) SetID(s string) {
	i.StatementID = s
}

// Run statisfies the Statement Interface
func (i *SetStatement) Run(s *influxclient.StressTest) {
	i.Tracer = influxclient.NewTracer(make(map[string]string))
	d := influxclient.NewDirective(strings.ToLower(i.Var), strings.ToLower(i.Value), i.Tracer)
	switch d.Property {
	// Needs to be set on both StressTest and client
	// Set the write percison for points generated
	case "precision":
		s.Precision = d.Value
		i.Tracer.Add(1)
		s.SendDirective(d)
	// Lives on StressTest
	// Set the date for the first point entered into the database
	case "startdate":
		s.Lock()
		s.StartDate = d.Value
		s.Unlock()
	// Lives on StressTest
	// Set the BatchSize for writes
	case "batchsize":
		s.Lock()
		s.BatchSize = parseInt(d.Value)
		s.Unlock()
	// All other variables live on client
	default:
		i.Tracer.Add(1)
		s.SendDirective(d)
	}
	i.Tracer.Wait()
}

// Report statisfies the Statement Interface
func (i *SetStatement) Report(s *influxclient.StressTest) string {
	return fmt.Sprintf("SET %v = '%v'", i.Var, i.Value)
}
