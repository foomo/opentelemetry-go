package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	// CircuitBreakerNameKey is the key for circuit_breaker.name.
	CircuitBreakerNameKey = attribute.Key("circuit_breaker.name")
	// CircuitBreakerStateKey is the key for circuit_breaker.state.
	CircuitBreakerStateKey = attribute.Key("circuit_breaker.state")
)

// CircuitBreakerName returns a new attribute.KeyValue for circuit_breaker.name.
func CircuitBreakerName(v string) attribute.KeyValue {
	return CircuitBreakerNameKey.String(v)
}

// CircuitBreakerState returns a new attribute.KeyValue for circuit_breaker.state.
func CircuitBreakerState(v string) attribute.KeyValue {
	return CircuitBreakerStateKey.String(v)
}
