package iql

import (
	"reflect"
	"testing"
	"time"

	"github.com/influxdata/influx-stress/internal/engine"
	"github.com/influxdata/influx-stress/internal/iql/parse"
	"github.com/influxdata/influx-stress/internal/workload"
)

// Pulls the default configFile and makes sure it parses
func TestParseStatements(t *testing.T) {
	stmts, err := ParseStatements("../../examples/iql/file.iql")
	if err != nil {
		t.Error(err)
	}
	expected := 15
	got := len(stmts)
	if expected != got {
		t.Errorf("expected: %v\ngot: %v\n", expected, got)
	}
}

func TestBuildStatementConvertsAST(t *testing.T) {
	timestamp := &workload.Timestamp{Count: 100, Duration: time.Second}
	templates := workload.Templates{
		&workload.Template{Tags: []string{"west", "east"}},
		&workload.Template{Function: &workload.Function{Type: "int", Fn: "rand", Argument: 10, Count: 2}},
	}

	tests := []struct {
		name string
		in   parse.Statement
		want engine.Statement
	}{
		{
			name: "query",
			in: &parse.QueryStatement{
				StatementID:    "q1",
				Name:           "basicCount",
				TemplateString: "SELECT count(%v) FROM cpu",
				Args:           []string{"%f"},
				Count:          10,
			},
			want: &engine.QueryStatement{
				StatementID:    "q1",
				Name:           "basicCount",
				TemplateString: "SELECT count(%v) FROM cpu",
				Args:           []string{"%f"},
				Count:          10,
			},
		},
		{
			name: "insert",
			in: &parse.InsertStatement{
				StatementID:    "i1",
				Name:           "mockCpu",
				TemplateString: "cpu,host=%v value=%v %v",
				TagCount:       1,
				Timestamp:      timestamp,
				Templates:      templates,
			},
			want: &engine.InsertStatement{
				StatementID:    "i1",
				Name:           "mockCpu",
				TemplateString: "cpu,host=%v value=%v %v",
				TagCount:       1,
				Timestamp:      timestamp,
				Templates:      templates,
			},
		},
		{
			name: "exec",
			in:   &parse.ExecStatement{StatementID: "e1", Script: "load.sh"},
			want: &engine.ExecStatement{StatementID: "e1", Script: "load.sh"},
		},
		{
			name: "set",
			in:   &parse.SetStatement{StatementID: "s1", Var: "database", Value: "stress"},
			want: &engine.SetStatement{StatementID: "s1", Var: "database", Value: "stress"},
		},
		{
			name: "wait",
			in:   &parse.WaitStatement{StatementID: "w1"},
			want: &engine.WaitStatement{StatementID: "w1"},
		},
		{
			name: "go",
			in: &parse.GoStatement{
				StatementID: "g1",
				Statement:   &parse.QueryStatement{StatementID: "q2", Name: "basicCount", TemplateString: "SELECT count(free) FROM cpu", Count: 5},
			},
			want: &engine.GoStatement{
				StatementID: "g1",
				Statement:   &engine.QueryStatement{StatementID: "q2", Name: "basicCount", TemplateString: "SELECT count(free) FROM cpu", Count: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildStatement(tt.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("wrong statement:\nexpected %#v\ngot      %#v", tt.want, got)
			}
		})
	}
}
