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

### Circuit Breaker

| Constructor | Key | Type |
|---|---|---|
| `CircuitBreakerName(v)` | `circuit_breaker.name` | string |
| `CircuitBreakerState(v)` | `circuit_breaker.state` | string |

### Distributed Lock

| Constructor | Key | Type |
|---|---|---|
| `DistributedLockName(v)` | `distributed_lock.name` | string |
| `DistributedLockSystem(v)` | `distributed_lock.system` | string |
| `DistributedLockAcquireResult(v)` | `distributed_lock.acquire.result` | string |
| `DistributedLockOwnerID(v)` | `distributed_lock.owner.id` | string |

Use the logical lock name (`order-sync`), not the per-resource key (`lock:order:123`). `distributed_lock.owner.id` is unbounded — spans and logs only.

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
| `SamplingPriority(v)` | `sampling.priority` | int |
| `TrackingID(v)` | `tracking.id` | string |

`tracking.id` is unbounded — use on spans and logs only, never as a metric attribute.

#### Forcing a trace to be kept

Set `sampling.priority` on any span of the trace to ask a tail-based sampler to always keep it:

```go
span.SetAttributes(semconv.SamplingPriority(1))
```

Match it in the OpenTelemetry Collector's `tail_sampling` processor, ahead of the probabilistic policies:

```yaml
processors:
  tail_sampling:
    policies:
      - name: force-keep
        type: numeric_attribute
        numeric_attribute:
          key: sampling.priority
          min_value: 1
          max_value: 100
      - name: errors
        type: status_code
        status_code: { status_codes: [ERROR] }
      - name: baseline
        type: probabilistic
        probabilistic: { sampling_percentage: 10 }
```

::: tip
The tail sampler only sees spans that reach it: keep SDK head sampling at always-on (or parent-based always-on). With multiple collectors, front them with the `loadbalancing` exporter routed by `traceID` so all spans of a trace land on the same instance.
:::

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

## circuitbreakerconv -- Metric Instruments

The `circuitbreakerconv` sub-package provides library-agnostic circuit breaker instruments. Wire them into your breaker's state-change and rejection hooks:

```go
import "github.com/foomo/opentelemetry-go/semconv/circuitbreakerconv"
```

```go
stateChanges, err := circuitbreakerconv.NewStateChanges(meter)
if err != nil {
	return err
}
stateChanges.Add(ctx, 1, "payment-api", circuitbreakerconv.StateOpen)
```

Counters: `StateChanges` (`circuit_breaker.state_changes`), `Rejections` (`circuit_breaker.rejections`).

## distributedlockconv -- Metric Instruments

The `distributedlockconv` sub-package provides backend-agnostic distributed lock instruments (Redis, MongoDB, NATS, etcd, PostgreSQL, ...):

```go
import "github.com/foomo/opentelemetry-go/semconv/distributedlockconv"
```

```go
acquire, err := distributedlockconv.NewAcquireDuration(meter)
if err != nil {
	return err
}
start := time.Now()
// ... acquire lock ...
acquire.Record(ctx, time.Since(start).Seconds(), "order-sync",
	distributedlockconv.SystemRedis,
	distributedlockconv.AcquireResultAcquired,
)
```

Histograms: `AcquireDuration` (`distributed_lock.acquire.duration`), `HoldDuration` (`distributed_lock.hold.duration`). Counter: `Lost` (`distributed_lock.lost`).

For tracing, wrap acquire in a `CLIENT` span named `distributed_lock.acquire <name>`; record hold time via `HoldDuration`, not a span.

See the full API in the [semconv API reference](/api/semconv).
