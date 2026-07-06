package parse

import "github.com/influxdata/influx-stress/internal/workload"

// Statement is an IQL syntax tree node produced by this package.
type Statement interface {
	SetID(string)
}

// QueryStatement is the parsed form of an IQL QUERY statement.
type QueryStatement struct {
	StatementID    string
	Name           string
	TemplateString string
	Args           []string
	Count          int
}

// SetID sets the statement identifier.
func (s *QueryStatement) SetID(id string) { s.StatementID = id }

// InsertStatement is the parsed form of an IQL INSERT statement.
type InsertStatement struct {
	StatementID    string
	Name           string
	TemplateString string
	TagCount       int
	Timestamp      *workload.Timestamp
	Templates      workload.Templates
}

// SetID sets the statement identifier.
func (s *InsertStatement) SetID(id string) { s.StatementID = id }

// ExecStatement is the parsed form of an IQL EXEC statement.
type ExecStatement struct {
	StatementID string
	Script      string
}

// SetID sets the statement identifier.
func (s *ExecStatement) SetID(id string) { s.StatementID = id }

// SetStatement is the parsed form of an IQL SET statement.
type SetStatement struct {
	StatementID string
	Var         string
	Value       string
}

// SetID sets the statement identifier.
func (s *SetStatement) SetID(id string) { s.StatementID = id }

// WaitStatement is the parsed form of an IQL WAIT statement.
type WaitStatement struct {
	StatementID string
}

// SetID sets the statement identifier.
func (s *WaitStatement) SetID(id string) { s.StatementID = id }

// GoStatement is the parsed form of an IQL GO statement.
type GoStatement struct {
	StatementID string
	Statement   Statement
}

// SetID sets the statement identifier.
func (s *GoStatement) SetID(id string) { s.StatementID = id }
