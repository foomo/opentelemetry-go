// Package semconv provides attribute keys and constructors for semantic
// conventions that are not covered by the upstream
// [go.opentelemetry.io/otel/semconv] packages.
//
// Each convention is exposed as an [go.opentelemetry.io/otel/attribute.Key]
// constant named <Name>Key and a constructor <Name> that returns an
// [go.opentelemetry.io/otel/attribute.KeyValue]:
//
//	span.SetAttributes(
//		semconv.TrackingID(id),
//		semconv.ErrorType(err),
//	)
//
// Instrument wrappers for library-specific metrics live in subpackages such
// as [github.com/foomo/opentelemetry-go/semconv/natsconv].
package semconv
