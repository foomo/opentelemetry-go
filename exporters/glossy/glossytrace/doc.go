// Package glossytrace provides an OpenTelemetry span exporter that prints
// traces as styled, human-readable trees to a terminal or any [io.Writer].
//
// Spans are grouped by trace ID and printed as a parent/child tree with
// durations color-coded against configurable thresholds. Optional output
// includes span attributes and events, a nested flamegraph, and span links.
// Colors are rendered with lipgloss and disabled when the NO_COLOR
// environment variable is set.
//
// # Usage
//
// Register the [Exporter] with a tracer provider like any other
// [go.opentelemetry.io/otel/sdk/trace.SpanExporter]:
//
//	exporter, err := glossytrace.New(glossytrace.WithFlamegraph())
//	if err != nil {
//		return err
//	}
//	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
//
// In tests, use [NewTest] to write to the test output, or [NewTestMain] from
// TestMain.
package glossytrace
