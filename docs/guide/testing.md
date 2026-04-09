---
title: Testing
description: Integrate OpenTelemetry traces and metrics into your Go tests
---

# Testing

The `oteltesting` package provides helpers that wire an exporter into a fully configured OTel provider for a single test or `TestMain`, handling cleanup automatically.

```go
import oteltesting "github.com/foomo/opentelemetry-go/testing"
```

## Quick Start

::: code-group

```go [testing.T -- Traces]
func TestHandler(t *testing.T) {
	exporter := glossytrace.NewTest(t,
		glossytrace.WithFlamegraph(),
		glossytrace.WithSpanAttributes(),
	)
	tp := oteltesting.ReportTraces(t, exporter)

	tracer := tp.Tracer("test")
	_, span := tracer.Start(context.Background(), "test-op")
	span.End()
	// Traces are printed when the test completes via t.Cleanup
}
```

```go [testing.T -- Metrics]
func TestMetrics(t *testing.T) {
	exporter := glossymetric.NewTest(t)
	mp := oteltesting.ReportMetrics(t, exporter)

	meter := mp.Meter("test")
	counter, _ := meter.Int64Counter("requests")
	counter.Add(context.Background(), 1)
	// Metrics are collected and printed when the test completes via t.Cleanup
}
```

```go [TestMain]
func TestMain(m *testing.M) {
	exporter := glossymetric.NewTestMain(m)
	mp, flush := oteltesting.TestMainReportMetrics(m, exporter)

	meter := mp.Meter("test")
	// ... set up global meter ...

	code := m.Run()
	flush() // Collect and export metrics, then shutdown
	os.Exit(code)
}
```

:::

## How It Works

### ReportTraces

`ReportTraces` creates a `TracerProvider` with a `BatchSpanProcessor` and registers a `tb.Cleanup` handler that calls `tp.Shutdown`, which flushes all buffered spans through the exporter.

::: tip
The batch processor is configured with a 1-hour timeout, so it never auto-flushes during the test. All spans are exported at once during cleanup, giving you a complete trace tree.
:::

### ReportMetrics

`ReportMetrics` creates a `MeterProvider` backed by a `ManualReader`. On `tb.Cleanup`, it calls `reader.Collect` to gather all recorded metrics, exports them through the exporter, then shuts down the provider.

### TestMainReportMetrics

`TestMainReportMetrics` returns both a `*metric.MeterProvider` and a `flush` function. Call `flush()` after `m.Run()` to collect metrics, export them, and shut down the provider. This is necessary because `TestMain` doesn't have `tb.Cleanup`.

## Combining Traces and Metrics

You can use both in the same test:

```go
func TestFullObservability(t *testing.T) {
	traceExp := glossytrace.NewTest(t, glossytrace.WithFlamegraph())
	tp := oteltesting.ReportTraces(t, traceExp)

	metricExp := glossymetric.NewTest(t)
	mp := oteltesting.ReportMetrics(t, metricExp)

	tracer := tp.Tracer("test")
	meter := mp.Meter("test")

	ctx, span := tracer.Start(context.Background(), "handle-request")
	counter, _ := meter.Int64Counter("requests")
	counter.Add(ctx, 1)
	span.End()
}
```

See the [testing API reference](/api/testing) for full function signatures.
