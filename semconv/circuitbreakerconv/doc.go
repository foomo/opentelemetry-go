// Package circuitbreakerconv provides types and functionality for circuit
// breaker metrics following the OpenTelemetry semantic conventions pattern in
// the "circuit_breaker" namespace. It is library-agnostic: wire the instruments
// into the state-change and rejection hooks of the breaker implementation in use.
//
// Each instrument type wraps an OpenTelemetry metric instrument and exposes its
// semantic convention Name, Unit and Description. Constructors accept a nil
// [go.opentelemetry.io/otel/metric.Meter] and return a no-op instrument in that
// case.
//
// # Usage
//
//	stateChanges, err := circuitbreakerconv.NewStateChanges(meter)
//	if err != nil {
//		return err
//	}
//	stateChanges.Add(ctx, 1, "payment-api", circuitbreakerconv.StateOpen)
package circuitbreakerconv
