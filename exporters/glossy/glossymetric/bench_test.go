package glossymetric_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func makeMetrics(nScopes, nMetrics int) *metricdata.ResourceMetrics {
	rm := &metricdata.ResourceMetrics{
		Resource: res,
	}

	for i := range nScopes {
		sm := metricdata.ScopeMetrics{
			Scope: instrumentation.Scope{Name: "scope-" + string(rune('A'+i))},
		}

		for j := range nMetrics {
			_ = j

			sm.Metrics = append(sm.Metrics,
				metricdata.Metrics{
					Name:        "counter",
					Description: "A counter",
					Unit:        "1",
					Data: metricdata.Sum[int64]{
						IsMonotonic: true,
						Temporality: metricdata.DeltaTemporality,
						DataPoints: []metricdata.DataPoint[int64]{
							{
								Attributes: attribute.NewSet(attribute.String("key", "value")),
								StartTime:  now,
								Time:       now.Add(time.Second),
								Value:      42,
							},
						},
					},
				},
				metricdata.Metrics{
					Name:        "latency",
					Description: "Request latency",
					Unit:        "ms",
					Data: metricdata.Histogram[float64]{
						Temporality: metricdata.DeltaTemporality,
						DataPoints: []metricdata.HistogramDataPoint[float64]{
							{
								Attributes:   attribute.NewSet(attribute.String("key", "value")),
								StartTime:    now,
								Time:         now.Add(time.Second),
								Count:        100,
								Bounds:       []float64{1, 5, 10, 50, 100},
								BucketCounts: []uint64{10, 20, 30, 25, 10, 5},
								Sum:          1234.5,
								Min:          metricdata.NewExtrema[float64](0.5),
								Max:          metricdata.NewExtrema[float64](99.9),
							},
						},
					},
				},
			)
		}

		rm.ScopeMetrics = append(rm.ScopeMetrics, sm)
	}

	return rm
}

func BenchmarkExport(b *testing.B) {
	b.Setenv("NO_COLOR", "1")

	for _, tc := range []struct {
		name            string
		scopes, metrics int
	}{
		{"1x1", 1, 1},
		{"1x10", 1, 10},
		{"5x10", 5, 10},
	} {
		rm := makeMetrics(tc.scopes, tc.metrics)

		b.Run(tc.name, func(b *testing.B) {
			exp, err := glossymetric.New(glossymetric.WithWriter(io.Discard))
			if err != nil {
				b.Fatal(err)
			}

			ctx := context.Background()

			b.ResetTimer()
			b.ReportAllocs()

			for range b.N {
				if err := exp.Export(ctx, rm); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
