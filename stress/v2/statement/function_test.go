package statement

import (
	"strconv"
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
	mustPanic(t, "randf(0)", func() { randDecimalFloatAppender(0)(nil) })
	mustPanic(t, "zipfInt(0)", func() { zipfIntAppender(0)(nil) })
	mustPanic(t, "zipfFloat(0)", func() { zipfFloatAppender(0)(nil) })
	mustNotPanic(t, "incInt(0)", func() { incInt(0)() })
	mustNotPanic(t, "incFloat(0)", func() { incFloat(0)() })
	mustNotPanic(t, "normalFloat(0)", func() { normalFloatAppender(0)(nil) })
	mustNotPanic(t, "sinFloat(0)", func() { sinFloatAppender(0)(nil) })
	mustNotPanic(t, "walkFloat(0)", func() { walkFloatAppender(0)(nil) })
	mustNotPanic(t, "hashStr(0)", func() { hashStrAppender(0)(nil) })
}

func TestIntBuiltins(t *testing.T) {
	constant := constIntAppender(42)
	if got := string(constant(nil)); got != "42i" {
		t.Fatalf("const int: expected 42i, got %q", got)
	}
	if got := string(constant(nil)); got != "42i" {
		t.Fatalf("const int second call: expected 42i, got %q", got)
	}

	dec := decIntAppender(1)
	for _, want := range []string{"1i", "0i", "-1i"} {
		if got := string(dec(nil)); got != want {
			t.Fatalf("dec int: expected %q, got %q", want, got)
		}
	}

	zipf := zipfIntAppender(10)
	got := string(zipf(nil))
	if !strings.HasSuffix(got, "i") {
		t.Fatalf("zipf int should have i suffix: %q", got)
	}
	v := parseInt(got[:len(got)-1])
	if v < 0 || v >= 10 {
		t.Fatalf("zipf int out of range: %q", got)
	}
}

func TestFloatBuiltins(t *testing.T) {
	constant := constFloatAppender(42)
	if got := string(constant(nil)); got != "42" {
		t.Fatalf("const float: expected 42, got %q", got)
	}

	dec := decFloatAppender(1)
	for _, want := range []string{"1", "0", "-1"} {
		if got := string(dec(nil)); got != want {
			t.Fatalf("dec float: expected %q, got %q", want, got)
		}
	}

	randf := randDecimalFloatAppender(10)
	got := string(randf(nil))
	assertDecimalFloatToken(t, got)
	v := parseFloat64(t, got)
	if v < 0 || v >= 10 {
		t.Fatalf("randf out of range: %q", got)
	}

	for name, fn := range map[string]ValueAppender{
		"normal": normalFloatAppender(0),
		"sin":    sinFloatAppender(0),
		"walk":   walkFloatAppender(0),
	} {
		if got := string(fn(nil)); got != "0.0" {
			t.Fatalf("%s(0): expected 0.0, got %q", name, got)
		}
	}

	sin := sinFloatAppender(10)
	if got := string(sin(nil)); got != "0.0" {
		t.Fatalf("sin first value: expected 0.0, got %q", got)
	}
	got = string(sin(nil))
	assertDecimalFloatToken(t, got)
	if got == "0.0" {
		t.Fatal("sin second value should advance")
	}

	walk := walkFloatAppender(10)
	assertDecimalFloatToken(t, string(walk(nil)))

	zipf := zipfFloatAppender(10)
	got = string(zipf(nil))
	assertIntegerFloatToken(t, got)
	v = parseFloat64(t, got)
	if v < 0 || v >= 10 {
		t.Fatalf("zipf float out of range: %q", got)
	}
}

