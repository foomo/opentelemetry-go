// Package natsconv provides types and functionality for NATS-specific metrics
// following the OpenTelemetry semantic conventions pattern in the "nats" and
// "messaging.nats" namespaces. It complements the generic messaging semantic
// conventions with NATS-specific instruments (JetStream consumer state, core
// connection health).
//
// Each instrument type wraps an OpenTelemetry metric instrument and exposes its
// semantic convention Name, Unit and Description. Constructors accept a nil
// [go.opentelemetry.io/otel/metric.Meter] and return a no-op instrument in that
// case.
//
// # Usage
//
//	disconnects, err := natsconv.NewClientDisconnects(meter)
//	if err != nil {
//		return err
//	}
//	disconnects.Add(ctx, 1, nc.ConnectedUrl(),
//		disconnects.AttrClientName("my-service"),
//	)
package natsconv
