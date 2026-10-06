---
title: semconv
description: API reference for semantic convention attribute constructors
---

# semconv

```go
import "github.com/foomo/opentelemetry-go/semconv"
```

Custom semantic convention attribute keys and constructors. Part of the root module.

## Error Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `ErrorType(err error)` | `error.type` | `error` | `attribute.KeyValue` |

Extends upstream `semconv.ErrorType`: maps `context.Canceled` and `context.DeadlineExceeded` (also when wrapped) to the stable values `context.Canceled` and `context.DeadlineExceeded` instead of their reflected type names.

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

::: warning Deprecated
`http.request.*` is an upstream namespace. Use `semconv.HTTPRequestHeader("x-request-id", v)` and `semconv.HTTPRequestHeader("referer", v)` from `go.opentelemetry.io/otel/semconv` instead.
:::

## Keel Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `KeelServiceType(v string)` | `keel.service.type` | `string` | `attribute.KeyValue` |
| `KeelServiceName(v string)` | `keel.service.name` | `string` | `attribute.KeyValue` |
| `KeelServiceInst(v int)` | `keel.service.inst` | `int` | `attribute.KeyValue` |

## Circuit Breaker Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `CircuitBreakerName(v string)` | `circuit_breaker.name` | `string` | `attribute.KeyValue` |
| `CircuitBreakerState(v string)` | `circuit_breaker.state` | `string` | `attribute.KeyValue` |

## Distributed Lock Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `DistributedLockName(v string)` | `distributed_lock.name` | `string` | `attribute.KeyValue` |
| `DistributedLockSystem(v string)` | `distributed_lock.system` | `string` | `attribute.KeyValue` |
| `DistributedLockAcquireResult(v string)` | `distributed_lock.acquire.result` | `string` | `attribute.KeyValue` |
| `DistributedLockOwnerID(v string)` | `distributed_lock.owner.id` | `string` | `attribute.KeyValue` |

`distributed_lock.name` must be the logical lock name, not the per-resource key. `distributed_lock.owner.id` values are unbounded: use on spans and logs only, never on metrics.

## NATS Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `NATSClientName(v string)` | `nats.client.name` | `string` | `attribute.KeyValue` |
| `NATSClientErrorKind(v string)` | `nats.client.error.kind` | `string` | `attribute.KeyValue` |
| `MessagingNATSStream(v string)` | `messaging.nats.stream` | `string` | `attribute.KeyValue` |

`MessagingSystemNats` is a predefined `messaging.system="nats"` attribute.

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
| `SamplingPriority(v int)` | `sampling.priority` | `int` | `attribute.KeyValue` |

`sampling.priority > 0` asks samplers to keep the trace. Match it in the collector with a `tail_sampling` `numeric_attribute` policy (`min_value: 1`).

## Tracking Attributes

| Function | Key | Parameter | Return |
|---|---|---|---|
| `TrackingID(v string)` | `tracking.id` | `string` | `attribute.KeyValue` |

`tracking.id` values are unbounded: use on spans and logs only, never on metrics.

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

---

## natsconv

```go
import "github.com/foomo/opentelemetry-go/semconv/natsconv"
```

Typed metric instrument wrappers for NATS conventions in the `nats` and `messaging.nats` namespaces. Every constructor `New<Instrument>(m metric.Meter, opt ...)` returns a no-op instrument if the meter is `nil`. Each instrument exposes `Inst()`, `Name()`, `Unit()` and `Description()`.

### Counters

| Instrument | Name | Unit | Record |
|---|---|---|---|
| `ClientDisconnects` | `nats.client.disconnects` | `{event}` | `Add(ctx, incr, serverAddress, attrs...)` |
| `ClientReconnects` | `nats.client.reconnects` | `{event}` | `Add(ctx, incr, serverAddress, attrs...)` |
| `ClientAsyncErrors` | `nats.client.async_errors` | `{error}` | `Add(ctx, incr, kind AsyncErrorKindAttr, attrs...)` |

All counters also provide `AddSet(ctx, incr, set attribute.Set)`.

Optional attribute helpers:

- `ClientDisconnects`, `ClientReconnects`: `AttrServerPort(int)`, `AttrClientName(string)`
- `ClientAsyncErrors`: `AttrSubject(string)` (`messaging.destination.name`), `AttrServerAddress(string)`

### Observable Gauges

