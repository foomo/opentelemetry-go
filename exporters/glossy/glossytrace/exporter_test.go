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
	traceID2 = trace.TraceID{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30}
	rootSpan = trace.SpanID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
	childID  = trace.SpanID{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18}
	spanID3  = trace.SpanID{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28}
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

func newStubSpanWithLinks(name string, spanID, parentID trace.SpanID, tid trace.TraceID, start, end time.Time, attrs []attribute.KeyValue, links []sdktrace.Link) sdktrace.ReadOnlySpan {
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

func TestExportSpansWithLinksInBatch(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	linkedSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  childID,
	})
	spans := []sdktrace.ReadOnlySpan{
		newStubSpanWithLinks("spanA", rootSpan, trace.SpanID{}, traceID, now, now.Add(100*time.Millisecond), nil,
			[]sdktrace.Link{{SpanContext: linkedSC}}),
		newStubSpanWithLinks("spanB", childID, trace.SpanID{}, traceID, now.Add(10*time.Millisecond), now.Add(50*time.Millisecond), nil, nil),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	// In-batch links show span name with reference.
	if !strings.Contains(out, "⤡ link:") {
		t.Errorf("expected link prefix in output, got:\n%s", out)
	}

	if !strings.Contains(out, "spanB") {
		t.Errorf("expected linked span name in output, got:\n%s", out)
	}

	if !strings.Contains(out, "Links:") {
		t.Errorf("expected Links: section in output, got:\n%s", out)
	}

	// Should include the TraceID:SpanID reference.
	if !strings.Contains(out, traceID.String()) {
		t.Errorf("expected trace ID reference in output, got:\n%s", out)
	}
}

func TestExportSpansWithLinksNotInBatch(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	missingSpanID := trace.SpanID{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22}
	missingTraceID := trace.TraceID{0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7, 0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff, 0x00}
	linkedSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: missingTraceID,
		SpanID:  missingSpanID,
	})

	spans := []sdktrace.ReadOnlySpan{
		newStubSpanWithLinks("spanA", rootSpan, trace.SpanID{}, traceID, now, now.Add(100*time.Millisecond), nil,
			[]sdktrace.Link{{SpanContext: linkedSC}}),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	if !strings.Contains(out, "⤡ link:") {
		t.Errorf("expected link prefix in output, got:\n%s", out)
	}

	if !strings.Contains(out, missingTraceID.String()) {
		t.Errorf("expected missing trace ID in output, got:\n%s", out)
	}

	if !strings.Contains(out, missingSpanID.String()) {
		t.Errorf("expected missing span ID in output, got:\n%s", out)
	}
}

func TestExportSpansWithLinksAttributes(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
		glossytrace.WithSpanAttributes(),
	)
	if err != nil {
		t.Fatal(err)
	}

	missingSpanID := trace.SpanID{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22}
	linkedSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID2,
		SpanID:  missingSpanID,
	})

	spans := []sdktrace.ReadOnlySpan{
		newStubSpanWithLinks("spanA", rootSpan, trace.SpanID{}, traceID, now, now.Add(100*time.Millisecond), nil,
			[]sdktrace.Link{{
				SpanContext: linkedSC,
				Attributes:  []attribute.KeyValue{attribute.String("reason", "retry")},
			}}),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	if !strings.Contains(out, "reason") {
		t.Errorf("expected link attribute key in output, got:\n%s", out)
	}

	if !strings.Contains(out, "retry") {
		t.Errorf("expected link attribute value in output, got:\n%s", out)
	}
}

func TestExportSpansWithLinksCrossTrace(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	linkedSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID2,
		SpanID:  spanID3,
	})

	spans := []sdktrace.ReadOnlySpan{
		newStubSpanWithLinks("producer", rootSpan, trace.SpanID{}, traceID, now, now.Add(100*time.Millisecond), nil,
			[]sdktrace.Link{{SpanContext: linkedSC}}),
		newStubSpanWithLinks("consumer", spanID3, trace.SpanID{}, traceID2, now.Add(10*time.Millisecond), now.Add(60*time.Millisecond), nil, nil),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	// Both traces should appear (no collapsing).
	if strings.Count(out, "=== TRACE") != 2 {
		t.Errorf("expected 2 trace headers, got:\n%s", out)
	}

	// Link reference should show consumer's name (in-batch).
	if !strings.Contains(out, "⤡ link:") {
		t.Errorf("expected link prefix in output, got:\n%s", out)
	}

	if !strings.Contains(out, "consumer") {
		t.Errorf("expected consumer span name in link reference, got:\n%s", out)
	}
}

