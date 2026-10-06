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

Each file defines exported `<Name>Key` `attribute.Key` constants paired with exported `<Name>` constructor functions that return `attribute.KeyValue`:

```go
// Key constant
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

### Error

| Constructor | Key | Type |
|---|---|---|
| `ErrorType(err)` | `error.type` | error |

Drop-in for upstream `semconv.ErrorType` that reports `context.Canceled` and `context.DeadlineExceeded` (also when wrapped) as stable values instead of reflected type names:

```go
err := fmt.Errorf("fetch user: %w", context.DeadlineExceeded)
span.SetAttributes(semconv.ErrorType(err)) // error.type="context.DeadlineExceeded"
```

### HTTP

| Constructor | Key | Type |
|---|---|---|
| `HTTPXRequestID(v)` | `http.request.id` | string |
| `HTTPXRequestReferer(v)` | `http.request.referer` | string |

::: warning Deprecated
Both live in the upstream `http.request.*` namespace. Use `semconv.HTTPRequestHeader("x-request-id", v)` / `semconv.HTTPRequestHeader("referer", v)` from `go.opentelemetry.io/otel/semconv` instead.
:::

### NATS

| Constructor | Key | Type |
|---|---|---|
| `NATSClientName(v)` | `nats.client.name` | string |
| `NATSClientErrorKind(v)` | `nats.client.error.kind` | string |
| `MessagingNATSStream(v)` | `messaging.nats.stream` | string |
| `MessagingSystemNats` (var) | `messaging.system` = `nats` | — |

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

`tracking.id` is unbounded — use on spans and logs only, never as a metric attribute.

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

## natsconv -- Metric Instruments

The `natsconv` sub-package provides NATS client and JetStream consumer instruments:

```go
import "github.com/foomo/opentelemetry-go/semconv/natsconv"
```

```go
disconnects, err := natsconv.NewClientDisconnects(meter)
if err != nil {
	return err
}
disconnects.Add(ctx, 1, nc.ConnectedUrl(),
	disconnects.AttrClientName("my-service"),
)
```

Counters: `ClientDisconnects`, `ClientReconnects`, `ClientAsyncErrors`. Observable gauges (observe from a registered callback): `JetStreamConsumerPending`, `JetStreamConsumerAckPending`, `JetStreamConsumerRedelivered`.

See the full API in the [semconv API reference](/api/semconv).
