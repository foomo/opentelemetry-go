// Package natsconv provides types and functionality for NATS-specific metrics
// following the OpenTelemetry semantic conventions pattern in the "nats" and
// "messaging.nats" namespaces. This package extends messagingconv with
// NATS-specific instruments (JetStream consumer state, core connection health)
// that are not covered by the generic messaging semantic conventions.
package natsconv

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
)

var (
	addOptPool = &sync.Pool{New: func() any { return &[]metric.AddOption{} }}
	recOptPool = &sync.Pool{New: func() any { return &[]metric.RecordOption{} }}
)

// ConnectionStatusAttr is an attribute conforming to the
// nats.client.connection.status semantic conventions. It represents the state
// of a NATS client connection.
type ConnectionStatusAttr string

var (
	// ConnectionStatusDisconnected is the client is not currently connected.
	ConnectionStatusDisconnected ConnectionStatusAttr = "disconnected"
	// ConnectionStatusConnected is the client is connected to a server.
	ConnectionStatusConnected ConnectionStatusAttr = "connected"
	// ConnectionStatusClosed is the connection has been closed permanently.
	ConnectionStatusClosed ConnectionStatusAttr = "closed"
	// ConnectionStatusReconnecting is the client is attempting to reconnect.
	ConnectionStatusReconnecting ConnectionStatusAttr = "reconnecting"
	// ConnectionStatusConnecting is the client is establishing its initial
	// connection.
	ConnectionStatusConnecting ConnectionStatusAttr = "connecting"
	// ConnectionStatusDraining is the client is draining subscriptions or
	// publishers prior to close.
	ConnectionStatusDraining ConnectionStatusAttr = "draining"
)

// AsyncErrorKindAttr is an attribute conforming to the nats.client.error.kind
// semantic conventions. It represents a low-cardinality classification of async
// errors reported by the NATS client.
type AsyncErrorKindAttr string

var (
	// AsyncErrorKindSlowConsumer indicates a subscription could not keep up
	// with incoming messages and messages were dropped.
	AsyncErrorKindSlowConsumer AsyncErrorKindAttr = "slow_consumer"
	// AsyncErrorKindPermissionViolation indicates the server rejected a
	// publish or subscribe due to permissions.
	AsyncErrorKindPermissionViolation AsyncErrorKindAttr = "permission_violation"
	// AsyncErrorKindAuthExpired indicates credentials have expired.
	AsyncErrorKindAuthExpired AsyncErrorKindAttr = "auth_expired"
	// AsyncErrorKindAuthRevoked indicates credentials have been revoked.
	AsyncErrorKindAuthRevoked AsyncErrorKindAttr = "auth_revoked"
	// AsyncErrorKindOther is a fallback when no classification applies.
	AsyncErrorKindOther AsyncErrorKindAttr = "_OTHER"
)

// -----------------------------------------------------------------------------
// ClientDisconnects — nats.client.disconnects
// -----------------------------------------------------------------------------

// ClientDisconnects is an instrument used to record metric values conforming to
// the "nats.client.disconnects" semantic conventions. It represents the number
// of disconnect events observed by the NATS client.
type ClientDisconnects struct {
	metric.Int64Counter
}

var newClientDisconnectsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of disconnect events observed by the NATS client."),
	metric.WithUnit("{event}"),
}

// NewClientDisconnects returns a new ClientDisconnects instrument.
func NewClientDisconnects(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientDisconnects, error) {
	if m == nil {
		return ClientDisconnects{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientDisconnectsOpts
	} else {
		opt = append(opt, newClientDisconnectsOpts...)
	}

	i, err := m.Int64Counter("nats.client.disconnects", opt...)
	if err != nil {
		return ClientDisconnects{noop.Int64Counter{}}, err
	}

	return ClientDisconnects{i}, nil
}

// Inst returns the underlying metric instrument.
func (m ClientDisconnects) Inst() metric.Int64Counter { return m.Int64Counter }

// Name returns the semantic convention name of the instrument.
func (ClientDisconnects) Name() string { return "nats.client.disconnects" }

// Unit returns the semantic convention unit of the instrument.
func (ClientDisconnects) Unit() string { return "{event}" }

// Description returns the semantic convention description of the instrument.
func (ClientDisconnects) Description() string {
	return "Number of disconnect events observed by the NATS client."
}

// Add adds incr to the existing count for attrs.
//
// The serverAddress is the server the client was connected to at the time of
// the event.
//
// All additional attrs passed are included in the recorded value.
func (m ClientDisconnects) Add(
	ctx context.Context,
	incr int64,
	serverAddress string,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("server.address", serverAddress),
		))

		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(attrs[:len(attrs):len(attrs)],
			attribute.String("server.address", serverAddress),
		)...,
	))
	m.Int64Counter.Add(ctx, incr, *o...)
}

