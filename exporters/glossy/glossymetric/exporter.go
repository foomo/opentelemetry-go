package glossymetric

import (
	"context"
	"sync"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Exporter is a human-readable, styled metric exporter using lipgloss.
type Exporter struct {
	cfg config
	mu  sync.Mutex
}

// New creates a new glossy metric exporter.
func New(opts ...Option) (*Exporter, error) {
	return &Exporter{
		cfg: newConfig(opts),
	}, nil
}

// Export exports metrics in a human-readable format.
func (e *Exporter) Export(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	return e.prettyPrint(rm)
}

// ForceFlush is a no-op.
func (e *Exporter) ForceFlush(_ context.Context) error {
	return nil
}

// Shutdown shuts down the exporter.
func (e *Exporter) Shutdown(_ context.Context) error {
	return nil
}

// Temporality returns the configured temporality selector.
func (e *Exporter) Temporality(kind sdkmetric.InstrumentKind) metricdata.Temporality {
	return e.cfg.temporalitySelector(kind)
}

// Aggregation returns the configured aggregation selector.
func (e *Exporter) Aggregation(kind sdkmetric.InstrumentKind) sdkmetric.Aggregation {
	return e.cfg.aggregationSelector(kind)
}
