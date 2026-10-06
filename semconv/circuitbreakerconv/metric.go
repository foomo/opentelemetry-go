package circuitbreakerconv

import (
	"context"
	"slices"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/foomo/opentelemetry-go/semconv"
)

var addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}

// StateAttr is an attribute conforming to the circuit_breaker.state semantic
// conventions. It represents the state of a circuit breaker.
type StateAttr string

var (
	// StateClosed is the breaker is passing calls through.
	StateClosed StateAttr = "closed"
	// StateOpen is the breaker is rejecting calls.
	StateOpen StateAttr = "open"
	// StateHalfOpen is the breaker is letting a limited number of trial calls
	// through to probe recovery.
	StateHalfOpen StateAttr = "half_open"
)

// -----------------------------------------------------------------------------
// StateChanges — circuit_breaker.state_changes
// -----------------------------------------------------------------------------

// StateChanges is an instrument used to record metric values conforming to the
// "circuit_breaker.state_changes" semantic conventions. It represents the
// number of state transitions of a circuit breaker.
type StateChanges struct {
	metric.Int64Counter
}

var newStateChangesOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of circuit breaker state transitions."),
	metric.WithUnit("{change}"),
}

// NewStateChanges returns a new StateChanges instrument.
func NewStateChanges(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (StateChanges, error) {
	if m == nil {
		return StateChanges{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newStateChangesOpts
	} else {
		opt = append(opt, newStateChangesOpts...)
	}

	i, err := m.Int64Counter("circuit_breaker.state_changes", opt...)
	if err != nil {
		return StateChanges{noop.Int64Counter{}}, err
	}

	return StateChanges{i}, nil
}

// Inst returns the underlying metric instrument.
func (m StateChanges) Inst() metric.Int64Counter { return m.Int64Counter }

// Name returns the semantic convention name of the instrument.
func (StateChanges) Name() string { return "circuit_breaker.state_changes" }

// Unit returns the semantic convention unit of the instrument.
func (StateChanges) Unit() string { return "{change}" }

// Description returns the semantic convention description of the instrument.
func (StateChanges) Description() string {
	return "Number of circuit breaker state transitions."
}

// Add adds incr to the existing count for attrs.
//
// The name identifies the circuit breaker and state is the state it
// transitioned to.
//
// All additional attrs passed are included in the recorded value.
func (m StateChanges) Add(
	ctx context.Context,
	incr int64,
	name string,
	state StateAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			semconv.CircuitBreakerName(name),
			semconv.CircuitBreakerState(string(state)),
		))

		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(slices.Clip(attrs),
			semconv.CircuitBreakerName(name),
			semconv.CircuitBreakerState(string(state)),
		)...,
	))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m StateChanges) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Enabled(ctx) {
		return
	}

	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// -----------------------------------------------------------------------------
// Rejections — circuit_breaker.rejections
// -----------------------------------------------------------------------------

// Rejections is an instrument used to record metric values conforming to the
// "circuit_breaker.rejections" semantic conventions. It represents the number
// of calls rejected without execution because the circuit breaker was open or
// at its half-open trial limit.
type Rejections struct {
	metric.Int64Counter
}

var newRejectionsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of calls rejected by the circuit breaker."),
	metric.WithUnit("{call}"),
}

// NewRejections returns a new Rejections instrument.
func NewRejections(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Rejections, error) {
	if m == nil {
		return Rejections{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newRejectionsOpts
	} else {
		opt = append(opt, newRejectionsOpts...)
	}

	i, err := m.Int64Counter("circuit_breaker.rejections", opt...)
	if err != nil {
		return Rejections{noop.Int64Counter{}}, err
	}

	return Rejections{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Rejections) Inst() metric.Int64Counter { return m.Int64Counter }

// Name returns the semantic convention name of the instrument.
func (Rejections) Name() string { return "circuit_breaker.rejections" }

// Unit returns the semantic convention unit of the instrument.
func (Rejections) Unit() string { return "{call}" }

// Description returns the semantic convention description of the instrument.
func (Rejections) Description() string {
	return "Number of calls rejected by the circuit breaker."
}

// Add adds incr to the existing count for attrs.
//
// The name identifies the circuit breaker.
//
// All additional attrs passed are included in the recorded value.
func (m Rejections) Add(
	ctx context.Context,
	incr int64,
	name string,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			semconv.CircuitBreakerName(name),
		))

		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(slices.Clip(attrs),
			semconv.CircuitBreakerName(name),
		)...,
	))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m Rejections) AddSet(ctx context.Context, incr int64, set attribute.Set) {
	if !m.Enabled(ctx) {
		return
	}

	if set.Len() == 0 {
		m.Int64Counter.Add(ctx, incr)
		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AttrState returns an optional attribute for the "circuit_breaker.state"
// semantic convention. It represents the breaker state at rejection time.
func (Rejections) AttrState(val StateAttr) attribute.KeyValue {
	return semconv.CircuitBreakerState(string(val))
}
