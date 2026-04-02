package glossymetric_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossymetric"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	res = resource.NewSchemaless(
		semconv.ServiceName("test-service"),
	)
)

func mockMetrics() *metricdata.ResourceMetrics {
	return &metricdata.ResourceMetrics{
		Resource: res,
		ScopeMetrics: []metricdata.ScopeMetrics{
			{
				Scope: instrumentation.Scope{Name: "test-scope", Version: "0.0.1"},
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
						Name:        "cpu.time",
						Description: "Accumulated CPU time",
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
						Name:        "memory.usage",
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
						Description: "CPU temperature",
						Unit:        "cel",
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
					{
						Name:        "latency",
						Description: "Request latency",
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
									Min:          metricdata.NewExtrema[float64](1.2),
									Max:          metricdata.NewExtrema[float64](9.8),
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestExportMetrics(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossymetric.New(glossymetric.WithWriter(&buf))
	if err != nil {
		t.Fatal(err)
	}

	if err := exp.Export(context.Background(), mockMetrics()); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	for _, want := range []string{
		"METRICS",
		"Scope: test-scope",
		"Metric: requests",
		"Value: 5",
		"Metric: cpu.time",
		"Value: 0.500",
		"Metric: memory.usage",
		"Value: 100",
		"Metric: temperature",
		"Value: 32.400",
		"Metric: latency",
		"Histogram:",
		"min=1.200",
		"max=9.800",
		"count=10",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestExportMetricsWithoutHistograms(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossymetric.New(
		glossymetric.WithWriter(&buf),
		glossymetric.WithoutHistograms(),
	)
	if err != nil {
		t.Fatal(err)
	}

	if err := exp.Export(context.Background(), mockMetrics()); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "Histogram:") {
		t.Error("expected no histogram detail when WithoutHistograms is set")
	}

	if !strings.Contains(out, "Metric: latency") {
		t.Error("expected latency metric header to still be present")
	}
}

func TestExportMetricsAttributes(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossymetric.New(glossymetric.WithWriter(&buf))
	if err != nil {
		t.Fatal(err)
	}

	if err := exp.Export(context.Background(), mockMetrics()); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "server=central") {
		t.Errorf("expected attribute in output, got:\n%s", out)
	}
}

func TestExportContextCanceled(t *testing.T) {
	exp, err := glossymetric.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := exp.Export(ctx, mockMetrics()); err == nil {
		t.Error("expected error for canceled context")
	}
}

func TestShutdownAndForceFlush(t *testing.T) {
	exp, err := glossymetric.New()
	if err != nil {
		t.Fatal(err)
	}

	if err := exp.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}

	if err := exp.ForceFlush(context.Background()); err != nil {
		t.Fatal(err)
	}
}
