---
title: Basic Metrics
description: Record counters and histograms, export with glossymetric
---

# Basic Metrics

This example creates a glossy metric exporter, wires it into a `MeterProvider`, and records various metric types.

## Setup and Export

```go
package main

import (
	"context"
	"fmt"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func main() {
	ctx := context.Background()

	exporter, err := glossymetric.New()
	if err != nil {
		panic(err)
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("my-api"),
	)

	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	meter := mp.Meter("example")

	// Counter
	requests, _ := meter.Int64Counter("http.requests",
		metric.WithDescription("Number of HTTP requests"),
		metric.WithUnit("1"),
	)
	requests.Add(ctx, 1, metric.WithAttributes(
		attribute.String("method", "GET"),
		attribute.String("path", "/api/users"),
	))

	// Histogram
	latency, _ := meter.Float64Histogram("http.latency",
		metric.WithDescription("Request latency"),
		metric.WithUnit("ms"),
	)
	latency.Record(ctx, 12.5, metric.WithAttributes(
		attribute.String("method", "GET"),
	))
}
```

## Example Output

```
METRICS

Scope: example
──────────────────────────────────────────────────
Metric: http.requests
Description: Number of HTTP requests
Unit: 1
  attrs: method=GET path=/api/users
Value: 1
──────────────────────────────────────────────────
Metric: http.latency
Description: Request latency
Unit: ms
  attrs: method=GET
Histogram: min=12.500 max=12.500 avg=12.500 count=1
──────────────────────────────────────────────────
```

## Without Histograms

Use `WithoutHistograms()` to suppress histogram detail output:

```go
exporter, err := glossymetric.New(
	glossymetric.WithoutHistograms(),
)
```

This shows only the metric name and value without the min/max/avg/count breakdown.