func TestExportSpansWithLinksAndFlamegraph(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
		glossytrace.WithFlamegraph(),
	)
	if err != nil {
		t.Fatal(err)
	}

	linkedSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  childID,
	})

	spans := []sdktrace.ReadOnlySpan{
		newStubSpanWithLinks("spanA", rootSpan, trace.SpanID{}, traceID, now, now.Add(100*time.Millisecond), nil,
			[]sdktrace.Link{{SpanContext: linkedSC}}),
		newStubSpanWithLinks("spanB", childID, trace.SpanID{}, traceID, now.Add(10*time.Millisecond), now.Add(50*time.Millisecond), nil, nil),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()

	if !strings.Contains(out, "Flamegraph") {
		t.Errorf("expected Flamegraph header in output, got:\n%s", out)
	}

	// Links section should appear in both tree and flamegraph.
	if count := strings.Count(out, "Links:"); count < 2 {
		t.Errorf("expected Links: in both tree and flamegraph, got %d in:\n%s", count, out)
	}
}

func TestExportSpansWithLinksIntegration(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Hour)),
	)
	tracer := tp.Tracer("test")

	_, producer := tracer.Start(context.Background(), "producer")
	producerSC := producer.SpanContext()
	producer.End()

	_, consumer := tracer.Start(context.Background(), "consumer",
		trace.WithLinks(trace.Link{SpanContext: producerSC}))
	consumer.End()

	if err := tp.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	t.Logf("output:\n%s", out)

	if !strings.Contains(out, "⤡ link:") {
		t.Errorf("expected link prefix in output, got:\n%s", out)
	}

	// Both producer and consumer should appear, with link reference showing producer's name.
	if !strings.Contains(out, "producer") {
		t.Errorf("expected producer span in output, got:\n%s", out)
	}

	if !strings.Contains(out, "consumer") {
		t.Errorf("expected consumer span in output, got:\n%s", out)
	}
}

func TestExportSpansWithLinksToParent(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	var buf bytes.Buffer

	exp, err := glossytrace.New(
		glossytrace.WithWriter(&buf),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Child span links to its own parent (common pattern with SpanContextFromContext).
	parentSC := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  rootSpan,
	})

	spans := []sdktrace.ReadOnlySpan{
		newStubSpanWithLinks("parentOp", rootSpan, trace.SpanID{}, traceID, now, now.Add(100*time.Millisecond), nil, nil),
		newStubSpanWithLinks("childOp", childID, rootSpan, traceID, now.Add(10*time.Millisecond), now.Add(50*time.Millisecond), nil,
			[]sdktrace.Link{{SpanContext: parentSC}}),
	}

	if err := exp.ExportSpans(context.Background(), spans); err != nil {
		t.Fatal(err)
	}

	out := buf.String()
	t.Logf("output:\n%s", out)

	// Parent should appear normally in the tree.
	if !strings.Contains(out, "└─ parentOp") {
		t.Errorf("expected parent span in tree, got:\n%s", out)
	}

	// Child should appear as a child of parent.
	if !strings.Contains(out, "└─ childOp") {
		t.Errorf("expected child span in tree, got:\n%s", out)
	}

	// Child's link to parent should show with parent's name.
	if !strings.Contains(out, "Links:") {
		t.Errorf("expected Links: section under child, got:\n%s", out)
	}

	if !strings.Contains(out, "⤡ link:") {
		t.Errorf("expected link prefix in output, got:\n%s", out)
	}

	// Link reference should include the parent's name (resolved from spanIndex).
	if !strings.Contains(out, "parentOp") {
		t.Errorf("expected parent name in link reference, got:\n%s", out)
	}
}
