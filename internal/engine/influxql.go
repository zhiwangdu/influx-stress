package engine

import (
	influxclient "github.com/influxdata/influx-stress/internal/influx"
	"github.com/influxdata/influx-stress/internal/report"
)

// InfluxqlStatement is a Statement Implementation that allows statements that parse in InfluxQL to be passed directly to the target instance
type InfluxqlStatement struct {
	StatementID string
	Query       string
	Tracer      *influxclient.Tracer
}

func (i *InfluxqlStatement) tags() map[string]string {
	tags := make(map[string]string)
	return tags
}

// SetID statisfies the Statement Interface
func (i *InfluxqlStatement) SetID(s string) {
	i.StatementID = s
}

// Run statisfies the Statement Interface
func (i *InfluxqlStatement) Run(s *influxclient.StressTest) {

	// Set the tracer
	i.Tracer = influxclient.NewTracer(i.tags())

	// Make the Package
	p := influxclient.NewPackage(influxclient.Query, []byte(i.Query), i.StatementID, i.Tracer)

	// Increment the tracer
	i.Tracer.Add(1)

	// Send the Package
	s.SendPackage(p)

	// Wait for all operations to finish
	i.Tracer.Wait()
}

// Report statisfies the Statement Interface
// No test coverage, fix
func (i *InfluxqlStatement) Report(s *influxclient.StressTest) (out string) {
	allData := s.GetStatementResults(i.StatementID, "query")

	return report.InfluxQL(i.Query, allData[0].Series[0].Columns, allData[0].Series[0].Values)
}
