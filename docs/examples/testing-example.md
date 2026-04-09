---
title: Testing Example
description: Full test files using the test helpers
---

# Testing Example

Complete examples showing how to integrate OpenTelemetry into Go tests.

## Trace Testing

```go
package handler_test

import (
	"context"
	"testing"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace"
	oteltesting "github.com/foomo/opentelemetry-go/testing"
)

func TestHandler(t *testing.T) {
	// Create a test exporter with flamegraph and attribute output
	exporter := glossytrace.NewTest(t,
		glossytrace.WithFlamegraph(),
		glossytrace.WithSpanAttributes(),
	)

	// Wire into a TracerProvider with automatic cleanup
	tp := oteltesting.ReportTraces(t, exporter)
	tracer := tp.Tracer("handler-test")

	// Your test code: create spans as usual
	ctx, span := tracer.Start(context.Background(), "TestHandler")
	defer span.End()

	_, child := tracer.Start(ctx, "db.Query")
	// ... exercise handler ...
	child.End()

	// When the test ends, t.Cleanup flushes all spans through the exporter.
	// The trace tree and flamegraph are printed to the test output.
}
```

## Metric Testing

```go
package service_test

import (
	"context"
	"testing"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric"
	oteltesting "github.com/foomo/opentelemetry-go/testing"
	"go.opentelemetry.io/otel/sdk/metric"
)

func TestMetrics(t *testing.T) {
	// Create a test exporter (writes to t.Output())
	exporter := glossymetric.NewTest(t)

	// Wire into a MeterProvider with automatic cleanup
	mp := oteltesting.ReportMetrics(t, exporter)
	meter := mp.Meter("service-test")

	// Record metrics
	counter, _ := meter.Int64Counter("requests")
	counter.Add(context.Background(), 5)

	histogram, _ := meter.Float64Histogram("latency",
		metric.WithUnit("ms"),
	)
	histogram.Record(context.Background(), 12.5)

	// When the test ends, t.Cleanup collects and exports all metrics
}
```

## TestMain Example

Use `TestMainReportMetrics` when you need a shared `MeterProvider` across all tests in a package:

```go
package integration_test

import (
	"os"
	"testing"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric"
	oteltesting "github.com/foomo/opentelemetry-go/testing"
	"go.opentelemetry.io/otel"
)

var meter = otel.Meter("integration")

func TestMain(m *testing.M) {
	exporter := glossymetric.NewTestMain(m)
	mp, flush := oteltesting.TestMainReportMetrics(m, exporter)

	// Set global meter provider so all tests in this package use it
	otel.SetMeterProvider(mp)

	code := m.Run()
	flush() // Collect, export, and shutdown
	os.Exit(code)
}

func TestFeatureA(t *testing.T) {
	counter, _ := meter.Int64Counter("feature_a.calls")
	counter.Add(t.Context(), 1)
}

func TestFeatureB(t *testing.T) {
	counter, _ := meter.Int64Counter("feature_b.calls")
	counter.Add(t.Context(), 1)
}
```

## Combined Traces and Metrics

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
	// Both traces and metrics are printed when the test completes
}
```
