package distributedlockconv

import (
	"context"
	"slices"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"

	"github.com/foomo/opentelemetry-go/semconv"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
)

// SystemAttr is an attribute conforming to the distributed_lock.system semantic
// conventions. It represents the backend providing the lock.
type SystemAttr string

var (
	// SystemRedis is a lock backed by Redis.
	SystemRedis SystemAttr = "redis"
	// SystemMongoDB is a lock backed by MongoDB.
	SystemMongoDB SystemAttr = "mongodb"
	// SystemNATS is a lock backed by NATS (e.g. JetStream KV).
	SystemNATS SystemAttr = "nats"
	// SystemEtcd is a lock backed by etcd.
	SystemEtcd SystemAttr = "etcd"
	// SystemPostgreSQL is a lock backed by PostgreSQL (e.g. advisory locks).
	SystemPostgreSQL SystemAttr = "postgresql"
)

// AcquireResultAttr is an attribute conforming to the
// distributed_lock.acquire.result semantic conventions. It represents the
// outcome of an acquire attempt.
type AcquireResultAttr string

var (
	// AcquireResultAcquired is the lock was acquired.
	AcquireResultAcquired AcquireResultAttr = "acquired"
	// AcquireResultContended is the lock is held by another owner and the
	// attempt gave up without waiting.
	AcquireResultContended AcquireResultAttr = "contended"
	// AcquireResultTimeout is the attempt waited for the lock and timed out.
	AcquireResultTimeout AcquireResultAttr = "timeout"
	// AcquireResultError is the attempt failed with an error from the backend.
	AcquireResultError AcquireResultAttr = "error"
)

// -----------------------------------------------------------------------------
// AcquireDuration — distributed_lock.acquire.duration
// -----------------------------------------------------------------------------

// AcquireDuration is an instrument used to record metric values conforming to
// the "distributed_lock.acquire.duration" semantic conventions. It represents
// the duration of lock acquire attempts, including time spent waiting.
type AcquireDuration struct {
	metric.Float64Histogram
}

var newAcquireDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("Duration of distributed lock acquire attempts."),
	metric.WithUnit("s"),
	metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.075, 0.1, 0.25, 0.5, 0.75, 1, 2.5, 5, 7.5, 10, 30, 60),
}

// NewAcquireDuration returns a new AcquireDuration instrument.
func NewAcquireDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (AcquireDuration, error) {
	if m == nil {
		return AcquireDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newAcquireDurationOpts
	} else {
		opt = append(opt, newAcquireDurationOpts...)
	}

	i, err := m.Float64Histogram("distributed_lock.acquire.duration", opt...)
	if err != nil {
		return AcquireDuration{noop.Float64Histogram{}}, err
	}

	return AcquireDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m AcquireDuration) Inst() metric.Float64Histogram { return m.Float64Histogram }

// Name returns the semantic convention name of the instrument.
func (AcquireDuration) Name() string { return "distributed_lock.acquire.duration" }

// Unit returns the semantic convention unit of the instrument.
func (AcquireDuration) Unit() string { return "s" }

// Description returns the semantic convention description of the instrument.
func (AcquireDuration) Description() string {
	return "Duration of distributed lock acquire attempts."
}

// Record records val to the current distribution for attrs.
//
// The name is the logical lock name, system the lock backend and result the
// outcome of the attempt.
//
// All additional attrs passed are included in the recorded value.
func (m AcquireDuration) Record(
	ctx context.Context,
	val float64,
	name string,
	system SystemAttr,
	result AcquireResultAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val, metric.WithAttributes(
			semconv.DistributedLockName(name),
			semconv.DistributedLockSystem(string(system)),
			semconv.DistributedLockAcquireResult(string(result)),
		))

		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; recOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(slices.Clip(attrs),
			semconv.DistributedLockName(name),
			semconv.DistributedLockSystem(string(system)),
			semconv.DistributedLockAcquireResult(string(result)),
		)...,
	))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m AcquireDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Enabled(ctx) {
		return
	}

	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; recOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// AttrErrorType returns an optional attribute for the "error.type" semantic
// convention. Use it together with AcquireResultError.
func (AcquireDuration) AttrErrorType(err error) attribute.KeyValue {
	return semconv.ErrorType(err)
}

// -----------------------------------------------------------------------------
// HoldDuration — distributed_lock.hold.duration
// -----------------------------------------------------------------------------

