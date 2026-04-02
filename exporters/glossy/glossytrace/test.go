package glossytrace

import (
	"os"
	"testing"

	testingx "github.com/foomo/go/testing"
	"go.opentelemetry.io/otel/sdk/trace"
)

// NewTest creates a trace.SpanExporter for testing, using the provided testing.TB and optional configuration options.
func NewTest(tb testing.TB, opts ...Option) trace.SpanExporter {
	tb.Helper()

	exporter, err := New(append([]Option{WithWriter(tb.Output()), WithoutTimestamps()}, opts...)...)
	if err != nil {
		tb.Fatal(err)
	}

	return exporter
}

// NewTestMain creates a trace.SpanExporter for testing, using the provided testingx.M and optional configuration options.
func NewTestMain(m testingx.M, opts ...Option) trace.SpanExporter {
	exporter, err := New(append([]Option{WithWriter(os.Stdout), WithoutTimestamps()}, opts...)...)
	if err != nil {
		panic(err)
	}

	return exporter
}
