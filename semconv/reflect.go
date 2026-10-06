package semconv

import (
	"reflect"

	"go.opentelemetry.io/otel/attribute"
)

const (
	// ReflectTypeKey is the key for reflect.type.
	ReflectTypeKey = attribute.Key("reflect.type")
)

// ReflectType returns a new attribute.KeyValue for reflect.type set to the
// dynamic type name of v, such as "*errors.errorString". It panics if v is nil.
func ReflectType(v any) attribute.KeyValue {
	return ReflectTypeKey.String(reflect.TypeOf(v).String())
}
