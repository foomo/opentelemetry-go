// Package glossymetric provides an OpenTelemetry metric exporter that prints
// metrics as styled, human-readable blocks to a terminal or any [io.Writer].
//
// Each export prints every metric grouped by instrumentation scope, with its
// description, unit, attributes and value. Histograms are summarized as
// min, max, average and count. Colors are rendered with lipgloss and disabled
// when the NO_COLOR environment variable is set.
//
// # Usage
//
// Register the [Exporter] with a reader like any other
// [go.opentelemetry.io/otel/sdk/metric.Exporter]:
//
//	exporter, err := glossymetric.New(glossymetric.WithoutHistograms())
//	if err != nil {
//		return err
//	}
//	mp := sdkmetric.NewMeterProvider(
//		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
//	)
//
// In tests, use [NewTest] to write to the test output, or [NewTestMain] from
// TestMain.
package glossymetric