// AddSet adds incr to the existing count for set.
func (m ClientDisconnects) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrServerPort returns an optional attribute for the "server.port" semantic
// convention.
func (ClientDisconnects) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}

// AttrClientName returns an optional attribute for the "nats.client.name"
// semantic convention. It represents the client name registered via
// nats.Name(...) on Connect.
func (ClientDisconnects) AttrClientName(val string) attribute.KeyValue {
	return attribute.String("nats.client.name", val)
}

// -----------------------------------------------------------------------------
// ClientReconnects — nats.client.reconnects
// -----------------------------------------------------------------------------

// ClientReconnects is an instrument used to record metric values conforming to
// the "nats.client.reconnects" semantic conventions. It represents the number
// of successful reconnects performed by the NATS client.
type ClientReconnects struct {
	metric.Int64Counter
}

var newClientReconnectsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of successful reconnects performed by the NATS client."),
	metric.WithUnit("{event}"),
}

// NewClientReconnects returns a new ClientReconnects instrument.
func NewClientReconnects(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientReconnects, error) {
	if m == nil {
		return ClientReconnects{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientReconnectsOpts
	} else {
		opt = append(opt, newClientReconnectsOpts...)
	}

	i, err := m.Int64Counter("nats.client.reconnects", opt...)
	if err != nil {
		return ClientReconnects{noop.Int64Counter{}}, err
	}

	return ClientReconnects{i}, nil
}

func (m ClientReconnects) Inst() metric.Int64Counter { return m.Int64Counter }
func (ClientReconnects) Name() string                { return "nats.client.reconnects" }
func (ClientReconnects) Unit() string                { return "{event}" }
func (ClientReconnects) Description() string {
	return "Number of successful reconnects performed by the NATS client."
}

// Add adds incr to the existing count for attrs.
//
// The serverAddress is the server the client reconnected to.
func (m ClientReconnects) Add(
	ctx context.Context,
	incr int64,
	serverAddress string,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("server.address", serverAddress),
		))

		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(attrs[:len(attrs):len(attrs)],
			attribute.String("server.address", serverAddress),
		)...,
	))
	m.Int64Counter.Add(ctx, incr, *o...)
}

func (m ClientReconnects) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

func (ClientReconnects) AttrServerPort(val int) attribute.KeyValue {
	return attribute.Int("server.port", val)
}
func (ClientReconnects) AttrClientName(val string) attribute.KeyValue {
	return attribute.String("nats.client.name", val)
}

// -----------------------------------------------------------------------------
// ClientAsyncErrors — nats.client.async_errors
// -----------------------------------------------------------------------------

// ClientAsyncErrors is an instrument used to record metric values conforming to
// the "nats.client.async_errors" semantic conventions. It represents the number
// of asynchronous errors reported by the NATS client (slow consumer,
// permission violations, auth expiry, etc.).
type ClientAsyncErrors struct {
	metric.Int64Counter
}

var newClientAsyncErrorsOpts = []metric.Int64CounterOption{
	metric.WithDescription("Number of asynchronous errors reported by the NATS client."),
	metric.WithUnit("{error}"),
}

// NewClientAsyncErrors returns a new ClientAsyncErrors instrument.
func NewClientAsyncErrors(
	m metric.Meter,
	opt ...metric.Int64CounterOption,
) (ClientAsyncErrors, error) {
	if m == nil {
		return ClientAsyncErrors{noop.Int64Counter{}}, nil
	}

	if len(opt) == 0 {
		opt = newClientAsyncErrorsOpts
	} else {
		opt = append(opt, newClientAsyncErrorsOpts...)
	}

	i, err := m.Int64Counter("nats.client.async_errors", opt...)
	if err != nil {
		return ClientAsyncErrors{noop.Int64Counter{}}, err
	}

	return ClientAsyncErrors{i}, nil
}

func (m ClientAsyncErrors) Inst() metric.Int64Counter { return m.Int64Counter }
func (ClientAsyncErrors) Name() string                { return "nats.client.async_errors" }
func (ClientAsyncErrors) Unit() string                { return "{error}" }
func (ClientAsyncErrors) Description() string {
	return "Number of asynchronous errors reported by the NATS client."
}

