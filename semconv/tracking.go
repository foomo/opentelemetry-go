package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	// TrackingIDKey is the key for tracking.id. Its values are unbounded: use
	// it on spans and logs only, never on metrics.
	TrackingIDKey = attribute.Key("tracking.id")
)

// TrackingID returns a new attribute.KeyValue for tracking.id.
func TrackingID(v string) attribute.KeyValue {
	return TrackingIDKey.String(v)
}
