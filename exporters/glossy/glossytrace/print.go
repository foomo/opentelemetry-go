package glossytrace

import (
	"fmt"
	"sort"
	"strings"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func (e *Exporter) printTrace(traceID string, spans []sdktrace.ReadOnlySpan) {
	w := e.cfg.writer

	header := fmt.Sprintf("=== TRACE %s ===", traceID)
	_, _ = fmt.Fprintln(w, renderStyled(styleHeader, header))

	sort.Slice(spans, func(i, j int) bool {
		return spans[i].StartTime().Before(spans[j].StartTime())
	})

	children := map[string][]sdktrace.ReadOnlySpan{}

	var roots []sdktrace.ReadOnlySpan

	for _, s := range spans {
		pid := s.Parent().SpanID().String()
		if pid == "0000000000000000" || pid == "" {
			roots = append(roots, s)
		} else {
			children[pid] = append(children[pid], s)
		}
	}

	for _, root := range roots {
		e.printNode(root, children, 0)
	}

	if e.cfg.flamegraph {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, renderStyled(styleHeader, "Nested Flamegraph:"))

		for _, root := range roots {
			e.printFlameNode(root, children, 0, root.EndTime().Sub(root.StartTime()))
		}
	}

	_, _ = fmt.Fprintln(w)
}

func (e *Exporter) printNode(span sdktrace.ReadOnlySpan, children map[string][]sdktrace.ReadOnlySpan, depth int) {
	w := e.cfg.writer
	prefix := strings.Repeat("  ", depth)
	duration := span.EndTime().Sub(span.StartTime())
	durStr := formatDuration(duration)
	durStyled := renderStyled(durationStyle(duration, e.cfg.durationWarn, e.cfg.durationCritical), durStr)

	connector := renderStyled(styleConnector, "└─")
	name := renderStyled(styleSpanName, span.Name())

	if e.cfg.timestamps {
		ts := renderStyled(styleTimestamp, span.StartTime().Format("15:04:05.000"))
		_, _ = fmt.Fprintf(w, "%s%s %s (%s) [%s]\n", prefix, connector, name, durStyled, ts)
	} else {
		_, _ = fmt.Fprintf(w, "%s%s %s (%s)\n", prefix, connector, name, durStyled)
	}

	if e.cfg.showAttributes {
		e.printAttributes(span, prefix+"    ")
		e.printEvents(span, prefix+"    ")
	}

	for _, child := range children[span.SpanContext().SpanID().String()] {
		e.printNode(child, children, depth+1)
	}
}

func (e *Exporter) printFlameNode(span sdktrace.ReadOnlySpan, children map[string][]sdktrace.ReadOnlySpan, depth int, total time.Duration) {
	w := e.cfg.writer
	duration := span.EndTime().Sub(span.StartTime())

	barLen := max(int(float64(30)*duration.Seconds()/total.Seconds()), 1)

	prefix := strings.Repeat("  ", depth)
	bar := renderStyled(styleFlameBar, strings.Repeat("█", barLen))
	name := renderStyled(styleSpanName, fmt.Sprintf("%-15s", span.Name()))
	durStr := formatDuration(duration)
	durStyled := renderStyled(durationStyle(duration, e.cfg.durationWarn, e.cfg.durationCritical), durStr)

	_, _ = fmt.Fprintf(w, "%s%s %s %s\n", prefix, name, bar, durStyled)

	for _, child := range children[span.SpanContext().SpanID().String()] {
		e.printFlameNode(child, children, depth+1, total)
	}
}

func (e *Exporter) printAttributes(span sdktrace.ReadOnlySpan, indent string) {
	w := e.cfg.writer

	attrs := span.Attributes()
	if len(attrs) == 0 {
		return
	}

	for _, attr := range attrs {
		key := renderStyled(styleAttrKey, string(attr.Key))
		val := renderStyled(styleAttrVal, attr.Value.Emit())
		_, _ = fmt.Fprintf(w, "%s%s=%s\n", indent, key, val)
	}
}

func (e *Exporter) printEvents(span sdktrace.ReadOnlySpan, indent string) {
	w := e.cfg.writer

	events := span.Events()
	if len(events) == 0 {
		return
	}

	for _, event := range events {
		name := renderStyled(styleEventName, event.Name)
		ts := renderStyled(styleTimestamp, event.Time.Format("15:04:05.000"))
		_, _ = fmt.Fprintf(w, "%s@ %s [%s]\n", indent, name, ts)

		for _, attr := range event.Attributes {
			key := renderStyled(styleAttrKey, string(attr.Key))
			val := renderStyled(styleAttrVal, attr.Value.Emit())
			_, _ = fmt.Fprintf(w, "%s  %s=%s\n", indent, key, val)
		}
	}
}
