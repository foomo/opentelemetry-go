---
title: Custom Semconv
description: Define and use custom semantic conventions
---

# Custom Semantic Conventions

Examples showing how to use the `semconv` and `gotsrpcconv` packages.

## Attaching Attributes to Spans

Use the `semconv` constructors to add typed attributes to spans:

```go
package main

import (
	"context"

	"github.com/foomo/opentelemetry-go/semconv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func handleRequest(ctx context.Context, trackingID string) {
	tracer := otel.Tracer("my-service")

	ctx, span := tracer.Start(ctx, "gotsrpc.call",
		trace.WithAttributes(
			semconv.GoTSRPCService("UserService"),
			semconv.GoTSRPCFunc("GetProfile"),
			semconv.GoTSRPCPackage("github.com/foomo/myapp"),
			semconv.TrackingID(trackingID),
		),
	)
	defer span.End()

	// ... handler logic ...
}
```

## Using Keel Service Attributes

```go
ctx, span := tracer.Start(ctx, "keel.init",
	trace.WithAttributes(
		semconv.KeelServiceType("http"),
		semconv.KeelServiceName("api-gateway"),
		semconv.KeelServiceInst(0),
	),
)
```

## Recording GoTSRPC Execution Duration

The `gotsrpcconv.ExecutionDuration` wraps a `Float64Histogram` with pre-configured metadata:

```go
package main

import (
	"context"
	"time"

	"github.com/foomo/opentelemetry-go/semconv/gotsrpcconv"
	"go.opentelemetry.io/otel"
)

func main() {
	meter := otel.Meter("my-service")

	duration, err := gotsrpcconv.NewExecutionDuration(meter)
	if err != nil {
		panic(err)
	}

	// Record a GoTSRPC execution duration
	start := time.Now()
	// ... execute RPC ...
	elapsed := time.Since(start).Seconds()

	duration.Record(
		context.Background(),
		elapsed,
		"github.com/foomo/myapp",  // package
		"UserService",             // service
		"GetProfile",              // function
	)
}
```

The `Record` method automatically adds `gotsrpc.package`, `gotsrpc.service`, and `gotsrpc.func` attributes.

## Recording with Error Tracking

Use `AttrError` to tag executions that resulted in errors:

```go
duration.Record(
	ctx,
	elapsed,
	"github.com/foomo/myapp",
	"UserService",
	"GetProfile",
	duration.AttrError(true),
)
```

## Defining New Conventions

Follow the same pattern used in `semconv/` to define your own:

```go
package mysemconv

import "go.opentelemetry.io/otel/attribute"

const (
	TenantIDKey = attribute.Key("tenant.id")
	RegionKey   = attribute.Key("region")
)

func TenantID(v string) attribute.KeyValue {
	return TenantIDKey.String(v)
}

func Region(v string) attribute.KeyValue {
	return RegionKey.String(v)
}
```

Then use them like any other attribute constructor:

```go
span.SetAttributes(
	mysemconv.TenantID("acme-corp"),
	mysemconv.Region("eu-west-1"),
)
```
