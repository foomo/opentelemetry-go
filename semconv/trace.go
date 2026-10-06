package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	// TraceIDKey is the key for trace.id.
	TraceIDKey = attribute.Key("trace.id")
	// SpanIDKey is the key for span.id.
	SpanIDKey = attribute.Key("span.id")
	// SamplingPriorityKey is the key for sampling.priority.
	SamplingPriorityKey = attribute.Key("sampling.priority")
)

// TraceID returns a new attribute.KeyValue for trace.id.
func TraceID(v string) attribute.KeyValue {
	return TraceIDKey.String(v)
}

// SpanID returns a new attribute.KeyValue for span.id.
func SpanID(v string) attribute.KeyValue {
	return SpanIDKey.String(v)
}

// SamplingPriority returns a new attribute.KeyValue for sampling.priority.
// A value > 0 asks samplers to keep the trace, 0 to drop it. Set it on any
// span and match it with a tail sampling numeric_attribute policy.
func SamplingPriority(v int) attribute.KeyValue {
	return SamplingPriorityKey.Int(v)
}
