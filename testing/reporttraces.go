package testing

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/trace"
)

// ReportTraces configures a trace provider for collecting and reporting spans, with cleanup on test completion.
//
//		func TestWithTrace(t *testing.T) {
//	    exporter := glossytrace.NewTest(t, glossytrace.WithFlamegraph(), glossytrace.WithSpanAttributes())
//		  tp := testingx.ReportTraces(t, exporter)
//		}
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
		ctx, cancel := context.WithTimeout(context.WithoutCancel(tb.Context()), time.Second)
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
