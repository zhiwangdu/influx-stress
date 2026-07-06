package workload

import (
	"strconv"
	"strings"
)

type PointRenderer struct {
	literals []string
	values   []ValueAppender
	time     func() int64
}

func NewPointRenderer(template string, templates Templates, seriesCount int, time func() int64) *PointRenderer {
	return &PointRenderer{
		literals: strings.Split(template, "%v"),
		values:   templates.InitAppenders(seriesCount),
		time:     time,
	}
}

func (r *PointRenderer) AppendPoint(dst []byte) []byte {
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

func (r *PointRenderer) EstimatedPointSize() int {
	n := 0
	for _, lit := range r.literals {
		n += len(lit)
	}
	return n + (len(r.values)+1)*20
}

func NewBatchBuffer(renderer *PointRenderer, batchSize int) []byte {
	if renderer == nil || batchSize <= 0 {
		return []byte{}
	}
	return make([]byte, 0, (renderer.EstimatedPointSize()+1)*batchSize)
}
