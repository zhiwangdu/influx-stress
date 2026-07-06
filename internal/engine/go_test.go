package engine

import (
	"testing"

	influxclient "github.com/influxdata/influx-stress/internal/influx"
)

func TestGoSetID(t *testing.T) {
	e := newTestGo()
	newID := "oaijnifo"
	e.SetID(newID)
	if e.StatementID != newID {
		t.Errorf("Expected: %v\nGot: %v\n", newID, e.StatementID)
	}
}

func TestGoRun(t *testing.T) {
	e := newTestGo()
	s, _, _ := influxclient.NewTestStressTest()
	e.Run(s)
	if e == nil {
		t.Fail()
	}
}

func TestGoReport(t *testing.T) {
	e := newTestGo()
	s, _, _ := influxclient.NewTestStressTest()
	report := e.Report(s)
	if report != "Go " {
		t.Errorf("Expected: %v\nGot: %v\n", "Go ", report)
	}
}

func newTestGo() *GoStatement {
	return &GoStatement{
		Statement:   newTestExec(),
		StatementID: "fooID",
	}
}