// Add adds incr to the existing count for attrs.
//
// The kind is a low-cardinality classification of the error.
func (m ClientAsyncErrors) Add(
	ctx context.Context,
	incr int64,
	kind AsyncErrorKindAttr,
	attrs ...attribute.KeyValue,
) {
	if !m.Enabled(ctx) {
		return
	}

	if len(attrs) == 0 {
		m.Int64Counter.Add(ctx, incr, metric.WithAttributes(
			attribute.String("nats.client.error.kind", string(kind)),
		))

		return
	}

	o := addOptPool.Get().(*[]metric.AddOption) //nolint:forcetypeassert

	defer func() { *o = (*o)[:0]; addOptPool.Put(o) }()

	*o = append(*o, metric.WithAttributes(
		append(attrs[:len(attrs):len(attrs)],
			attribute.String("nats.client.error.kind", string(kind)),
		)...,
	))
	m.Int64Counter.Add(ctx, incr, *o...)
}

func (m ClientAsyncErrors) AddSet(ctx context.Context, incr int64, set attribute.Set) {
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

// AttrSubject returns an optional attribute for the
// "messaging.destination.name" semantic convention. For async errors it
// represents the subject on which the error was reported, if available.
func (ClientAsyncErrors) AttrSubject(val string) attribute.KeyValue {
	return attribute.String("messaging.destination.name", val)
}

// AttrServerAddress returns an optional attribute for the "server.address"
// semantic convention.
func (ClientAsyncErrors) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// -----------------------------------------------------------------------------
// JetStreamConsumerPending — nats.jetstream.consumer.pending
// -----------------------------------------------------------------------------

// JetStreamConsumerPending is an instrument used to record metric values
// conforming to the "nats.jetstream.consumer.pending" semantic conventions. It
// represents the number of messages in the stream that have not yet been
// delivered to the consumer.
type JetStreamConsumerPending struct {
	metric.Int64ObservableGauge
}

var newJetStreamConsumerPendingOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Messages in the stream not yet delivered to the consumer."),
	metric.WithUnit("{message}"),
}

// NewJetStreamConsumerPending returns a new JetStreamConsumerPending
// instrument. The caller is responsible for registering a callback with the
// meter that observes this instrument.
func NewJetStreamConsumerPending(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (JetStreamConsumerPending, error) {
	if m == nil {
		return JetStreamConsumerPending{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newJetStreamConsumerPendingOpts
	} else {
		opt = append(opt, newJetStreamConsumerPendingOpts...)
	}

	i, err := m.Int64ObservableGauge("nats.jetstream.consumer.pending", opt...)
	if err != nil {
		return JetStreamConsumerPending{noop.Int64ObservableGauge{}}, err
	}

	return JetStreamConsumerPending{i}, nil
}

func (m JetStreamConsumerPending) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}
func (JetStreamConsumerPending) Name() string { return "nats.jetstream.consumer.pending" }
func (JetStreamConsumerPending) Unit() string { return "{message}" }
func (JetStreamConsumerPending) Description() string {
	return "Messages in the stream not yet delivered to the consumer."
}

// Observe records val for the given stream and consumer within an async
// callback registered with the meter.
//
// The stream is the JetStream stream name.
//
// The consumerGroupName is the durable consumer name (conforms to
// messaging.consumer.group.name).
func (m JetStreamConsumerPending) Observe(
	o metric.Int64Observer,
	val int64,
	stream string,
	consumerGroupName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		o.Observe(val, metric.WithAttributes(
			attribute.String("messaging.nats.stream", stream),
			attribute.String("messaging.consumer.group.name", consumerGroupName),
		))

		return
	}

	o.Observe(val, metric.WithAttributes(
		append(attrs[:len(attrs):len(attrs)],
			attribute.String("messaging.nats.stream", stream),
			attribute.String("messaging.consumer.group.name", consumerGroupName),
		)...,
	))
}

// ObserveSet observes val for set within an async callback.
func (m JetStreamConsumerPending) ObserveSet(o metric.Int64Observer, val int64, set attribute.Set) {
	o.Observe(val, metric.WithAttributeSet(set))
}

func (JetStreamConsumerPending) AttrServerAddress(val string) attribute.KeyValue {
	return attribute.String("server.address", val)
}

// -----------------------------------------------------------------------------
// JetStreamConsumerAckPending — nats.jetstream.consumer.ack_pending
// -----------------------------------------------------------------------------

// JetStreamConsumerAckPending represents the number of messages that have been
// delivered but not yet acknowledged.
type JetStreamConsumerAckPending struct {
	metric.Int64ObservableGauge
}

var newJetStreamConsumerAckPendingOpts = []metric.Int64ObservableGaugeOption{
	metric.WithDescription("Messages delivered to the consumer but not yet acknowledged."),
	metric.WithUnit("{message}"),
}

