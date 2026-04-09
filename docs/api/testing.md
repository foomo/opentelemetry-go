---
title: oteltesting
description: API reference for the testing helpers
---

# oteltesting

```go
import oteltesting "github.com/foomo/opentelemetry-go/testing"
```

Test helpers for wiring OpenTelemetry providers with automatic cleanup. Part of the root module.

## Functions

### ReportTraces

```go
func ReportTraces(tb testing.TB, exporter trace.SpanExporter) *trace.TracerProvider
```

Creates a `TracerProvider` with a `BatchSpanProcessor` configured to never auto-flush (1-hour batch timeout). Registers a `tb.Cleanup` handler that calls `tp.Shutdown`, flushing all buffered spans through the exporter.

**Parameters:**
- `tb` -- the test instance (`*testing.T`, `*testing.B`, etc.)
- `exporter` -- a span exporter (e.g., from `glossytrace.NewTest`)

**Returns:** a configured `*trace.TracerProvider` ready to create tracers

```go
func TestMyHandler(t *testing.T) {
	exporter := glossytrace.NewTest(t, glossytrace.WithFlamegraph())
	tp := oteltesting.ReportTraces(t, exporter)

	tracer := tp.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "handler")
	// ... test code ...
	span.End()
}
```

### ReportMetrics

```go
func ReportMetrics(tb testing.TB, exporter metric.Exporter) *metric.MeterProvider
```

Creates a `MeterProvider` backed by a `ManualReader`. On `tb.Cleanup`, collects all metrics via `reader.Collect`, exports them through the exporter, and shuts down the provider.

**Parameters:**
- `tb` -- the test instance
- `exporter` -- a metric exporter (e.g., from `glossymetric.NewTest`)

**Returns:** a configured `*metric.MeterProvider` ready to create meters

```go
func TestMyMetrics(t *testing.T) {
	exporter := glossymetric.NewTest(t)
	mp := oteltesting.ReportMetrics(t, exporter)

	meter := mp.Meter("test")
	counter, _ := meter.Int64Counter("requests")
	counter.Add(context.Background(), 1)
}
```

### TestMainReportMetrics

```go
func TestMainReportMetrics(m *testing.M, exporter metric.Exporter) (*metric.MeterProvider, func())
```

Sets up a `MeterProvider` with a `ManualReader` for use in `TestMain`. Returns the provider and a flush function that must be called manually after `m.Run()`.

The flush function:
1. Calls `reader.Collect` to gather all recorded metrics
2. Exports them through the exporter
3. Calls `mp.Shutdown`

**Parameters:**
- `m` -- the test main instance
- `exporter` -- a metric exporter (e.g., from `glossymetric.NewTestMain`)

**Returns:** `(*metric.MeterProvider, func())`

```go
func TestMain(m *testing.M) {
	exporter := glossymetric.NewTestMain(m)
	mp, flush := oteltesting.TestMainReportMetrics(m, exporter)

	// Use mp to create meters for your test suite
	_ = mp

	code := m.Run()
	flush()
	os.Exit(code)
}
```
