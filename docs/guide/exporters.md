---
title: Exporters
description: Overview of the glossytrace and glossymetric terminal exporters
---

# Exporters

The library ships two exporters that render OpenTelemetry data as styled terminal output using [lipgloss](https://github.com/charmbracelet/lipgloss):

- **glossytrace** -- trace spans as tree views with optional flamegraphs
- **glossymetric** -- metrics (Sum, Gauge, Histogram) in a readable table format

Both follow the **functional options** pattern used throughout the OpenTelemetry Go ecosystem.

## Functional Options Pattern

Each exporter defines an `Option` interface and `With*` constructor functions:

```go
// Option applies a configuration option to the Exporter.
type Option interface {
	apply(cfg config) config
}

// Usage
exporter, err := glossytrace.New(
	glossytrace.WithFlamegraph(),
	glossytrace.WithMinDuration(time.Millisecond),
)
```

This pattern keeps the API extensible without breaking changes. The `config` struct is unexported -- all configuration flows through options.

## glossytrace

The trace exporter renders spans in a tree structure showing parent-child relationships:

```
=== TRACE 0102030405060708090a0b0c0d0e0f10 ===
└─ HTTP GET /api/users (150.00 ms) (0102030405060708090a0b0c0d0e0f10:0102030405060708)
  └─ db.Query (40.00 ms) (0102030405060708090a0b0c0d0e0f10:1112131415161718)
```

### Key Features

**Duration color-coding** -- Spans are colored based on configurable thresholds:
- Green: faster than `warn` (default 10ms)
- Yellow: between `warn` and `critical`
- Red: slower than `critical` (default 100ms)

**Flamegraph mode** -- Enable with `WithFlamegraph()` to get a proportional bar chart after the tree:

```
Nested Flamegraph:
HTTP GET /api/users ██████████████████████████████ 150.00 ms
  db.Query        ████████ 40.00 ms
```

**Span attributes** -- Enable with `WithSpanAttributes()` to print attributes and events below each span.

**Min-duration filter** -- Use `WithMinDuration(d)` to suppress short spans and focus on slow operations.

See the full option list in the [glossytrace API reference](/api/glossytrace).

## glossymetric

The metric exporter renders all metric types in a structured format:

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
Metric: latency
Description: Time spend processing received requests
Unit: ms
  attrs: server=central
Histogram: min=0.000 max=0.000 avg=5.700 count=10
```

### Supported Metric Types

- **Sum[int64]** and **Sum[float64]** -- counters and up-down counters
- **Gauge[int64]** and **Gauge[float64]** -- point-in-time values
- **Histogram[int64]** and **Histogram[float64]** -- distribution data with min/max/avg/count

### Key Options

- `WithoutHistograms()` -- suppress histogram detail rendering
- `WithTemporalitySelector(fn)` -- override the default cumulative temporality
- `WithAggregationSelector(fn)` -- override the default aggregation strategy

See the full option list in the [glossymetric API reference](/api/glossymetric).

## Using Exporters in Tests

::: warning
In tests, prefer `NewTest` and `NewTestMain` over `New`. These helpers automatically set the writer to `tb.Output()` (or `os.Stdout` for TestMain) and disable timestamps for deterministic output.
:::

```go
// In a test function
exporter := glossytrace.NewTest(t, glossytrace.WithFlamegraph())

// In TestMain
exporter := glossymetric.NewTestMain(m)
```

See the [Testing Guide](/guide/testing) for complete examples.
