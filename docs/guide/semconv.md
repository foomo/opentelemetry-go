---
title: Semantic Conventions
description: Custom attribute constructors and metric instruments
---

# Semantic Conventions

The `semconv` package provides custom semantic convention attribute keys and constructors, following the same pattern as the upstream [OpenTelemetry semantic conventions](https://pkg.go.dev/go.opentelemetry.io/otel/semconv).

```go
import "github.com/foomo/opentelemetry-go/semconv"
```

## Pattern

Each file defines unexported `attribute.Key` constants paired with exported constructor functions that return `attribute.KeyValue`:

```go
// Unexported key constant
const GoTSRPCFuncKey = attribute.Key("gotsrpc.func")

// Exported constructor
func GoTSRPCFunc(v string) attribute.KeyValue {
	return GoTSRPCFuncKey.String(v)
}
```

This pattern provides type safety (you can't accidentally pass an `int` to a string attribute) while keeping the key names consistent.

## Available Conventions

### GoTSRPC

Attributes for the [GoTSRPC](https://github.com/foomo/gotsrpc) framework:

| Constructor | Key | Type |
|---|---|---|
| `GoTSRPCFunc(v)` | `gotsrpc.func` | string |
| `GoTSRPCService(v)` | `gotsrpc.service` | string |
| `GoTSRPCPackage(v)` | `gotsrpc.package` | string |
| `GoTSRPCMarshalling(v)` | `gotsrpc.marshalling` | int64 |
| `GoTSRPCUnmarshalling(v)` | `gotsrpc.unmarshalling` | int64 |
| `GoTSRPCPayload(v)` | `gotsrpc.payload` | string |
| `GoTSRPCErrorCode(v)` | `gotsrpc.error.code` | int |
| `GoTSRPCErrorMessage(v)` | `gotsrpc.error.message` | string |
| `GoTSRPCErrorType(v)` | `gotsrpc.error.type` | string |

### HTTP

| Constructor | Key | Type |
|---|---|---|
| `HTTPXRequestID(v)` | `http.request.id` | string |
| `HTTPXRequestReferer(v)` | `http.request.referer` | string |

### Keel

Attributes for the [Keel](https://github.com/foomo/keel) service framework:

| Constructor | Key | Type |
|---|---|---|
| `KeelServiceType(v)` | `keel.service.type` | string |
| `KeelServiceName(v)` | `keel.service.name` | string |
| `KeelServiceInst(v)` | `keel.service.inst` | int |

### Other

| Constructor | Key | Type |
|---|---|---|
| `ProfileName(v)` | `profile.name` | string |
| `ReflectType(v)` | `reflect.type` | any |
| `TraceID(v)` | `trace.id` | string |
| `SpanID(v)` | `span.id` | string |
| `TrackingID(v)` | `tracking.id` | string |

## Usage

Attach attributes to spans or use them as metric labels:

```go
ctx, span := tracer.Start(ctx, "gotsrpc.call",
	trace.WithAttributes(
		semconv.GoTSRPCService("UserService"),
		semconv.GoTSRPCFunc("GetProfile"),
		semconv.GoTSRPCPackage("github.com/foomo/myapp"),
	),
)
defer span.End()
```

## gotsrpcconv -- Metric Instruments

The `gotsrpcconv` sub-package provides typed metric instrument wrappers:

```go
import "github.com/foomo/opentelemetry-go/semconv/gotsrpcconv"
```

### ExecutionDuration

A `Float64Histogram` wrapper for recording GoTSRPC execution durations:

```go
meter := mp.Meter("my-service")
duration, err := gotsrpcconv.NewExecutionDuration(meter)
if err != nil {
	panic(err)
}

// Record a duration with GoTSRPC context
duration.Record(ctx, 0.150, "github.com/foomo/myapp", "UserService", "GetProfile")
```

The instrument is pre-configured with:
- **Name**: `gotsrpc.execution.duration`
- **Unit**: `s`
- **Description**: `Duration of GOTSRPC execution.`

See the full API in the [semconv API reference](/api/semconv).
