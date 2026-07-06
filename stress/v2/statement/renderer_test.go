package statement

import "testing"

var benchmarkRendererBytes []byte

func TestPointRendererAppendPoint(t *testing.T) {
	renderer := newPointRenderer(
		"cpu,host=%v value=%v %v",
		Templates{
			{Tags: []string{"server-a", "server-b"}},
			{Function: &Function{Type: "int", Fn: "inc", Argument: 7, Count: 0}},
		},
		1,
		func() int64 { return 123 },
	)

	got := string(renderer.AppendPoint(nil))
	want := "cpu,host=server-a value=7i 123"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestPointRendererPreservesLiteralPercent(t *testing.T) {
	renderer := newPointRenderer(
		"cpu,host=x usage=100%,value=%v %v",
		Templates{
			{Function: &Function{Type: "float", Fn: "inc", Argument: 42, Count: 0}},
		},
		1,
		func() int64 { return 456 },
	)

	got := string(renderer.AppendPoint(nil))
	want := "cpu,host=x usage=100%,value=42 456"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestPointRendererTimestampFormat(t *testing.T) {
	renderer := newPointRenderer(
		"cpu value=%v %v",
		Templates{
			{Function: &Function{Type: "float", Fn: "inc", Argument: 1, Count: 0}},
		},
		1,
		func() int64 { return 999999999 },
	)

	got := string(renderer.AppendPoint(nil))
	want := "cpu value=1 999999999"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func BenchmarkInsertRenderPoint(b *testing.B) {
	renderer := benchmarkPointRenderer()
	dst := make([]byte, 0, renderer.estimatedPointSize())

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst = dst[:0]
		dst = renderer.AppendPoint(dst)
	}
	benchmarkRendererBytes = dst
}

func BenchmarkInsertRenderBatch(b *testing.B) {
	const batchSize = 5000
	renderer := benchmarkPointRenderer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst := newBatchBuffer(renderer, batchSize)
		for j := 0; j < batchSize; j++ {
			dst = renderer.AppendPoint(dst)
			dst = append(dst, '\n')
		}
		benchmarkRendererBytes = dst
	}
}

func benchmarkPointRenderer() *pointRenderer {
	var ts int64 = 1451606400
	return newPointRenderer(
		"cpu,host=%v,device=%v busy=%v,free=%v %v",
		Templates{
			{Tags: []string{"us-west", "us-east", "eu-north"}},
			{Function: &Function{Type: "str", Fn: "rand", Argument: 8, Count: 100}},
			{Function: &Function{Type: "int", Fn: "rand", Argument: 1000, Count: 0}},
			{Function: &Function{Type: "float", Fn: "inc", Argument: 0, Count: 0}},
		},
		300,
		func() int64 {
			ts += 10
			return ts
		},
	)
}
