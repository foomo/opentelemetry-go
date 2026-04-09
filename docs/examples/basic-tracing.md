---
title: Basic Tracing
description: Instrument an HTTP handler with spans and export to terminal
---

# Basic Tracing

This example creates a glossy trace exporter, wires it into a `TracerProvider`, and produces a trace with parent-child spans.

## Tree View

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func main() {
	ctx := context.Background()

	exporter, err := glossytrace.New(
		glossytrace.WithSpanAttributes(),
	)
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("my-api"),
			semconv.ServiceVersion("1.0.0"),
		)),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	tracer := otel.Tracer("example")

	ctx, parent := tracer.Start(ctx, "HTTP GET /api/users")
	parent.SetAttributes(attribute.String("http.method", "GET"))

	_, child := tracer.Start(ctx, "db.Query")
	time.Sleep(5 * time.Millisecond) // simulate work
	child.End()

	parent.End()
}
```

Output:

```
=== TRACE abc123... ===
└─ HTTP GET /api/users (5.23 ms) [2025-01-15 10:30:00]
    http.method=GET
  └─ db.Query (5.10 ms) [2025-01-15 10:30:00]
```

## With Flamegraph

Add `WithFlamegraph()` to get a proportional bar chart after the tree:

```go
exporter, err := glossytrace.New(
	glossytrace.WithFlamegraph(),
	glossytrace.WithDurationThresholds(5*time.Millisecond, 50*time.Millisecond),
)
```

Output:

```
=== TRACE 0102030405060708090a0b0c0d0e0f10 ===
└─ HTTP GET /api/users (150.00 ms)
  └─ db.Query (40.00 ms)

Nested Flamegraph:
HTTP GET /api/users ██████████████████████████████ 150.00 ms
  db.Query        ████████ 40.00 ms
```

## Filtering Short Spans

Use `WithMinDuration` to focus on slow operations:

```go
exporter, err := glossytrace.New(
	glossytrace.WithMinDuration(time.Millisecond),
)
```

Spans shorter than 1ms are excluded from the output.
