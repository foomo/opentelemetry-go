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
	// └─ HTTP GET /api/users (150.00 ms) (0102030405060708090a0b0c0d0e0f10:0102030405060708)
	//   └─ db.Query (40.00 ms) (0102030405060708090a0b0c0d0e0f10:1112131415161718)
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
	// └─ HTTP GET /api/users (150.00 ms) (0102030405060708090a0b0c0d0e0f10:0102030405060708)
	//   └─ db.Query (40.00 ms) (0102030405060708090a0b0c0d0e0f10:1112131415161718)
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
	// └─ HTTP GET /api/users (150.00 ms) (0102030405060708090a0b0c0d0e0f10:0102030405060708)
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
	// └─ HTTP GET /api/users (150.00 ms) (0102030405060708090a0b0c0d0e0f10:0102030405060708)
}

func Example_links() {
	_ = os.Setenv("NO_COLOR", "1")

	exp, _ := glossytrace.New(
		glossytrace.WithWriter(os.Stdout),
		glossytrace.WithoutTimestamps(),
	)

	trace2 := trace.TraceID{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30}
	spanB := trace.SpanID{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28}
	spanC := trace.SpanID{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38}

	linkedSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace2,
		SpanID:  spanB,
	})

	missingSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff, 0x00},
		SpanID:  trace.SpanID{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22},
	})

	spans := []sdktrace.ReadOnlySpan{
		newExampleSpanWithLinks("HTTP POST /dispatch", rootSpan, trace.SpanID{}, traceID,
			now, now.Add(200*time.Millisecond), nil,
			[]sdktrace.Link{{SpanContext: linkedSC}, {SpanContext: missingSC}}),
		newExampleSpanWithLinks("async.process", spanB, trace.SpanID{}, trace2,
			now.Add(50*time.Millisecond), now.Add(150*time.Millisecond), nil, nil),
		newExampleSpanWithLinks("db.Insert", spanC, spanB, trace2,
			now.Add(60*time.Millisecond), now.Add(90*time.Millisecond), nil, nil),
	}

	_ = exp.ExportSpans(context.Background(), spans)

	// Output:
	// === TRACE 0102030405060708090a0b0c0d0e0f10 ===
	// └─ HTTP POST /dispatch (200.00 ms) (0102030405060708090a0b0c0d0e0f10:0102030405060708)
	//   Links:
	//   ⤡ link: async.process (2122232425262728292a2b2c2d2e2f30:2122232425262728)
	//   ⤡ link: f1f2f3f4f5f6f7f8f9fafbfcfdfeff00:aabbccddeeff1122
	//
	// === TRACE 2122232425262728292a2b2c2d2e2f30 ===
	// └─ async.process (100.00 ms) (2122232425262728292a2b2c2d2e2f30:2122232425262728)
	//   └─ db.Insert (30.00 ms) (2122232425262728292a2b2c2d2e2f30:3132333435363738)
}

func newExampleSpanWithLinks(name string, spanID, parentID trace.SpanID, tid trace.TraceID, start, end time.Time, attrs []attribute.KeyValue, links []sdktrace.Link) sdktrace.ReadOnlySpan {
	stub := tracetest.SpanStub{
		Name: name,
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: tid,
			SpanID:  spanID,
		}),
		StartTime:  start,
		EndTime:    end,
		Status:     sdktrace.Status{Code: codes.Ok},
		Attributes: attrs,
		Links:      links,
	}
	if parentID.IsValid() {
		stub.Parent = trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: tid,
			SpanID:  parentID,
		})
	}

	return stub.Snapshot()
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
