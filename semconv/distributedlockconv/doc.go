// Package distributedlockconv provides types and functionality for distributed
// lock metrics following the OpenTelemetry semantic conventions pattern in the
// "distributed_lock" namespace. It is backend-agnostic: the same instruments
// cover locks backed by Redis, MongoDB, NATS, etcd, PostgreSQL, etc.
//
// Each instrument type wraps an OpenTelemetry metric instrument and exposes its
// semantic convention Name, Unit and Description. Constructors accept a nil
// [go.opentelemetry.io/otel/metric.Meter] and return a no-op instrument in that
// case.
//
// # Usage
//
//	acquire, err := distributedlockconv.NewAcquireDuration(meter)
//	if err != nil {
//		return err
//	}
//	start := time.Now()
//	// ... acquire lock ...
//	acquire.Record(ctx, time.Since(start).Seconds(), "order-sync",
//		distributedlockconv.SystemRedis,
//		distributedlockconv.AcquireResultAcquired,
//	)
package distributedlockconv
