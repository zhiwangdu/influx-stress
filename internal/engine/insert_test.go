package engine

import (
	"strings"
	"testing"
	"time"

	influxclient "github.com/influxdata/influx-stress/internal/influx"
	"github.com/influxdata/influx-stress/internal/workload"
)

func TestInsertSetID(t *testing.T) {
	e := newTestInsert()
	newID := "oaijnifo"
	e.SetID(newID)
	if e.StatementID != newID {
		t.Errorf("Expected: %v\nGot: %v\n", newID, e.StatementID)
	}
}

func TestInsertRun(t *testing.T) {
	i := newTestInsert()
	s, packageCh, _ := influxclient.NewTestStressTest()
	// Listen to the other side of the directiveCh
	go func() {
		for pkg := range packageCh {
			countPoints := i.Timestamp.Count
			batchSize := s.BatchSize
			got := len(strings.Split(string(pkg.Body), "\n"))
			switch got {
			case countPoints % batchSize:
			case batchSize:
			default:
				t.Errorf("countPoints: %v\nbatchSize: %v\ngot: %v\n", countPoints, batchSize, got)
			}
			pkg.Tracer.Done()
		}
	}()
	i.Run(s)
}

func newTestInsert() *InsertStatement {
	return &InsertStatement{
		TestID:         "foo_test",
		StatementID:    "foo_ID",
		Name:           "foo_name",
		TemplateString: "cpu,%v %v %v",
		Timestamp:      &workload.Timestamp{Count: 20, Duration: time.Second},
		Templates: workload.Templates{
			{Tags: []string{"thing", "other_thing"}},
			{Function: &workload.Function{Type: "int", Fn: "inc", Argument: 0, Count: 0}},
		},
		TagCount: 1,
	}
}
