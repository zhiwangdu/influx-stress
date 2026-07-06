package influx

// Directive enables SET statements to change backend client state.
type Directive struct {
	Property string
	Value    string
	Tracer   *Tracer
}

// NewDirective creates a new instance of a Directive with the appropriate state variable to change
func NewDirective(property string, value string, tracer *Tracer) Directive {
	d := Directive{
		Property: property,
		Value:    value,
		Tracer:   tracer,
	}
	return d
}
