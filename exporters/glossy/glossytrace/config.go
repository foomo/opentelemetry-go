package glossytrace

import (
	"io"
	"os"
	"time"
)

const (
	defaultDurationWarn     = 10 * time.Millisecond
	defaultDurationCritical = 100 * time.Millisecond
)

type config struct {
	writer           io.Writer
	flamegraph       bool
	showAttributes   bool
	timestamps       bool
	minDuration      time.Duration
	durationWarn     time.Duration
	durationCritical time.Duration
}

func newConfig(options []Option) config {
	cfg := config{
		writer:           os.Stdout,
		timestamps:       true,
		durationWarn:     defaultDurationWarn,
		durationCritical: defaultDurationCritical,
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

// WithFlamegraph enables the nested flamegraph visualization.
func WithFlamegraph() Option {
	return optionFunc(func(cfg config) config {
		cfg.flamegraph = true
		return cfg
	})
}

// WithSpanAttributes enables printing of span attributes and events.
func WithSpanAttributes() Option {
	return optionFunc(func(cfg config) config {
		cfg.showAttributes = true
		return cfg
	})
}

// WithoutTimestamps disables timestamp output.
func WithoutTimestamps() Option {
	return optionFunc(func(cfg config) config {
		cfg.timestamps = false
		return cfg
	})
}

// WithMinDuration sets the minimum span duration to display.
// Spans shorter than this duration are filtered out.
func WithMinDuration(d time.Duration) Option {
	return optionFunc(func(cfg config) config {
		cfg.minDuration = d
		return cfg
	})
}

// WithDurationThresholds sets the duration thresholds for color-coding.
// Spans faster than warn are green, between warn and critical are yellow,
// and slower than critical are red.
func WithDurationThresholds(warn, critical time.Duration) Option {
	return optionFunc(func(cfg config) config {
		cfg.durationWarn = warn
		cfg.durationCritical = critical

		return cfg
	})
}
