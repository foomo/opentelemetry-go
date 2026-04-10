---
title: glossytrace
description: API reference for the glossy trace exporter
---

# glossytrace

```go
import "github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace"
```

Styled terminal exporter for OpenTelemetry traces. Renders spans as tree views with optional flamegraph visualization.

## Types

### Exporter

```go
type Exporter struct {
	// contains filtered or unexported fields
}
```

Implements `go.opentelemetry.io/otel/sdk/trace.SpanExporter`.

## Constructors

### New

```go
func New(options ...Option) (*Exporter, error)
```

Creates a new glossy trace exporter. Output defaults to `os.Stdout`.

### NewTest

```go
func NewTest(tb testing.TB, opts ...Option) trace.SpanExporter
```

Creates an exporter for use in tests. Automatically applies `WithWriter(tb.Output())` and `WithoutTimestamps()`. Calls `tb.Fatal()` on initialization error.

### NewTestMain

```go
func NewTestMain(m testingx.M, opts ...Option) trace.SpanExporter
```

Creates an exporter for use in `TestMain`. Automatically applies `WithWriter(os.Stdout)` and `WithoutTimestamps()`. Panics on initialization error.

::: tip
`testingx.M` is from `github.com/foomo/go/testing`, not `*testing.M`.
:::

## Methods

### ExportSpans

```go
func (e *Exporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error
```

Exports a batch of spans in human-readable format. Groups spans by trace ID and renders them as a tree. Returns an error if the context is canceled.

### Shutdown

```go
func (e *Exporter) Shutdown(ctx context.Context) error
```

Shuts down the exporter. Currently a no-op.

## Options

| Option | Signature | Description | Default |
|---|---|---|---|
| `WithWriter` | `WithWriter(w io.Writer) Option` | Output destination | `os.Stdout` |
| `WithFlamegraph` | `WithFlamegraph() Option` | Enable nested flamegraph visualization | disabled |
| `WithSpanAttributes` | `WithSpanAttributes() Option` | Print span attributes and events | disabled |
| `WithoutTimestamps` | `WithoutTimestamps() Option` | Disable timestamp output | timestamps enabled |
| `WithMinDuration` | `WithMinDuration(d time.Duration) Option` | Filter spans shorter than `d` | `0` (show all) |
| `WithDurationThresholds` | `WithDurationThresholds(warn, critical time.Duration) Option` | Color-coding thresholds (green/yellow/red) | `10ms` / `100ms` |

## Example

```go
exporter, err := glossytrace.New(
	glossytrace.WithFlamegraph(),
	glossytrace.WithSpanAttributes(),
	glossytrace.WithDurationThresholds(5*time.Millisecond, 50*time.Millisecond),
)
```

Output:

```
=== TRACE 0102030405060708090a0b0c0d0e0f10 ===
└─ HTTP GET /api/users (150.00 ms) (0102030405060708090a0b0c0d0e0f10:0102030405060708)
    http.method=GET
    http.target=/api/users
  └─ db.Query (40.00 ms) (0102030405060708090a0b0c0d0e0f10:1112131415161718)

Nested Flamegraph:
HTTP GET /api/users ██████████████████████████████ 150.00 ms
  db.Query        ████████ 40.00 ms
```
