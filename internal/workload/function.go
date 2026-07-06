package workload

import (
	"crypto/sha256"
	"math"
	"math/rand"
	"strconv"
)

const (
	sinePeriod = 100
	zipfS      = 1.5
	zipfV      = 1
)

// Built-in function argument meanings:
// const=value, inc/dec=start, rand/randf/zipf=upper bound, normal=scale,
// sin=amplitude with a fixed period, walk=max step, id=counter width,
// hash=str-rand-compatible hex length.

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
	case "inc":
		return incStrAppender(arg)
	case "id":
		return idStrAppender(arg)
	case "hash":
		return hashStrAppender(arg)
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
	case "const":
		return constFloatAppender(arg)
	case "dec":
		return decFloatAppender(arg)
	case "rand":
		return randFloatAppender(arg)
	case "randf":
		return randDecimalFloatAppender(arg)
	case "inc":
		return incFloatAppender(arg)
	case "normal":
		return normalFloatAppender(arg)
	case "sin":
		return sinFloatAppender(arg)
	case "walk":
		return walkFloatAppender(arg)
	case "zipf":
		return zipfFloatAppender(arg)
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
	case "const":
		return constIntAppender(arg)
	case "dec":
		return decIntAppender(arg)
	case "rand":
		return randIntAppender(arg)
	case "inc":
		return incIntAppender(arg)
	case "zipf":
		return zipfIntAppender(arg)
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

func constIntAppender(n int) ValueAppender {
	return func(dst []byte) []byte {
		dst = strconv.AppendInt(dst, int64(n), 10)
		return append(dst, 'i')
	}
}

func decIntAppender(n int) ValueAppender {
	i := n
	return func(dst []byte) []byte {
		dst = strconv.AppendInt(dst, int64(i), 10)
		i--
		return append(dst, 'i')
	}
}

func constFloatAppender(n int) ValueAppender {
	return func(dst []byte) []byte {
		return strconv.AppendInt(dst, int64(n), 10)
	}
}

func decFloatAppender(n int) ValueAppender {
	i := n
	return func(dst []byte) []byte {
		dst = strconv.AppendInt(dst, int64(i), 10)
		i--
		return dst
	}
}

func randDecimalFloatAppender(n int) ValueAppender {
	return func(dst []byte) []byte {
		if n <= 0 {
			panic("randf requires positive bound")
		}
		return appendDecimalFloat(dst, rand.Float64()*float64(n))
	}
}

func normalFloatAppender(n int) ValueAppender {
	return func(dst []byte) []byte {
		if n == 0 {
			return appendDecimalFloat(dst, 0)
		}
		return appendDecimalFloat(dst, rand.NormFloat64()*float64(n))
	}
}

func sinFloatAppender(n int) ValueAppender {
	step := 0
	return func(dst []byte) []byte {
		v := math.Sin(2*math.Pi*float64(step)/sinePeriod) * float64(n)
		step++
		return appendDecimalFloat(dst, v)
	}
}

func walkFloatAppender(n int) ValueAppender {
	var v float64
	return func(dst []byte) []byte {
		if n != 0 {
			v += (rand.Float64()*2 - 1) * float64(n)
		}
		return appendDecimalFloat(dst, v)
	}
}

func zipfIntAppender(n int) ValueAppender {
	return zipfAppender(n, true)
}

func zipfFloatAppender(n int) ValueAppender {
	return zipfAppender(n, false)
}

func zipfAppender(n int, intSuffix bool) ValueAppender {
	var z *rand.Zipf
	return func(dst []byte) []byte {
		if n <= 0 {
			panic("zipf requires positive bound")
		}
		if z == nil {
			r := rand.New(rand.NewSource(rand.Int63()))
			z = rand.NewZipf(r, zipfS, zipfV, uint64(n-1))
		}
		dst = strconv.AppendUint(dst, z.Uint64(), 10)
		if intSuffix {
			dst = append(dst, 'i')
		}
		return dst
	}
}

func appendDecimalFloat(dst []byte, v float64) []byte {
	start := len(dst)
	dst = strconv.AppendFloat(dst, v, 'f', -1, 64)
	for _, b := range dst[start:] {
		if b == '.' {
			return dst
		}
	}
	return append(dst, '.', '0')
}

func incStrAppender(n int) ValueAppender {
	i := n
	return func(dst []byte) []byte {
		dst = strconv.AppendInt(dst, int64(i), 10)
		i++
		return dst
	}
}

func idStrAppender(width int) ValueAppender {
	i := 0
	return func(dst []byte) []byte {
		dst = append(dst, "id-"...)
		s := strconv.FormatInt(int64(i), 10)
		for pad := width - len(s); pad > 0; pad-- {
			dst = append(dst, '0')
		}
		dst = append(dst, s...)
		i++
		return dst
	}
}

func hashStrAppender(n int) ValueAppender {
	i := 0
	return func(dst []byte) []byte {
		const hex = "0123456789abcdef"
		need := n / 2
		for block := 0; need > 0; block++ {
			var scratch [40]byte
			b := strconv.AppendInt(scratch[:0], int64(i), 10)
			b = append(b, ':')
			b = strconv.AppendInt(b, int64(block), 10)
			sum := sha256.Sum256(b)
			for j := 0; j < len(sum) && need > 0; j++ {
				dst = append(dst, hex[sum[j]>>4], hex[sum[j]&0x0f])
				need--
			}
		}
		i++
		return dst
	}
}