| Instrument | Name | Unit | Record |
|---|---|---|---|
| `JetStreamConsumerPending` | `nats.jetstream.consumer.pending` | `{message}` | `Observe(o, val, stream, consumerGroupName, attrs...)` |
| `JetStreamConsumerAckPending` | `nats.jetstream.consumer.ack_pending` | `{message}` | `Observe(o, val, stream, consumerGroupName, attrs...)` |
| `JetStreamConsumerRedelivered` | `nats.jetstream.consumer.redelivered` | `{message}` | `Observe(o, val, stream, consumerGroupName, attrs...)` |

All gauges also provide `ObserveSet(o, val, set attribute.Set)`. Call `Observe` from a callback registered with the meter. `stream` maps to `messaging.nats.stream`, `consumerGroupName` to `messaging.consumer.group.name`.

### Enums

- `AsyncErrorKindAttr` (`nats.client.error.kind`): `AsyncErrorKindSlowConsumer`, `AsyncErrorKindPermissionViolation`, `AsyncErrorKindAuthExpired`, `AsyncErrorKindAuthRevoked`, `AsyncErrorKindOther` (`_OTHER`)
- `ConnectionStatusAttr` (`nats.client.connection.status`): `ConnectionStatusConnected`, `ConnectionStatusDisconnected`, `ConnectionStatusReconnecting`, `ConnectionStatusConnecting`, `ConnectionStatusDraining`, `ConnectionStatusClosed`

---

## circuitbreakerconv

```go
import "github.com/foomo/opentelemetry-go/semconv/circuitbreakerconv"
```

Library-agnostic metric instrument wrappers for circuit breakers in the `circuit_breaker` namespace. Every constructor `New<Instrument>(m metric.Meter, opt ...)` returns a no-op instrument if the meter is `nil`. Each instrument exposes `Inst()`, `Name()`, `Unit()` and `Description()`.

### Counters

| Instrument | Name | Unit | Record |
|---|---|---|---|
| `StateChanges` | `circuit_breaker.state_changes` | `{change}` | `Add(ctx, incr, name, state StateAttr, attrs...)` |
| `Rejections` | `circuit_breaker.rejections` | `{call}` | `Add(ctx, incr, name, attrs...)` |

All counters also provide `AddSet(ctx, incr, set attribute.Set)`. `name` maps to `circuit_breaker.name`; for `StateChanges`, `state` is the state transitioned to.

Optional attribute helpers:

- `Rejections`: `AttrState(StateAttr)`

### Enums

- `StateAttr` (`circuit_breaker.state`): `StateClosed`, `StateOpen`, `StateHalfOpen`

---

## distributedlockconv

```go
import "github.com/foomo/opentelemetry-go/semconv/distributedlockconv"
```

Backend-agnostic metric instrument wrappers for distributed locks in the `distributed_lock` namespace. Every constructor `New<Instrument>(m metric.Meter, opt ...)` returns a no-op instrument if the meter is `nil`. Each instrument exposes `Inst()`, `Name()`, `Unit()` and `Description()`.

### Histograms

| Instrument | Name | Unit | Record |
|---|---|---|---|
| `AcquireDuration` | `distributed_lock.acquire.duration` | `s` | `Record(ctx, val, name, system SystemAttr, result AcquireResultAttr, attrs...)` |
| `HoldDuration` | `distributed_lock.hold.duration` | `s` | `Record(ctx, val, name, system SystemAttr, attrs...)` |

All histograms also provide `RecordSet(ctx, val, set attribute.Set)`. `AcquireDuration` includes time spent waiting; its count doubles as the number of acquire attempts.

Optional attribute helpers:

- `AcquireDuration`: `AttrErrorType(error)` (`error.type`)

### Counters

| Instrument | Name | Unit | Record |
|---|---|---|---|
| `Lost` | `distributed_lock.lost` | `{lock}` | `Add(ctx, incr, name, system SystemAttr, attrs...)` |

Also provides `AddSet(ctx, incr, set attribute.Set)`. Counts locks lost while still held (lease expired or taken over).

### Enums

- `SystemAttr` (`distributed_lock.system`): `SystemRedis`, `SystemMongoDB`, `SystemNATS`, `SystemEtcd`, `SystemPostgreSQL`
- `AcquireResultAttr` (`distributed_lock.acquire.result`): `AcquireResultAcquired`, `AcquireResultContended`, `AcquireResultTimeout`, `AcquireResultError`