func NewJetStreamConsumerAckPending(
	m metric.Meter,
	opt ...metric.Int64ObservableGaugeOption,
) (JetStreamConsumerAckPending, error) {
	if m == nil {
		return JetStreamConsumerAckPending{noop.Int64ObservableGauge{}}, nil
	}

	if len(opt) == 0 {
		opt = newJetStreamConsumerAckPendingOpts
	} else {
		opt = append(opt, newJetStreamConsumerAckPendingOpts...)
	}

	i, err := m.Int64ObservableGauge("nats.jetstream.consumer.ack_pending", opt...)
	if err != nil {
		return JetStreamConsumerAckPending{noop.Int64ObservableGauge{}}, err
	}

	return JetStreamConsumerAckPending{i}, nil
}

func (m JetStreamConsumerAckPending) Inst() metric.Int64ObservableGauge {
	return m.Int64ObservableGauge
}
func (JetStreamConsumerAckPending) Name() string { return "nats.jetstream.consumer.ack_pending" }
func (JetStreamConsumerAckPending) Unit() string { return "{message}" }
func (JetStreamConsumerAckPending) Description() string {
	return "Messages delivered to the consumer but not yet acknowledged."
}

func (m JetStreamConsumerAckPending) Observe(
	o metric.Int64Observer,
	val int64,
	stream string,
	consumerGroupName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		o.Observe(val, metric.WithAttributes(
			attribute.String("messaging.nats.stream", stream),
			attribute.String("messaging.consumer.group.name", consumerGroupName),
		))

		return
	}

	o.Observe(val, metric.WithAttributes(
		append(attrs[:len(attrs):len(attrs)],
			attribute.String("messaging.nats.stream", stream),
			attribute.String("messaging.consumer.group.name", consumerGroupName),
		)...,
	))
}

func (m JetStreamConsumerAckPending) ObserveSet(o metric.Int64Observer, val int64, set attribute.Set) {
	o.Observe(val, metric.WithAttributeSet(set))
}

// -----------------------------------------------------------------------------
// JetStreamConsumerRedelivered — nats.jetstream.consumer.redelivered
// -----------------------------------------------------------------------------

// JetStreamConsumerRedelivered represents the total number of messages
// redelivered to the consumer.
type JetStreamConsumerRedelivered struct {
	metric.Int64ObservableCounter
}

var newJetStreamConsumerRedeliveredOpts = []metric.Int64ObservableCounterOption{
	metric.WithDescription("Messages that have been redelivered to the consumer."),
	metric.WithUnit("{message}"),
}

func NewJetStreamConsumerRedelivered(
	m metric.Meter,
	opt ...metric.Int64ObservableCounterOption,
) (JetStreamConsumerRedelivered, error) {
	if m == nil {
		return JetStreamConsumerRedelivered{noop.Int64ObservableCounter{}}, nil
	}

	if len(opt) == 0 {
		opt = newJetStreamConsumerRedeliveredOpts
	} else {
		opt = append(opt, newJetStreamConsumerRedeliveredOpts...)
	}

	i, err := m.Int64ObservableCounter("nats.jetstream.consumer.redelivered", opt...)
	if err != nil {
		return JetStreamConsumerRedelivered{noop.Int64ObservableCounter{}}, err
	}

	return JetStreamConsumerRedelivered{i}, nil
}

func (m JetStreamConsumerRedelivered) Inst() metric.Int64ObservableCounter {
	return m.Int64ObservableCounter
}
func (JetStreamConsumerRedelivered) Name() string { return "nats.jetstream.consumer.redelivered" }
func (JetStreamConsumerRedelivered) Unit() string { return "{message}" }
func (JetStreamConsumerRedelivered) Description() string {
	return "Messages that have been redelivered to the consumer."
}

func (m JetStreamConsumerRedelivered) Observe(
	o metric.Int64Observer,
	val int64,
	stream string,
	consumerGroupName string,
	attrs ...attribute.KeyValue,
) {
	if len(attrs) == 0 {
		o.Observe(val, metric.WithAttributes(
			attribute.String("messaging.nats.stream", stream),
			attribute.String("messaging.consumer.group.name", consumerGroupName),
		))

		return
	}

	o.Observe(val, metric.WithAttributes(
		append(attrs[:len(attrs):len(attrs)],
			attribute.String("messaging.nats.stream", stream),
			attribute.String("messaging.consumer.group.name", consumerGroupName),
		)...,
	))
}

func (m JetStreamConsumerRedelivered) ObserveSet(o metric.Int64Observer, val int64, set attribute.Set) {
	o.Observe(val, metric.WithAttributeSet(set))
}

// Keep the unused pool symbol live in case a future instrument needs it.
var _ = recOptPool
