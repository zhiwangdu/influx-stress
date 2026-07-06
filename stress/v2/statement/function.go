package statement

import (
	"math/rand"
	"strconv"
)

// ################
// # Function     #
// ################

// Function is a struct that holds information for generating values in templated points
type Function struct {
	Type     string
	Fn       string
	Argument int
	Count    int
}

// ValueAppender appends the next generated value to dst.
type ValueAppender func(dst []byte) []byte

// NewStringer creates a new Stringer
func (f *Function) NewStringer(series int) Stringer {
	return stringerFromAppender(f.NewAppender(series))
}

// NewAppender creates a value generator for templated writes.
func (f *Function) NewAppender(series int) ValueAppender {
	var fn ValueAppender
	switch f.Type {
	case "int":
		fn = newIntAppender(f.Fn, f.Argument)
	case "float":
		fn = newFloatAppender(f.Fn, f.Argument)
	case "str":
		fn = newStrAppender(f.Fn, f.Argument)
	default:
		fn = errorAppender("STRINGER ERROR")
	}

	if f.Count != 0 {
		return cycleAppender(f.Count, fn)
	}

	return nTimesAppender(series, fn)
}

// ################
// # Stringers    #
// ################

// Stringers is a collection of Stringer
type Stringers []Stringer

// Eval returns an array of all the Stringer functions evaluated once
func (s Stringers) Eval(time func() int64) []interface{} {
	arr := make([]interface{}, len(s)+1)

	for i, st := range s {
		arr[i] = st()
	}

	arr[len(s)] = time()

	return arr
}

// Stringer is a function that returns a string
type Stringer func() string

func stringerFromAppender(fn ValueAppender) Stringer {
	return func() string {
		return string(fn(nil))
	}
}

func errorAppender(s string) ValueAppender {
	return func(dst []byte) []byte {
		return append(dst, s...)
	}
}

func randStrAppender(n int) ValueAppender {
	return func(dst []byte) []byte {
		const hex = "0123456789abcdef"
		for i := 0; i < n/2; i++ {
			b := rand.Intn(256)
			dst = append(dst, hex[b>>4], hex[b&0x0f])
		}
		return dst
	}
}

func randStr(n int) func() string {
	return stringerFromAppender(randStrAppender(n))
}

func newStrAppender(fn string, arg int) ValueAppender {
	switch fn {
	case "rand":
		return randStrAppender(arg)
	default:
		return errorAppender("STR ERROR")
	}
}

// NewStrFunc reates a new striger to create strings for templated writes
func NewStrFunc(fn string, arg int) Stringer {
	return stringerFromAppender(newStrAppender(fn, arg))
}

func randFloatAppender(n int) ValueAppender {
	return func(dst []byte) []byte {
		return strconv.AppendInt(dst, int64(rand.Intn(n)), 10)
	}
}

func randFloat(n int) func() string {
	return stringerFromAppender(randFloatAppender(n))
}

func incFloatAppender(n int) ValueAppender {
	i := n
	return func(dst []byte) []byte {
		dst = strconv.AppendInt(dst, int64(i), 10)
		i++
		return dst
	}
}

func incFloat(n int) func() string {
	return stringerFromAppender(incFloatAppender(n))
}

func newFloatAppender(fn string, arg int) ValueAppender {
	switch fn {
	case "rand":
		return randFloatAppender(arg)
	case "inc":
		return incFloatAppender(arg)
	default:
		return errorAppender("FLOAT ERROR")
	}
}

// NewFloatFunc reates a new striger to create float values for templated writes
func NewFloatFunc(fn string, arg int) Stringer {
	return stringerFromAppender(newFloatAppender(fn, arg))
}

func randIntAppender(n int) ValueAppender {
	return func(dst []byte) []byte {
		dst = strconv.AppendInt(dst, int64(rand.Intn(n)), 10)
		return append(dst, 'i')
	}
}

func randInt(n int) Stringer {
	return stringerFromAppender(randIntAppender(n))
}

func incIntAppender(n int) ValueAppender {
	i := n
	return func(dst []byte) []byte {
		dst = strconv.AppendInt(dst, int64(i), 10)
		i++
		return append(dst, 'i')
	}
}

func incInt(n int) Stringer {
	return stringerFromAppender(incIntAppender(n))
}

func newIntAppender(fn string, arg int) ValueAppender {
	switch fn {
	case "rand":
		return randIntAppender(arg)
	case "inc":
		return incIntAppender(arg)
	default:
		return errorAppender("INT ERROR")
	}
}

// NewIntFunc reates a new striger to create int values for templated writes
func NewIntFunc(fn string, arg int) Stringer {
	return stringerFromAppender(newIntAppender(fn, arg))
}

// nTimesAppender returns the previous return value of a function n-many times
// before calling the function again.
func nTimesAppender(n int, fn ValueAppender) ValueAppender {
	i := 0
	t := fn(nil)
	return func(dst []byte) []byte {
		i++
		if i > n {
			t = fn(t[:0])
			i = 1
		}
		return append(dst, t...)
	}
}

// cycleAppender cycles through a fixed list of generated values before repeating them.
func cycleAppender(n int, fn ValueAppender) ValueAppender {
	if n == 0 {
		return fn
	}
	i := 0
	cache := make([][]byte, n)
	for j := range cache {
		cache[j] = fn(nil)
	}

	return func(dst []byte) []byte {
		dst = append(dst, cache[i]...)
		i = (i + 1) % n
		return dst
	}
}
