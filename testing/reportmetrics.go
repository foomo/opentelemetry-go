package testing

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// TestMainReportMetrics returns a MeterProvider backed by a manual reader and
// a flush function for use in TestMain. Calling flush collects all metrics,
// writes them to exporter prefixed with the caller's file and line, and shuts
// down the provider. flush panics on any collect, export or shutdown error.
//
//	func TestMain(m *testing.M) {
//		mp, flush := oteltesting.TestMainReportMetrics(m, glossymetric.NewTestMain(m))
//		otel.SetMeterProvider(mp)
//		code := m.Run()
//		flush()
//		os.Exit(code)
//	}
func TestMainReportMetrics(m *testing.M, exporter metric.Exporter) (*metric.MeterProvider, func()) {
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))

	var pcs [1]uintptr

	n := runtime.Callers(2, pcs[:])

	return mp, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var rm metricdata.ResourceMetrics
		if err := reader.Collect(ctx, &rm); err != nil {
			panic(err)
		}

		if n > 0 {
			f, _ := runtime.CallersFrames([]uintptr{pcs[0]}).Next()
			_, _ = fmt.Fprintf(os.Stdout, "%s:%d: ", filepath.Base(f.File), f.Line)
		}

		if err := exporter.Export(ctx, &rm); err != nil {
			panic(err)
		}

		if err := mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

// ReportMetrics returns a MeterProvider backed by a manual reader. When the
// test finishes, a tb.Cleanup hook collects all metrics, writes them to
// exporter prefixed with the caller's file and line, and shuts down the
// provider. Errors are reported with tb.Fatal.
//
//	func TestWithMetrics(t *testing.T) {
//		mp := oteltesting.ReportMetrics(t, glossymetric.NewTest(t))
//		counter, _ := mp.Meter("test").Int64Counter("requests")
//		counter.Add(t.Context(), 1)
//	}
func ReportMetrics(tb testing.TB, exporter metric.Exporter) *metric.MeterProvider {
	tb.Helper()

	var pcs [1]uintptr

	n := runtime.Callers(2, pcs[:])

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))

	tb.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		var rm metricdata.ResourceMetrics
		if err := reader.Collect(ctx, &rm); err != nil {
			tb.Fatal(err)
		}

		if n > 0 {
			f, _ := runtime.CallersFrames([]uintptr{pcs[0]}).Next()
			_, _ = fmt.Fprintf(tb.Output(), "%s:%d: ", filepath.Base(f.File), f.Line)
		}

		if err := exporter.Export(ctx, &rm); err != nil {
			tb.Fatal(err)
		}

		if err := mp.Shutdown(ctx); err != nil {
			tb.Fatal(err)
		}
	})

	return mp
}