// HoldDuration is an instrument used to record metric values conforming to the
// "distributed_lock.hold.duration" semantic conventions. It represents the
// time a lock was held, from acquire until release.
type HoldDuration struct {
	metric.Float64Histogram
}

var newHoldDurationOpts = []metric.Float64HistogramOption{
	metric.WithDescription("Duration a distributed lock was held."),
	metric.WithUnit("s"),
	metric.WithExplicitBucketBoundaries(0.01, 0.05, 0.1, 0.5, 1, 5, 10, 30, 60, 300, 600, 1800, 3600),
}

// NewHoldDuration returns a new HoldDuration instrument.
func NewHoldDuration(
	m metric.Meter,
	opt ...metric.Float64HistogramOption,
) (HoldDuration, error) {
	if m == nil {
		return HoldDuration{noop.Float64Histogram{}}, nil
	}

	if len(opt) == 0 {
		opt = newHoldDurationOpts
	} else {
		opt = append(opt, newHoldDurationOpts...)
	}

	i, err := m.Float64Histogram("distributed_lock.hold.duration", opt...)
	if err != nil {
		return HoldDuration{noop.Float64Histogram{}}, err
	}

	return HoldDuration{i}, nil
}

// Inst returns the underlying metric instrument.
func (m HoldDuration) Inst() metric.Float64Histogram { return m.Float64Histogram }

// Name returns the semantic convention name of the instrument.
func (HoldDuration) Name() string { return "distributed_lock.hold.duration" }

// Unit returns the semantic convention unit of the instrument.
func (HoldDuration) Unit() string { return "s" }

// Description returns the semantic convention description of the instrument.
func (HoldDuration) Description() string {
	return "Duration a distributed lock was held."
}

// Record records val to the current distribution for attrs.
//
// The name is the logical lock name and system the lock backend.
//
// All additional attrs passed are included in the recorded value.
func (m HoldDuration) Record(
	ctx context.Context,
	val float64,
	name string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Float64Histogram.Record(ctx, val, metric.WithAttributes(
			semconv.DistributedLockName(name),
			semconv.DistributedLockSystem(string(system)),
		))

		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; recOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(slices.Clip(attrs),
			semconv.DistributedLockName(name),
			semconv.DistributedLockSystem(string(system)),
		)...,
	))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// RecordSet records val to the current distribution for set.
func (m HoldDuration) RecordSet(ctx context.Context, val float64, set attribute.Set) {
	if !m.Enabled(ctx) {
		return
	}

	if set.Len() == 0 {
		m.Float64Histogram.Record(ctx, val)
		return
	}

	o := recOptPool.Get().(*[]metric.RecordOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; recOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributeSet(set))
	m.Float64Histogram.Record(ctx, val, *o...)
}

// -----------------------------------------------------------------------------
// Lost — distributed_lock.lost
// -----------------------------------------------------------------------------

// Lost is an instrument used to record metric values conforming to the
// "distributed_lock.lost" semantic conventions. It represents the number of
// locks lost while still held, e.g. because the lease expired or was taken
// over by another owner.
type Lost struct {
	metric.Int64Counter
}

var newLostOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of distributed locks lost while still held."),
	metric.WithUnit("{lock}"),
}

// NewLost returns a new Lost instrument.
func NewLost(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (Lost, error) {
	if m == nil {
		return Lost{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newLostOpts
	} else {
		opt = append(opt, newLostOpts...)
	}

	i, err := m.Int64Counter("distributed_lock.lost", opt...)
	if err != nil {
		return Lost{noop.Int64Counter{}}, err
	}

	return Lost{i}, nil
}

// Inst returns the underlying metric instrument.
func (m Lost) Inst() metric.Int64Counter { return m.Int64Counter }

// Name returns the semantic convention name of the instrument.
func (Lost) Name() string { return "distributed_lock.lost" }

// Unit returns the semantic convention unit of the instrument.
func (Lost) Unit() string { return "{lock}" }

// Description returns the semantic convention description of the instrument.
func (Lost) Description() string {
	return "Number of distributed locks lost while still held."
}

// Add adds incr to the existing count for attrs.
//
// The name is the logical lock name and system the lock backend.
//
// All additional attrs passed are included in the recorded value.
func (m Lost) Add(
	ctx context.Context,
	incr int64,
	name string,
	system SystemAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			semconv.DistributedLockName(name),
			semconv.DistributedLockSystem(string(system)),
		))

		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(slices.Clip(attrs),
			semconv.DistributedLockName(name),
			semconv.DistributedLockSystem(string(system)),
		)...,
	))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m Lost) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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
