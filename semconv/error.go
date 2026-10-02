package semconv

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

// ErrorType returns a new attribute.KeyValue for error.type. It extends the
// upstream semconv.ErrorType by mapping context.Canceled and
// context.DeadlineExceeded to stable values instead of their reflected type
// names.
func ErrorType(err error) attribute.KeyValue {
	switch {
	case errors.Is(err, context.Canceled):
		return semconv.ErrorTypeKey.String("context.Canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return semconv.ErrorTypeKey.String("context.DeadlineExceeded")
	default:
		return semconv.ErrorType(err)
	}
}
