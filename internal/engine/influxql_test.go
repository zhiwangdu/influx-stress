package engine

import (
	"testing"

	influxclient "github.com/influxdata/influx-stress/internal/influx"
)

func TestInfluxQlSetID(t *testing.T) {
	e := newTestInfluxQl()
	newID := "oaijnifo"
	e.SetID(newID)
	if e.StatementID != newID {
		t.Errorf("Expected: %v\nGot: %v\n", newID, e.StatementID)
	}
}

func TestInfluxQlRun(t *testing.T) {
	e := newTestInfluxQl()
	s, packageCh, _ := influxclient.NewTestStressTest()
	go func() {
		for pkg := range packageCh {
			if pkg.T != influxclient.Query {
				t.Errorf("Expected package to be Query\nGot: %v", pkg.T)
			}
			if string(pkg.Body) != e.Query {
				t.Errorf("Expected query: %v\nGot: %v", e.Query, string(pkg.Body))
			}
			if pkg.StatementID != e.StatementID {
				t.Errorf("Expected statementID: %v\nGot: %v", e.StatementID, pkg.StatementID)
			}
			pkg.Tracer.Done()
		}
	}()
	e.Run(s)
}

func newTestInfluxQl() *InfluxqlStatement {
	return &InfluxqlStatement{
		Query:       "CREATE DATABASE foo",
		Tracer:      influxclient.NewTracer(make(map[string]string)),
		StatementID: "fooID",
	}
}
