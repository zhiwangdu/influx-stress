package engine

import (
	"strings"
	"testing"

	influxclient "github.com/influxdata/influx-stress/internal/influx"
)

func TestWaitSetID(t *testing.T) {
	e := newTestWait()
	newID := "oaijnifo"
	e.SetID(newID)
	if e.StatementID != newID {
		t.Errorf("Expected: %v\ngott: %v\n", newID, e.StatementID)
	}
}

func TestWaitRun(t *testing.T) {
	e := newTestWait()
	s, _, _ := influxclient.NewTestStressTest()
	e.Run(s)
	if e == nil {
		t.Fail()
	}
}

func TestWaitReport(t *testing.T) {
	e := newTestWait()
	s, _, _ := influxclient.NewTestStressTest()
	rpt := e.Report(s)
	if !strings.Contains(rpt, "WAIT") {
		t.Fail()
	}
}

func newTestWait() *WaitStatement {
	return &WaitStatement{
		StatementID: "fooID",
	}
}
