package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	ProfileNameKey = attribute.Key("profile.name")
)

// ProfileName returns a new attribute.KeyValue for profile.name.
func ProfileName(v string) attribute.KeyValue {
	return ProfileNameKey.String(v)
}
