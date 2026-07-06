package statement

import (
	"strings"
	"testing"
)

func TestNewStrRandStringer(t *testing.T) {
	function := newStrRandFunction()
	strRandStringer := function.NewStringer(10)
	s := strRandStringer()
	if len(s) != function.Argument {
		t.Errorf("Expected: %v\nGot: %v\n", function.Argument, len(s))
	}
}

func TestStrRandLengthAndCharset(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{n: 0, want: 0},
		{n: 1, want: 0},
		{n: 7, want: 6},
		{n: 8, want: 8},
	}

	for _, tt := range tests {
		got := randStr(tt.n)()
		if len(got) != tt.want {
			t.Errorf("randStr(%d) length: expected %d, got %d", tt.n, tt.want, len(got))
		}
		for _, ch := range got {
			if !strings.ContainsRune("0123456789abcdef", ch) {
				t.Errorf("randStr(%d) emitted non-hex rune %q in %q", tt.n, ch, got)
			}
		}
		if strings.ContainsAny(got, ", =\"") {
			t.Errorf("randStr(%d) emitted line-protocol metacharacter in %q", tt.n, got)
		}
	}
}

func TestNewIntIncStringer(t *testing.T) {
	function := newIntIncFunction()
	intIncStringer := function.NewStringer(10)
	s := intIncStringer()
	if s != "0i" {
		t.Errorf("Expected: 0i\nGot: %v\n", s)
	}
}

func TestNewIntRandStringer(t *testing.T) {
	function := newIntRandFunction()
	intRandStringer := function.NewStringer(10)
	s := intRandStringer()
	if parseInt(s[:len(s)-1]) > function.Argument {
		t.Errorf("Expected value below: %v\nGot value: %v\n", function.Argument, s)
	}
}

func TestNewFloatIncStringer(t *testing.T) {
	function := newFloatIncFunction()
	floatIncStringer := function.NewStringer(10)
	s := floatIncStringer()
	assertIntegerFloatToken(t, s)
	if parseFloat(s) != function.Argument {
		t.Errorf("Expected value: %v\nGot: %v\n", function.Argument, s)
	}
}
func TestNewFloatRandStringer(t *testing.T) {
	function := newFloatRandFunction()
	floatRandStringer := function.NewStringer(10)
	s := floatRandStringer()
	assertIntegerFloatToken(t, s)
	if parseFloat(s) > function.Argument {
		t.Errorf("Expected value below: %v\nGot value: %v\n", function.Argument, s)
	}
}

func TestFunctionAppenderCardinality(t *testing.T) {
	cycled := (&Function{Type: "int", Fn: "inc", Argument: 0, Count: 3}).NewAppender(10)
	for _, want := range []string{"0i", "1i", "2i", "0i", "1i"} {
		if got := string(cycled(nil)); got != want {
			t.Fatalf("cycled value: expected %q, got %q", want, got)
		}
	}

	repeated := (&Function{Type: "int", Fn: "inc", Argument: 0, Count: 0}).NewAppender(3)
	for _, want := range []string{"0i", "0i", "0i", "1i", "1i"} {
		if got := string(repeated(nil)); got != want {
			t.Fatalf("repeated value: expected %q, got %q", want, got)
		}
	}
}

func TestRandZeroPanicBehavior(t *testing.T) {
	mustPanic(t, "randInt(0)", func() { randInt(0)() })
	mustPanic(t, "randFloat(0)", func() { randFloat(0)() })
	mustNotPanic(t, "incInt(0)", func() { incInt(0)() })
	mustNotPanic(t, "incFloat(0)", func() { incFloat(0)() })
}

func TestStringersEval(t *testing.T) {
	// Make the *Function(s)
	strRandFunction := newStrRandFunction()
	intIncFunction := newIntIncFunction()
	intRandFunction := newIntRandFunction()
	floatIncFunction := newFloatIncFunction()
	floatRandFunction := newFloatRandFunction()
	// Make the *Stringer(s)
	strRandStringer := strRandFunction.NewStringer(10)
	intIncStringer := intIncFunction.NewStringer(10)
	intRandStringer := intRandFunction.NewStringer(10)
	floatIncStringer := floatIncFunction.NewStringer(10)
	floatRandStringer := floatRandFunction.NewStringer(10)
	// Make the *Stringers
	stringers := Stringers([]Stringer{strRandStringer, intIncStringer, intRandStringer, floatIncStringer, floatRandStringer})
	// Spoff the Time function
	// Call *Stringers.Eval
	values := stringers.Eval(spoofTime)
	// Check the strRandFunction
	if len(values[0].(string)) != strRandFunction.Argument {
		t.Errorf("Expected: %v\nGot: %v\n", strRandFunction.Argument, len(values[0].(string)))
	}
	// Check the intIncFunction
	if values[1].(string) != "0i" {
		t.Errorf("Expected: 0i\nGot: %v\n", values[1].(string))
	}
	// Check the intRandFunction
	s := values[2].(string)
	if parseInt(s[:len(s)-1]) > intRandFunction.Argument {
		t.Errorf("Expected value below: %v\nGot value: %v\n", intRandFunction.Argument, s)
	}
	// Check the floatIncFunction
	assertIntegerFloatToken(t, values[3].(string))
	if parseFloat(values[3].(string)) != floatIncFunction.Argument {
		t.Errorf("Expected value: %v\nGot: %v\n", floatIncFunction.Argument, values[3])
	}
	// Check the floatRandFunction
	assertIntegerFloatToken(t, values[4].(string))
	if parseFloat(values[4].(string)) > floatRandFunction.Argument {
		t.Errorf("Expected value below: %v\nGot value: %v\n", floatRandFunction.Argument, values[4])
	}
	// Check the spoofTime func
	if values[5] != 8 {

	}
}

func spoofTime() int64 {
	return int64(8)
}

func newStrRandFunction() *Function {
	return &Function{
		Type:     "str",
		Fn:       "rand",
		Argument: 8,
		Count:    1000,
	}
}

func newIntIncFunction() *Function {
	return &Function{
		Type:     "int",
		Fn:       "inc",
		Argument: 0,
		Count:    0,
	}
}

func newIntRandFunction() *Function {
	return &Function{
		Type:     "int",
		Fn:       "rand",
		Argument: 100,
		Count:    1000,
	}
}

func newFloatIncFunction() *Function {
	return &Function{
		Type:     "float",
		Fn:       "inc",
		Argument: 0,
		Count:    1000,
	}
}

func newFloatRandFunction() *Function {
	return &Function{
		Type:     "float",
		Fn:       "rand",
		Argument: 100,
		Count:    1000,
	}
}

func assertIntegerFloatToken(t *testing.T, s string) {
	t.Helper()
	if s == "" {
		t.Fatal("float token is empty")
	}
	if strings.ContainsAny(s, ".i") {
		t.Fatalf("float token should be integer-looking without '.' or 'i': %q", s)
	}
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			t.Fatalf("float token should contain only digits: %q", s)
		}
	}
}

func mustPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatalf("%s did not panic", name)
		}
	}()
	fn()
}

func mustNotPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s panicked: %v", name, r)
		}
	}()
	fn()
}
