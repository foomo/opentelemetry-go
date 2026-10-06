// Package testing provides helpers that wire an OpenTelemetry exporter into a
// fully configured tracer or meter provider for a single test or TestMain.
//
// Telemetry is collected while the test runs and exported once on cleanup, so
// the output appears after the test body. Pair the helpers with the test
// constructors of the glossy exporters:
//
//	func TestFoo(t *testing.T) {
//		tp := oteltesting.ReportTraces(t, glossytrace.NewTest(t))
//		mp := oteltesting.ReportMetrics(t, glossymetric.NewTest(t))
//		// use tp.Tracer(...) and mp.Meter(...)
//	}
package testing
