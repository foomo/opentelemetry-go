package testing

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
)

// TestMainReportMetrics sets up a MeterProvider with a metric exporter for testing and returns it along with a flush function.
//
//	func TestMain(m *testing.M) {
//		mp, flush := oteltesting.TestMainReportMetrics(m, glossymetric.NewTestingM(m))
//	}
func TestMainReportMetrics(m *testing.M, exporter metric.Exporter) (*metric.MeterProvider, func()) {
	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(0))
	mp := metric.NewMeterProvider(metric.WithReader(reader))

	otel.SetMeterProvider(mp)

	return mp, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := mp.Shutdown(ctx); err != nil {
			panic(err)
		}
	}
}

func ReportMetrics(tb testing.TB, exporter metric.Exporter) (*metric.MeterProvider, func()) {
	tb.Helper()

	reader := metric.NewPeriodicReader(exporter, metric.WithInterval(0))
	mp := metric.NewMeterProvider(metric.WithReader(reader))

	prev := otel.GetMeterProvider()

	otel.SetMeterProvider(mp)

	// tb.Cleanup(func() {
	// 	ctx, cancel := context.WithTimeout(context.WithoutCancel(tb.Context()), time.Second)
	// 	defer cancel()
	// 	if err := mp.Shutdown(ctx); err != nil {
	// 		tb.Fatal(err)
	// 	}
	// 	otel.SetMeterProvider(prev)
	// })

	return mp, func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := mp.Shutdown(ctx); err != nil {
			tb.Fatal(err)
		}

		otel.SetMeterProvider(prev)
	}
}
