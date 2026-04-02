package glossytrace_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

var (
	traceID  = trace.TraceID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
	rootSpan = trace.SpanID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	childID  = trace.SpanID{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}
	now      = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
)

func newStubSpan(name string, spanID, parentID trace.SpanID, start, end time.Time, attrs []attribute.KeyValue) sdktrace.ReadOnlySpan {
	stub := tracetest.SpanStub{
		Name: name,
		SpanContext: trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		}),
		StartTime: start,
		EndTime:   end,
		Status: sdktrace.Status{
			Code: codes.Ok,
		},
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

func TestExportSpansTree(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	spans := []sdktrace.ReadOnlySpan{
		newStubSpan("root", rootSpan, trace.SpanID{}, now, now.Add(100*time.Millisecond), nil),
		newStubSpan("child", childID, rootSpan, now.Add(10*time.Millisecond), now.Add(50*time.Millisecond), nil),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "TRACE") {
		t.Error("expected TRACE header in output")
	}

	if !strings.Contains(out, "root") {
		t.Error("expected root span name in output")
	}

	if !strings.Contains(out, "child") {
		t.Error("expected child span name in output")
	}

	if !strings.Contains(out, "100.00 ms") {
		t.Errorf("expected root duration in output, got:\n%s", out)
	}

	if !strings.Contains(out, "40.00 ms") {
		t.Errorf("expected child duration in output, got:\n%s", out)
	}
}

func TestExportSpansWithFlamegraph(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithFlamegraph(),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	spans := []sdktrace.ReadOnlySpan{
		newStubSpan("root", rootSpan, trace.SpanID{}, now, now.Add(100*time.Millisecond), nil),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "Flamegraph") {
		t.Error("expected Flamegraph header in output")
	}

	if !strings.Contains(out, "█") {
		t.Error("expected flamegraph bar in output")
	}
}

func TestExportSpansWithAttributes(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithSpanAttributes(),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	spans := []sdktrace.ReadOnlySpan{
		newStubSpan("root", rootSpan, trace.SpanID{}, now, now.Add(100*time.Millisecond),
			[]attribute.KeyValue{attribute.String("http.method", "GET")}),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "http.method") {
		t.Error("expected attribute key in output")
	}

	if !strings.Contains(out, "GET") {
		t.Error("expected attribute value in output")
	}
}

func TestExportSpansMinDuration(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
		glossytrace.WithMinDuration(60*time.Millisecond),
	)
	if err != nil {
		t.Fatal(err)
	}

	spans := []sdktrace.ReadOnlySpan{
		newStubSpan("fast", rootSpan, trace.SpanID{}, now, now.Add(10*time.Millisecond), nil),
		newStubSpan("slow", childID, trace.SpanID{}, now, now.Add(100*time.Millisecond), nil),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if strings.Contains(out, "fast") {
		t.Error("expected fast span to be filtered out")
	}

	if !strings.Contains(out, "slow") {
		t.Error("expected slow span to be present")
	}
}

func TestExportSpansWithTimestamps(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
	)
	if err != nil {
		t.Fatal(err)
	}

	spans := []sdktrace.ReadOnlySpan{
		newStubSpan("root", rootSpan, trace.SpanID{}, now, now.Add(100*time.Millisecond), nil),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	if !strings.Contains(out, "00:00:00.000") {
		t.Errorf("expected timestamp in output, got:\n%s", out)
	}
}

func TestExportSpansContextCanceled(t *testing.T) {
	exp, err := glossytrace.New()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := exp.ExportSpans(ctx, nil); err == nil {
		t.Error("expected error for canceled context")
	}
}

func TestShutdown(t *testing.T) {
	exp, err := glossytrace.New()
	if err != nil {
		t.Fatal(err)
	}

	if err := exp.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
}
