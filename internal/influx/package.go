package influx

// Package enables runtime statements to send write or query payloads to the backend client.
// Packages carry either writes or queries in Body.
type Package struct {
	T           Type
	Body        []byte
	StatementID string
	Tracer      *Tracer
}

// NewPackage creates a new package with the appropriate payload
func NewPackage(t Type, body []byte, statementID string, tracer *Tracer) Package {
	p := Package{
		T:           t,
		Body:        body,
		StatementID: statementID,
		Tracer:      tracer,
	}

	return p
}
