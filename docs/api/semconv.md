---
title: semconv
description: API reference for semantic convention attribute constructors
---

# semconv

```go
import "github.com/foomo/opentelemetry-go/semconv"
```

Custom semantic convention attribute keys and constructors. Part of the root module.

## GoTSRPC Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `GoTSRPCFunc(v string)` | `gotsrpc.func` | `string` | `attribute.KeyValue` |
| `GoTSRPCService(v string)` | `gotsrpc.service` | `string` | `attribute.KeyValue` |
| `GoTSRPCPackage(v string)` | `gotsrpc.package` | `string` | `attribute.KeyValue` |
| `GoTSRPCMarshalling(v int64)` | `gotsrpc.marshalling` | `int64` | `attribute.KeyValue` |
| `GoTSRPCUnmarshalling(v int64)` | `gotsrpc.unmarshalling` | `int64` | `attribute.KeyValue` |
| `GoTSRPCPayload(v string)` | `gotsrpc.payload` | `string` | `attribute.KeyValue` |
| `GoTSRPCErrorCode(v int)` | `gotsrpc.error.code` | `int` | `attribute.KeyValue` |
| `GoTSRPCErrorMessage(v string)` | `gotsrpc.error.message` | `string` | `attribute.KeyValue` |
| `GoTSRPCErrorType(v string)` | `gotsrpc.error.type` | `string` | `attribute.KeyValue` |

## HTTP Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `HTTPXRequestID(v string)` | `http.request.id` | `string` | `attribute.KeyValue` |
| `HTTPXRequestReferer(v string)` | `http.request.referer` | `string` | `attribute.KeyValue` |

## Keel Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `KeelServiceType(v string)` | `keel.service.type` | `string` | `attribute.KeyValue` |
| `KeelServiceName(v string)` | `keel.service.name` | `string` | `attribute.KeyValue` |
| `KeelServiceInst(v int)` | `keel.service.inst` | `int` | `attribute.KeyValue` |

## Profile Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `ProfileName(v string)` | `profile.name` | `string` | `attribute.KeyValue` |

## Reflect Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `ReflectType(v any)` | `reflect.type` | `any` | `attribute.KeyValue` |

## Trace Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `TraceID(v string)` | `trace.id` | `string` | `attribute.KeyValue` |
| `SpanID(v string)` | `span.id` | `string` | `attribute.KeyValue` |

## Tracking Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `TrackingID(v string)` | `tracking.id` | `string` | `attribute.KeyValue` |

---

## gotsrpcconv

```go
import "github.com/foomo/opentelemetry-go/semconv/gotsrpcconv"
```

Typed metric instrument wrappers for GoTSRPC conventions.

### ExecutionDuration

```go
type ExecutionDuration struct {
	metric.Float64Histogram
}
```

A `Float64Histogram` wrapper pre-configured with name `gotsrpc.execution.duration`, unit `s`, and description `Duration of GOTSRPC execution.`

#### NewExecutionDuration

```go
func NewExecutionDuration(m metric.Meter, opt ...metric.Float64HistogramOption) (ExecutionDuration, error)
```

Creates a new `ExecutionDuration` instrument. Returns a no-op if the meter is `nil`.

#### Methods

```go
func (m ExecutionDuration) Inst() metric.Float64Histogram
```

Returns the underlying histogram instrument.

```go
func (ExecutionDuration) Name() string        // "gotsrpc.execution.duration"
func (ExecutionDuration) Unit() string        // "s"
func (ExecutionDuration) Description() string // "Duration of GOTSRPC execution."
```

```go
func (m ExecutionDuration) Record(
	ctx context.Context,
	val float64,
	pkg string,   // gotsrpc.package
	svs string,   // gotsrpc.service
	fnc string,   // gotsrpc.func
	attrs ...attribute.KeyValue,
)
```

Records a duration value with GoTSRPC context attributes automatically added.

```go
func (m ExecutionDuration) RecordSet(ctx context.Context, val float64, set attribute.Set)
```

Records a duration value with a pre-built attribute set.

```go
func (ExecutionDuration) AttrError(val bool) attribute.KeyValue
```

Returns `attribute.Bool("gotsprc.error", val)`.
