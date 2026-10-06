package glossytrace

import (
	"os"
	"testing"

	testingx "github.com/foomo/go/testing"
	"go.opentelemetry.io/otel/sdk/trace"
)

// NewTest returns an Exporter configured by opts that writes to tb.Output()
// without timestamps, so its output is attributed to the test and stable
// across runs. Options in opts take precedence. It calls tb.Fatal if the
// exporter cannot be created.
func NewTest(tb testing.TB, opts ...Option) trace.SpanExporter {
	tb.Helper()

	exporter, err := New(append([]Option{WithWriter(tb.Output()), WithoutTimestamps()}, opts...)...)
	if err != nil {
		tb.Fatal(err)
	}

	return exporter
}

// NewTestMain returns an Exporter configured by opts for use in TestMain. It
// writes to os.Stdout without timestamps unless opts overrides either. m is
// accepted to restrict the call to TestMain and is otherwise unused. It
// panics if the exporter cannot be created.
func NewTestMain(m testingx.M, opts ...Option) trace.SpanExporter {
	exporter, err := New(append([]Option{WithWriter(os.Stdout), WithoutTimestamps()}, opts...)...)
	if err != nil {
		panic(err)
	}

	return exporter
}
