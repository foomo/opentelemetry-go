package semconv

import (
	"go.opentelemetry.io/otel/attribute"
)

const (
	// DistributedLockNameKey is the key for distributed_lock.name.
	DistributedLockNameKey = attribute.Key("distributed_lock.name")
	// DistributedLockSystemKey is the key for distributed_lock.system.
	DistributedLockSystemKey = attribute.Key("distributed_lock.system")
	// DistributedLockAcquireResultKey is the key for distributed_lock.acquire.result.
	DistributedLockAcquireResultKey = attribute.Key("distributed_lock.acquire.result")
	// DistributedLockOwnerIDKey is the key for distributed_lock.owner.id.
	DistributedLockOwnerIDKey = attribute.Key("distributed_lock.owner.id")
)

// DistributedLockName returns a new attribute.KeyValue for distributed_lock.name.
// Use the logical lock name, not the per-resource key, to keep cardinality low.
func DistributedLockName(v string) attribute.KeyValue {
	return DistributedLockNameKey.String(v)
}

// DistributedLockSystem returns a new attribute.KeyValue for distributed_lock.system.
func DistributedLockSystem(v string) attribute.KeyValue {
	return DistributedLockSystemKey.String(v)
}

// DistributedLockAcquireResult returns a new attribute.KeyValue for distributed_lock.acquire.result.
func DistributedLockAcquireResult(v string) attribute.KeyValue {
	return DistributedLockAcquireResultKey.String(v)
}

// DistributedLockOwnerID returns a new attribute.KeyValue for distributed_lock.owner.id.
// Values are unbounded: use on spans and logs only, never on metrics.
func DistributedLockOwnerID(v string) attribute.KeyValue {
	return DistributedLockOwnerIDKey.String(v)
}
