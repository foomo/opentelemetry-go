package testing

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
)

// ReportTraces returns a TracerProvider that buffers all ended spans. When the
// test finishes, a tb.Cleanup hook shuts down the provider, which exports the
// buffered spans to exporter. Shutdown errors are reported with tb.Fatal.
//
//	func TestWithTrace(t *testing.T) {
//		exporter := glossytrace.NewTest(t, glossytrace.WithFlamegraph(), glossytrace.WithSpanAttributes())
//		tp := oteltesting.ReportTraces(t, exporter)
//		_, span := tp.Tracer("test").Start(t.Context(), "op")
//		span.End()
//	}
func ReportTraces(tb testing.TB, exporter trace.SpanExporter) *trace.TracerProvider {
	tb.Helper()

	// var pcs [1]uintptr
	// n := runtime.Callers(2, pcs[:])

	sp := trace.NewBatchSpanProcessor(exporter,
		trace.WithBatchTimeout(time.Hour), // never auto-flush
	)
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(sp),
	)

	tb.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// if n > 0 {
		// 	f, _ := runtime.CallersFrames([]uintptr{pcs[0]}).Next()
		// 	_, _ = fmt.Fprintf(tb.Output(), "%s:%d: ", filepath.Base(f.File), f.Line)
		// }

		if err := tp.Shutdown(ctx); err != nil {
			tb.Fatal(err)
		}
	})

	return tp
}
