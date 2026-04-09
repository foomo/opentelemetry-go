package glossytrace_test

import (
	"context"
	"os"
	"time"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"
)

func Example() {
	_ = os.Setenv("NO_COLOR", "1")

	ctx := context.Background()

	exporter, _ := glossytrace.New(
		glossytrace.WithWriter(os.Stdout),
		glossytrace.WithoutTimestamps(),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("stdout-example"),
			semconv.ServiceVersion("0.0.1"),
		)),
	)
	otel.SetTracerProvider(tracerProvider)

	spans := []sdktrace.ReadOnlySpan{
		newExampleSpan("HTTP GET /api/users", rootSpan, trace.SpanID{}, now, now.Add(150*time.Millisecond), nil),
		newExampleSpan("db.Query", childID, rootSpan, now.Add(10*time.Millisecond), now.Add(50*time.Millisecond), nil),
	}

	_ = exporter.ExportSpans(ctx, spans)

	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			panic(err)
		}
	}()

	// Output:
	// === TRACE 0102030405060708090a0b0c0d0e0f10 ===
	// └─ HTTP GET /api/users (150.00 ms)
	//   └─ db.Query (40.00 ms)
}

func Example_flamegraph() {
	_ = os.Setenv("NO_COLOR", "1")

	exp, _ := glossytrace.New(
		glossytrace.WithWriter(os.Stdout),
		glossytrace.WithoutTimestamps(),
		glossytrace.WithFlamegraph(),
	)

	spans := []sdktrace.ReadOnlySpan{
		newExampleSpan("HTTP GET /api/users", rootSpan, trace.SpanID{}, now, now.Add(150*time.Millisecond), nil),
		newExampleSpan("db.Query", childID, rootSpan, now.Add(10*time.Millisecond), now.Add(50*time.Millisecond), nil),
	}

	_ = exp.ExportSpans(context.Background(), spans)

	// Output:
	// === TRACE 0102030405060708090a0b0c0d0e0f10 ===
	// └─ HTTP GET /api/users (150.00 ms)
	//   └─ db.Query (40.00 ms)
	//
	// Nested Flamegraph:
	// HTTP GET /api/users ██████████████████████████████ 150.00 ms
	//   db.Query        ████████ 40.00 ms
}

func Example_attributes() {
	_ = os.Setenv("NO_COLOR", "1")

	exp, _ := glossytrace.New(
		glossytrace.WithWriter(os.Stdout),
		glossytrace.WithoutTimestamps(),
		glossytrace.WithSpanAttributes(),
	)

	spans := []sdktrace.ReadOnlySpan{
		newExampleSpan("HTTP GET /api/users", rootSpan, trace.SpanID{}, now, now.Add(150*time.Millisecond),
			[]attribute.KeyValue{attribute.String("http.method", "GET"), attribute.String("http.target", "/api/users")}),
	}

	_ = exp.ExportSpans(context.Background(), spans)

	// Output:
	// === TRACE 0102030405060708090a0b0c0d0e0f10 ===
	// └─ HTTP GET /api/users (150.00 ms)
	//     http.method=GET
	//     http.target=/api/users
}

func Example_minDuration() {
	_ = os.Setenv("NO_COLOR", "1")

	exp, _ := glossytrace.New(
		glossytrace.WithWriter(os.Stdout),
		glossytrace.WithoutTimestamps(),
		glossytrace.WithMinDuration(50*time.Millisecond),
	)

	spans := []sdktrace.ReadOnlySpan{
		newExampleSpan("HTTP GET /api/users", rootSpan, trace.SpanID{}, now, now.Add(150*time.Millisecond), nil),
		newExampleSpan("db.Query", childID, rootSpan, now.Add(10*time.Millisecond), now.Add(30*time.Millisecond), nil),
	}

	_ = exp.ExportSpans(context.Background(), spans)

	// Output:
	// === TRACE 0102030405060708090a0b0c0d0e0f10 ===
	// └─ HTTP GET /api/users (150.00 ms)
}

func newExampleSpan(name string, spanID, parentID trace.SpanID, start, end time.Time, attrs []attribute.KeyValue) sdktrace.ReadOnlySpan {
	stub := tracetest.SpanStub{
		Name: name,
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		}),
		StartTime:  start,
		EndTime:    end,
		Status:     sdktrace.Status{Code: codes.Ok},
		Attributes: attrs,
	}
	if parentID.IsValid() {
		stub.Parent = trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  parentID,
		})
	}

	return stub.Snapshot()
}
