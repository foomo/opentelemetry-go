package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	// TraceIDKey is the key for trace.id.
	TraceIDKey = attribute.Key("trace.id")
	// SpanIDKey is the key for span.id.
	SpanIDKey = attribute.Key("span.id")
)

// TraceID returns a new attribute.KeyValue for trace.id.
func TraceID(v string) attribute.KeyValue {
	return TraceIDKey.String(v)
}

// SpanID returns a new attribute.KeyValue for span.id.
func SpanID(v string) attribute.KeyValue {
	return SpanIDKey.String(v)
}
