---
title: Getting Started
description: Install opentelemetry-go and run your first trace
---

# Getting Started

## Prerequisites

- **Go 1.26+**
- Familiarity with the [OpenTelemetry Go SDK](https://opentelemetry.io/docs/languages/go/)

## Installation

Install the packages you need. The exporters are separate Go modules:

::: code-group

```sh [Trace exporter]
go get github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace
```

```sh [Metric exporter]
go get github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric
```

```sh [Semconv + testing helpers]
go get github.com/foomo/opentelemetry-go
```

:::

## Minimal Example

Create a glossy trace exporter, wire it into a `TracerProvider`, and produce a span:

```go
package main

import (
	"context"
	"fmt"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

func main() {
	ctx := context.Background()

	exporter, err := glossytrace.New()
	if err != nil {
		panic(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("my-service"),
			semconv.ServiceVersion("0.1.0"),
		)),
	)
	otel.SetTracerProvider(tp)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	tracer := otel.Tracer("example")
	_, span := tracer.Start(ctx, "hello-world")
	span.End()
}
```

Running this prints a styled trace tree to your terminal:

```
=== TRACE abc123... ===
└─ hello-world (0.12 ms)
```

::: tip
Set `NO_COLOR=1` to disable color output, useful for CI or piping to files.
:::

## Next Steps

- [Exporters Guide](/guide/exporters) -- learn about all exporter options including flamegraphs and duration thresholds
- [Testing Guide](/guide/testing) -- integrate traces and metrics into your Go tests
- [API Reference](/api/) -- full function signatures and option tables
