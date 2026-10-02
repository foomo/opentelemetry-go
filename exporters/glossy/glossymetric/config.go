package glossymetric

import (
	"io"
	"os"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type config struct {
	writer              io.Writer
	histograms          bool
	temporalitySelector sdkmetric.TemporalitySelector
	aggregationSelector sdkmetric.AggregationSelector
}

func newConfig(options []Option) config {
	cfg := config{
		writer:     os.Stdout,
		histograms: true,
		temporalitySelector: func(_ sdkmetric.InstrumentKind) metricdata.Temporality {
			return metricdata.CumulativeTemporality
		},
		aggregationSelector: sdkmetric.DefaultAggregationSelector,
	}
	for _, opt := range options {
		cfg = opt.apply(cfg)
	}

	return cfg
}

// Option applies a configuration option to the Exporter.
type Option interface {
	apply(cfg config) config
}

type optionFunc func(config) config

func (f optionFunc) apply(cfg config) config { return f(cfg) }

// WithWriter sets the output destination. Default is os.Stdout.
func WithWriter(w io.Writer) Option {
	return optionFunc(func(cfg config) config {
		cfg.writer = w
		return cfg
	})
}

// WithoutHistograms disables printing of histogram data points. Histogram
// metrics still print their name, description and unit.
func WithoutHistograms() Option {
	return optionFunc(func(cfg config) config {
		cfg.histograms = false
		return cfg
	})
}

// WithTemporalitySelector sets the selector used by [Exporter.Temporality].
// Default is cumulative temporality for all instrument kinds.
func WithTemporalitySelector(fn sdkmetric.TemporalitySelector) Option {
	return optionFunc(func(cfg config) config {
		cfg.temporalitySelector = fn
		return cfg
	})
}

// WithAggregationSelector sets the selector used by [Exporter.Aggregation].
// Default is [sdkmetric.DefaultAggregationSelector].
func WithAggregationSelector(fn sdkmetric.AggregationSelector) Option {
	return optionFunc(func(cfg config) config {
		cfg.aggregationSelector = fn
		return cfg
	})
}
