package testing

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

// ReportTraces configures a trace provider for collecting and reporting spans, with cleanup on test completion.
//
//		func TestWithTrace(t *testing.T) {
//	    exporter := glossytrace.NewTesting(glossytrace.WithFlamegraph(), glossytrace.WithSpanAttributes())
//		  testing.ReportTraces(t, exporter)
//		}
func ReportTraces(tb testing.TB, exporter trace.SpanExporter) {
	tb.Helper()

	sp := trace.NewSimpleSpanProcessor(exporter)
	tp := trace.NewTracerProvider(
		trace.WithSpanProcessor(sp),
	)
	prev := otel.GetTracerProvider()

	otel.SetTracerProvider(tp)

	tb.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(tb.Context()), time.Second)
		defer cancel()

		if err := tp.Shutdown(ctx); err != nil {
			tb.Fatal(err)
		}

		otel.SetTracerProvider(prev)
	})
}
