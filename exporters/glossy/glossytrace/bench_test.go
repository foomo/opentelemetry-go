package glossytrace_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/foomo/opentelemetry-go/exporters/glossy/glossytrace"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func makeSpans(n int) []sdktrace.ReadOnlySpan {
	spans := make([]sdktrace.ReadOnlySpan, n)
	for i := range n {
		parentID := trace.SpanID{}
		if i > 0 {
			parentID = rootSpan
		}

		spans[i] = newStubSpan(
			"span",
			trace.SpanID{byte(i + 1)},
			parentID,
			now,
			now.Add(time.Duration(i+1)*time.Millisecond),
			[]attribute.KeyValue{attribute.String("key", "value")},
		)
	}

	return spans
}

func BenchmarkExportSpans(b *testing.B) {
	b.Setenv("NO_COLOR", "1")

	for _, n := range []int{1, 10, 100} {
		spans := makeSpans(n)

		b.Run(
			func() string {
				switch n {
				case 1:
					return "1_span"
				case 10:
					return "10_spans"
				default:
					return "100_spans"
				}
			}(),
			func(b *testing.B) {
				exp, err := glossytrace.New(
					glossytrace.WithWriter(io.Discard),
					glossytrace.WithoutTimestamps(),
				)
				if err != nil {
					b.Fatal(err)
				}

				ctx := context.Background()

				b.ResetTimer()
				b.ReportAllocs()

				for range b.N {
					if err := exp.ExportSpans(ctx, spans); err != nil {
						b.Fatal(err)
					}
				}
			},
		)
	}
}

func BenchmarkExportSpansFlamegraph(b *testing.B) {
	b.Setenv("NO_COLOR", "1")

	spans := makeSpans(10)

	exp, err := glossytrace.New(
		glossytrace.WithWriter(io.Discard),
		glossytrace.WithFlamegraph(),
		glossytrace.WithoutTimestamps(),
	)
	if err != nil {
		b.Fatal(err)
	}

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		if err := exp.ExportSpans(ctx, spans); err != nil {
			b.Fatal(err)
		}
	}
}
