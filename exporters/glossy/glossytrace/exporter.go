package glossytrace

import (
	"context"
	"sort"
	"sync"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Exporter is a human-readable, styled trace exporter using lipgloss.
type Exporter struct {
	cfg config
	mu  sync.Mutex
}

// New creates a new glossy trace exporter.
func New(options ...Option) (*Exporter, error) {
	return &Exporter{
		cfg: newConfig(options),
	}, nil
}

// ExportSpans exports a batch of spans in a human-readable format.
func (e *Exporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	traces := map[string][]sdktrace.ReadOnlySpan{}

	for _, s := range spans {
		d := s.EndTime().Sub(s.StartTime())
		if e.cfg.minDuration > 0 && d < e.cfg.minDuration {
			continue
		}

		traceID := s.SpanContext().TraceID().String()
		traces[traceID] = append(traces[traceID], s)
	}

	// Sort trace IDs for deterministic output.
	traceIDs := make([]string, 0, len(traces))
	for id := range traces {
		traceIDs = append(traceIDs, id)
	}

	sort.Strings(traceIDs)

	for _, id := range traceIDs {
		e.printTrace(id, traces[id])
	}

	return nil
}

// Shutdown shuts down the exporter.
func (e *Exporter) Shutdown(_ context.Context) error {
	return nil
}
