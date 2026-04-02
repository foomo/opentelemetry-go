package glossymetric

import (
	"os"
	"testing"

	testingx "github.com/foomo/go/testing"
	"go.opentelemetry.io/otel/sdk/metric"
)

// NewTest creates a metric exporter with the provided options and associates it with the testing.TB instance.
// It uses tb.Output() as the writer and calls tb.Fatal() if initialization fails.
func NewTest(tb testing.TB, opts ...Option) metric.Exporter {
	tb.Helper()

	exporter, err := New(append([]Option{WithWriter(tb.Output())}, opts...)...)
	if err != nil {
		tb.Fatal(err)
	}

	return exporter
}

// NewTestMain creates a new metric exporter configured for testing and outputs to os.Stdout by default.
func NewTestMain(m testingx.M, opts ...Option) metric.Exporter {
	exporter, err := New(append([]Option{WithWriter(os.Stdout)}, opts...)...)
	if err != nil {
		panic(err)
	}

	return exporter
}
