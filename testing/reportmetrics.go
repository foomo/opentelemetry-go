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

// TestMainReportMetrics sets up a MeterProvider with a metric exporter for testing and returns it along with a flush function.
//
//		func TestMain(m *testing.M) {
//	   exporter := glossymetric.NewTestingM(m)
//			mp, flush := oteltesting.TestMainReportMetrics(m, exporter)
//		}
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

// ReportMetrics sets up a MeterProvider for metrics reporting and returns it along with a cleanup function.
//
//		func TestWithMetrics(t *testing.T) {
//	    exporter := glossymetric.NewTest(t)
//		  mp := oteltesting.ReportMetrics(t, exporter)
//		}
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
