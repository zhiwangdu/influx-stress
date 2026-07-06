package engine

import (
	"testing"

	influxclient "github.com/influxdata/influx-stress/internal/influx"
)

func TestQuerySetID(t *testing.T) {
	e := newTestQuery()
	newID := "oaijnifo"
	e.SetID(newID)
	if e.StatementID != newID {
		t.Errorf("Expected: %v\nGot: %v\n", newID, e.StatementID)
	}
}

func TestQueryRun(t *testing.T) {
	i := newTestQuery()
	s, packageCh, _ := influxclient.NewTestStressTest()
	// Listen to the other side of the directiveCh
	go func() {
		for pkg := range packageCh {
			if i.TemplateString != string(pkg.Body) {
				t.Fail()
			}
			pkg.Tracer.Done()
		}
	}()
	i.Run(s)
}

func newTestQuery() *QueryStatement {
	return &QueryStatement{
		StatementID:    "foo_ID",
		Name:           "foo_name",
		TemplateString: "SELECT count(value) FROM cpu",
		Args:           []string{},
		Count:          5,
		Tracer:         influxclient.NewTracer(map[string]string{}),
	}
}