func TestStrBuiltins(t *testing.T) {
	inc := incStrAppender(7)
	for _, want := range []string{"7", "8", "9"} {
		if got := string(inc(nil)); got != want {
			t.Fatalf("str inc: expected %q, got %q", want, got)
		}
	}

	id := idStrAppender(4)
	for _, want := range []string{"id-0000", "id-0001"} {
		if got := string(id(nil)); got != want {
			t.Fatalf("str id: expected %q, got %q", want, got)
		}
	}
	if got := string(idStrAppender(0)(nil)); got != "id-0" {
		t.Fatalf("str id width 0: expected id-0, got %q", got)
	}

	hash := hashStrAppender(7)
	got := string(hash(nil))
	if len(got) != 6 {
		t.Fatalf("hash(7) length: expected 6, got %d (%q)", len(got), got)
	}
	assertSafeHexString(t, got)
	next := string(hash(nil))
	if next == got {
		t.Fatalf("hash should advance by counter: repeated %q", got)
	}

	longHash := string(hashStrAppender(80)(nil))
	if len(longHash) != 80 {
		t.Fatalf("hash(80) length: expected 80, got %d", len(longHash))
	}
	assertSafeHexString(t, longHash)

	if got := string(hashStrAppender(0)(nil)); got != "" {
		t.Fatalf("hash(0): expected empty string, got %q", got)
	}

	for _, got := range []string{string(incStrAppender(0)(nil)), string(idStrAppender(2)(nil)), longHash} {
		if strings.ContainsAny(got, ", =\"") {
			t.Fatalf("string builtin emitted line-protocol metacharacter: %q", got)
		}
	}
}

func TestStatefulBuiltinsWrappedByCount(t *testing.T) {
	tests := []struct {
		name     string
		fn       *Function
		cycled   []string
		repeated []string
	}{
		{
			name:     "int dec",
			fn:       &Function{Type: "int", Fn: "dec", Argument: 3, Count: 2},
			cycled:   []string{"3i", "2i", "3i"},
			repeated: []string{"3i", "3i", "2i"},
		},
		{
			name:     "float dec",
			fn:       &Function{Type: "float", Fn: "dec", Argument: 3, Count: 2},
			cycled:   []string{"3", "2", "3"},
			repeated: []string{"3", "3", "2"},
		},
		{
			name:     "str inc",
			fn:       &Function{Type: "str", Fn: "inc", Argument: 3, Count: 2},
			cycled:   []string{"3", "4", "3"},
			repeated: []string{"3", "3", "4"},
		},
		{
			name:     "str id",
			fn:       &Function{Type: "str", Fn: "id", Argument: 2, Count: 2},
			cycled:   []string{"id-00", "id-01", "id-00"},
			repeated: []string{"id-00", "id-00", "id-01"},
		},
	}

	for _, tt := range tests {
		cycled := tt.fn.NewAppender(2)
		for _, want := range tt.cycled {
			if got := string(cycled(nil)); got != want {
				t.Fatalf("%s cycle: expected %q, got %q", tt.name, want, got)
			}
		}

		repeatedFn := *tt.fn
		repeatedFn.Count = 0
		repeated := repeatedFn.NewAppender(2)
		for _, want := range tt.repeated {
			if got := string(repeated(nil)); got != want {
				t.Fatalf("%s nTimes: expected %q, got %q", tt.name, want, got)
			}
		}
	}
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

func assertDecimalFloatToken(t *testing.T, s string) {
	t.Helper()
	if s == "" {
		t.Fatal("decimal float token is empty")
	}
	if !strings.Contains(s, ".") {
		t.Fatalf("decimal float token should contain '.': %q", s)
	}
	if strings.Contains(s, "i") {
		t.Fatalf("decimal float token should not contain 'i': %q", s)
	}
	if _, err := strconv.ParseFloat(s, 64); err != nil {
		t.Fatalf("decimal float token should parse as float: %q: %v", s, err)
	}
}

func assertSafeHexString(t *testing.T, s string) {
	t.Helper()
	for _, ch := range s {
		if !strings.ContainsRune("0123456789abcdef", ch) {
			t.Fatalf("expected hex string, got %q", s)
		}
	}
	if strings.ContainsAny(s, ", =\"") {
		t.Fatalf("hex string emitted line-protocol metacharacter: %q", s)
	}
}

func parseFloat64(t *testing.T, s string) float64 {
	t.Helper()
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		t.Fatalf("failed parsing float %q: %v", s, err)
	}
	return v
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
