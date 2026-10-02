package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	// HTTPXRequestIDKey is the key for http.request.id.
	//
	// Deprecated: http.request.* is an upstream namespace. Use
	// semconv.HTTPRequestHeader("x-request-id", v) from
	// go.opentelemetry.io/otel/semconv instead.
	HTTPXRequestIDKey = attribute.Key("http.request.id")
	// HTTPXRequestRefererKey is the key for http.request.referer.
	//
	// Deprecated: Use semconv.HTTPRequestHeader("referer", v) from
	// go.opentelemetry.io/otel/semconv instead.
	HTTPXRequestRefererKey = attribute.Key("http.request.referer")
)

// HTTPXRequestID returns a new attribute.KeyValue for http.request.id.
//
// Deprecated: http.request.* is an upstream namespace. Use
// semconv.HTTPRequestHeader("x-request-id", v) from
// go.opentelemetry.io/otel/semconv instead.
func HTTPXRequestID(v string) attribute.KeyValue {
	return HTTPXRequestIDKey.String(v)
}

// HTTPXRequestReferer returns a new attribute.KeyValue for http.request.referer.
//
// Deprecated: Use semconv.HTTPRequestHeader("referer", v) from
// go.opentelemetry.io/otel/semconv instead.
func HTTPXRequestReferer(v string) attribute.KeyValue {
	return HTTPXRequestRefererKey.String(v)
}
