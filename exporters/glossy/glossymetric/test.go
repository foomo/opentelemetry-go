package glossymetric

import (
	"os"
	"testing"

	testingx "github.com/foomo/go/testing"
	"go.opentelemetry.io/otel/sdk/metric"
)

// NewTest returns an Exporter configured by opts that writes to tb.Output(),
// so its output is attributed to the test. Options in opts take precedence.
// It calls tb.Fatal if the exporter cannot be created.
func NewTest(tb testing.TB, opts ...Option) metric.Exporter {
	tb.Helper()

	exporter, err := New(append([]Option{WithWriter(tb.Output())}, opts...)...)
	if err != nil {
		tb.Fatal(err)
	}

	return exporter
}

// NewTestMain returns an Exporter configured by opts for use in TestMain. It
// writes to os.Stdout unless opts sets [WithWriter]. m is accepted to restrict
// the call to TestMain and is otherwise unused. It panics if the exporter
// cannot be created.
func NewTestMain(m testingx.M, opts ...Option) metric.Exporter {
	exporter, err := New(append([]Option{WithWriter(os.Stdout)}, opts...)...)
	if err != nil {
		panic(err)
	}

	return exporter
}
