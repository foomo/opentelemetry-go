package glossymetric

import (
	"context"
	"sync"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Exporter is a human-readable, styled metric exporter using lipgloss.
// It implements [sdkmetric.Exporter] and is safe for concurrent use.
// The zero value is not usable; create one with [New].
type Exporter struct {
	cfg config
	mu  sync.Mutex
}

// New returns a new Exporter configured by opts. It writes to os.Stdout unless
// [WithWriter] is given. The returned error is always nil.
func New(opts ...Option) (*Exporter, error) {
	return &Exporter{
		cfg: newConfig(opts),
	}, nil
}

// Export writes rm to the configured writer in a human-readable format.
// It returns ctx.Err() without writing if ctx is already done. Write errors
// are ignored.
func (e *Exporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	return e.prettyPrint(rm)
}

// ForceFlush is a no-op because Export writes synchronously. It always
// returns nil.
func (e *Exporter) ForceFlush(_ context.Context) error {
	return nil
}

// Shutdown is a no-op; the exporter holds no resources and keeps accepting
// exports afterwards. It always returns nil.
func (e *Exporter) Shutdown(_ context.Context) error {
	return nil
}

// Temporality returns the temporality for kind as chosen by the selector set
// with [WithTemporalitySelector]. It defaults to cumulative temporality.
func (e *Exporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	return e.cfg.temporalitySelector(kind)
}

// Aggregation returns the aggregation for kind as chosen by the selector set
// with [WithAggregationSelector]. It defaults to
// [sdkmetric.DefaultAggregationSelector].
func (e *Exporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return e.cfg.aggregationSelector(kind)
}
