package glossymetric_test

import (
	"context"
	"os"
	"time"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func Example() {
	os.Setenv("NO_COLOR", "1")

	exporter, err := glossymetric.New(glossymetric.WithWriter(os.Stdout))
	if err != nil {
		panic(err)
	}

	// Register the exporter with an SDK via a periodic reader.
	sdk := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)

	mockData := metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "example", Version: "0.0.1"},
				Metrics: []metricdata.Metrics{
					{
						Name:        "requests",
						Description: "Number of requests received",
						Unit:        "1",
						Data: metricdata.Sum[int64]{
							IsMonotonic: true,
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("server", "central")),
									StartTime:  now,
									Time:       now.Add(1 * time.Second),
									Value:      5,
								},
							},
						},
					},
					{
						Name:        "system.cpu.time",
						Description: "Accumulated CPU time spent",
						Unit:        "s",
						Data: metricdata.Sum[float64]{
							IsMonotonic: true,
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(attribute.String("state", "user")),
									StartTime:  now,
									Time:       now.Add(1 * time.Second),
									Value:      0.5,
								},
							},
						},
					},
					{
						Name:        "requests.size",
						Description: "Size of received requests",
						Unit:        "kb",
						Data: metricdata.Histogram[int64]{
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.HistogramDataPoint[int64]{
								{
									Attributes:   attribute.NewSet(attribute.String("server", "central")),
									StartTime:    now,
									Time:         now.Add(1 * time.Second),
									Count:        10,
									Bounds:       []float64{1, 5, 10},
									BucketCounts: []uint64{1, 3, 6, 0},
									Sum:          128,
									Min:          metricdata.NewExtrema[int64](3),
									Max:          metricdata.NewExtrema[int64](30),
								},
							},
						},
					},
					{
						Name:        "latency",
						Description: "Time spend processing received requests",
						Unit:        "ms",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
								{
									Attributes:   attribute.NewSet(attribute.String("server", "central")),
									StartTime:    now,
									Time:         now.Add(1 * time.Second),
									Count:        10,
									Bounds:       []float64{1, 5, 10},
									BucketCounts: []uint64{1, 3, 6, 0},
									Sum:          57,
								},
							},
						},
					},
					{
						Name:        "system.memory.usage",
						Description: "Memory usage",
						Unit:        "By",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("state", "used")),
									Time:       now.Add(1 * time.Second),
									Value:      100,
								},
							},
						},
					},
					{
						Name:        "temperature",
						Description: "CPU global temperature",
						Unit:        "cel(1 K)",
						Data: metricdata.Gauge[float64]{
							DataPoints: []metricdata.DataPoint[float64]{
								{
									Attributes: attribute.NewSet(attribute.String("server", "central")),
									Time:       now.Add(1 * time.Second),
									Value:      32.4,
								},
							},
						},
					},
				},
			},
		},
	}

	ctx := context.Background()
	_ = exporter.Export(ctx, &mockData)

	// Ensure the periodic reader is cleaned up by shutting down the sdk.
	_ = sdk.Shutdown(ctx)

	// Output:
	// METRICS
	//
	// Scope: example
	// ──────────────────────────────────────────────────
	// Metric: requests
	// Description: Number of requests received
	// Unit: 1
	//   attrs: server=central
	// Value: 5
	// ──────────────────────────────────────────────────
	// Metric: system.cpu.time
	// Description: Accumulated CPU time spent
	// Unit: s
	//   attrs: state=user
	// Value: 0.500
	// ──────────────────────────────────────────────────
	// Metric: requests.size
	// Description: Size of received requests
	// Unit: kb
	//   attrs: server=central
	// Histogram: min=3 max=30 avg=12.800 count=10
	// ──────────────────────────────────────────────────
	// Metric: latency
	// Description: Time spend processing received requests
	// Unit: ms
	//   attrs: server=central
	// Histogram: min=0.000 max=0.000 avg=5.700 count=10
	// ──────────────────────────────────────────────────
	// Metric: system.memory.usage
	// Description: Memory usage
	// Unit: By
	//   attrs: state=used
	// Value: 100
	// ──────────────────────────────────────────────────
	// Metric: temperature
	// Description: CPU global temperature
	// Unit: cel(1 K)
	//   attrs: server=central
	// Value: 32.400
	// ──────────────────────────────────────────────────
}

func Example_withoutHistograms() {
	os.Setenv("NO_COLOR", "1")

	exporter, _ := glossymetric.New(
		glossymetric.WithWriter(os.Stdout),
		glossymetric.WithoutHistograms(),
	)

	mockData := metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "example"},
				Metrics: []metricdata.Metrics{
					{
						Name:        "requests",
						Description: "Number of requests",
						Unit:        "1",
						Data: metricdata.Sum[int64]{
							IsMonotonic: true,
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("method", "GET")),
									StartTime:  now,
									Time:       now.Add(1 * time.Second),
									Value:      42,
								},
							},
						},
					},
					{
						Name:        "latency",
						Description: "Request latency",
						Unit:        "ms",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.DeltaTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
								{
									Attributes:   attribute.NewSet(attribute.String("method", "GET")),
									StartTime:    now,
									Time:         now.Add(1 * time.Second),
									Count:        10,
									Bounds:       []float64{1, 5, 10},
									BucketCounts: []uint64{1, 3, 6, 0},
									Sum:          57,
									Min:          metricdata.NewExtrema[float64](0.5),
									Max:          metricdata.NewExtrema[float64](9.8),
								},
							},
						},
					},
				},
			},
		},
	}

	_ = exporter.Export(context.Background(), &mockData)

	// Output:
	// METRICS
	//
	// Scope: example
	// ──────────────────────────────────────────────────
	// Metric: requests
	// Description: Number of requests
	// Unit: 1
	//   attrs: method=GET
	// Value: 42
	// ──────────────────────────────────────────────────
	// Metric: latency
	// Description: Request latency
	// Unit: ms
	// ──────────────────────────────────────────────────
}
