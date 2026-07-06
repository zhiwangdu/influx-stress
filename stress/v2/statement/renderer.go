package statement

import (
	"strconv"
	"strings"
)

type pointRenderer struct {
	literals []string
	values   []ValueAppender
	time     func() int64
}

func newPointRenderer(template string, templates Templates, seriesCount int, time func() int64) *pointRenderer {
	return &pointRenderer{
		literals: strings.Split(template, "%v"),
		values:   templates.InitAppenders(seriesCount),
		time:     time,
	}
}

func (r *pointRenderer) AppendPoint(dst []byte) []byte {
	for i, lit := range r.literals {
		dst = append(dst, lit...)
		if i == len(r.literals)-1 {
			break
		}

		switch {
		case i < len(r.values):
			dst = r.values[i](dst)
		case i == len(r.values):
			dst = strconv.AppendInt(dst, r.time(), 10)
		default:
			dst = append(dst, "%!v(MISSING)"...)
		}
	}
	return dst
}

func (r *pointRenderer) estimatedPointSize() int {
	n := 0
	for _, lit := range r.literals {
		n += len(lit)
	}
	return n + (len(r.values)+1)*20
}
