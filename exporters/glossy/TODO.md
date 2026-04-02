# Glossy Exporter - Deferred Features

## Bench Mode (Metric Aggregation)

The `tmp/exporters/humanmetric` prototype includes a "bench mode" that accumulates metrics across
multiple `Export()` calls and prints a summary (Count, Sum, Min, Max, Avg) on demand.

This behavior is better implemented as a **custom `sdkmetric.Reader`** or a wrapper around an
existing reader that:

1. Intercepts `Export()` calls and accumulates data points per metric name
2. Exposes a `PrintSummary(w io.Writer)` method for on-demand output
3. Delegates to the underlying glossymetric exporter for real-time display when not in bench mode

This keeps the exporter focused on rendering and moves aggregation logic to where it belongs in
the OTel SDK pipeline.

## Latency Histogram (Trace)

The `tmp/exporters/humantrace` prototype collects span durations and prints a latency histogram
on `Shutdown()`.

This is better implemented as a **custom `sdktrace.SpanProcessor`** that:

1. Implements `OnEnd(s sdktrace.ReadOnlySpan)` to collect span durations
2. Maintains configurable histogram buckets (e.g., 1us, 10us, 100us, 1ms, 10ms, 100ms, 1s)
3. Prints the histogram on `Shutdown()` or exposes a `PrintHistogram(w io.Writer)` method
4. Can be registered alongside the glossytrace exporter via `sdktrace.WithSpanProcessor()`

This separates collection/aggregation from export rendering and follows OTel's processor model.

## Log Exporter (glossylog)

Add a `glossylog` subpackage under `exporters/glossy/` following the same pattern as glossytrace
and glossymetric:

- Implement `log.Exporter` interface (`Export`, `Shutdown`, `ForceFlush`)
- Pretty-print log records with lipgloss styling (severity colors, structured attribute rendering)
- Separate `go.mod` at `exporters/glossy/glossylog/`
