---
title: glossymetric
description: API reference for the glossy metric exporter
---

# glossymetric

```go
import "github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric"
```

Styled terminal exporter for OpenTelemetry metrics. Renders Sum, Gauge, and Histogram data in a readable format.

## Types

### Exporter

```go
type Exporter struct {
	// contains filtered or unexported fields
}
```

Implements `go.opentelemetry.io/otel/sdk/metric.Exporter`.

## Constructors

### New

```go
func New(opts ...Option) (*Exporter, error)
```

Creates a new glossy metric exporter. Output defaults to `os.Stdout`.

### NewTest

```go
func NewTest(tb testing.TB, opts ...Option) metric.Exporter
```

Creates an exporter for use in tests. Automatically applies `WithWriter(tb.Output())`. Calls `tb.Fatal()` on initialization error.

### NewTestMain

```go
func NewTestMain(m testingx.M, opts ...Option) metric.Exporter
```

Creates an exporter for use in `TestMain`. Automatically applies `WithWriter(os.Stdout)`. Panics on initialization error.

::: tip
`testingx.M` is from `github.com/foomo/go/testing`, not `*testing.M`.
:::

## Methods

### Export

```go
func (e *Exporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error
```

Exports metrics in human-readable format. Returns an error if the context is canceled.

### ForceFlush

```go
func (e *Exporter) ForceFlush(ctx context.Context) error
```

No-op. Required by the `metric.Exporter` interface.

### Shutdown

```go
func (e *Exporter) Shutdown(ctx context.Context) error
```

Shuts down the exporter. Currently a no-op.

### Temporality

```go
func (e *Exporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality
```

Returns the temporality for the given instrument kind using the configured selector.

### Aggregation

```go
func (e *Exporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation
```

Returns the aggregation for the given instrument kind using the configured selector.

## Options

| Option | Signature | Description | Default |
|---|---|---|---|
| `WithWriter` | `WithWriter(w io.Writer) Option` | Output destination | `os.Stdout` |
| `WithoutHistograms` | `WithoutHistograms() Option` | Disable histogram detail printing | histograms enabled |
| `WithTemporalitySelector` | `WithTemporalitySelector(fn sdkmetric.TemporalitySelector) Option` | Set temporality selector | Cumulative |
| `WithAggregationSelector` | `WithAggregationSelector(fn sdkmetric.AggregationSelector) Option` | Set aggregation selector | `DefaultAggregationSelector` |

## Example Output

```
METRICS

Scope: example
──────────────────────────────────────────────────
Metric: requests
Description: Number of requests received
Unit: 1
  attrs: server=central
Value: 5
──────────────────────────────────────────────────
Metric: requests.size
Description: Size of received requests
Unit: kb
  attrs: server=central
Histogram: min=3 max=30 avg=12.800 count=10
──────────────────────────────────────────────────
```
