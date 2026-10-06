package semconv

import (
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// MessagingSystemNats is the messaging.system attribute for NATS.
var MessagingSystemNats = semconv.MessagingSystemKey.String("nats")

const (
	// NATSClientNameKey is the key for nats.client.name.
	NATSClientNameKey = attribute.Key("nats.client.name")
	// NATSClientErrorKindKey is the key for nats.client.error.kind.
	NATSClientErrorKindKey = attribute.Key("nats.client.error.kind")
	// MessagingNATSStreamKey is the key for messaging.nats.stream.
	MessagingNATSStreamKey = attribute.Key("messaging.nats.stream")
)

// NATSClientName returns a new attribute.KeyValue for nats.client.name.
func NATSClientName(v string) attribute.KeyValue {
	return NATSClientNameKey.String(v)
}

// NATSClientErrorKind returns a new attribute.KeyValue for nats.client.error.kind.
func NATSClientErrorKind(v string) attribute.KeyValue {
	return NATSClientErrorKindKey.String(v)
}

// MessagingNATSStream returns a new attribute.KeyValue for messaging.nats.stream.
func MessagingNATSStream(v string) attribute.KeyValue {
	return MessagingNATSStreamKey.String(v)
}
