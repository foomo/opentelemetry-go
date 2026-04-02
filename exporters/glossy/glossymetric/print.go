package glossymetric

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func (e *Exporter) prettyPrint(rm *metricdata.ResourceMetrics) error {
	w := e.cfg.writer

	if len(rm.ScopeMetrics) == 0 {
		return nil
	}

	_, _ = fmt.Fprintln(w, renderStyled(styleHeader, "METRICS"))

	sep := renderStyled(styleSeparator, strings.Repeat("─", 50))

	for _, sm := range rm.ScopeMetrics {
		_, _ = fmt.Fprintf(w, "\n%s %s\n", renderStyled(styleLabel, "Scope:"), renderStyled(styleMetricName, sm.Scope.Name))
		for _, m := range sm.Metrics {
			_, _ = fmt.Fprintln(w, sep)
			_, _ = fmt.Fprintf(w, "%s %s\n", renderStyled(styleLabel, "Metric:"), renderStyled(styleMetricName, m.Name))
			_, _ = fmt.Fprintf(w, "%s %s\n", renderStyled(styleLabel, "Description:"), m.Description)
			_, _ = fmt.Fprintf(w, "%s %s\n", renderStyled(styleLabel, "Unit:"), renderStyled(styleUnit, m.Unit))

			switch data := m.Data.(type) {
			case metricdata.Sum[int64]:
				for _, dp := range data.DataPoints {
					printAttrs(dp.Attributes, w)
					_, _ = fmt.Fprintf(w, "%s %s\n", renderStyled(styleLabel, "Value:"), renderStyled(styleValue, fmt.Sprintf("%d", dp.Value)))
				}
			case metricdata.Sum[float64]:
				for _, dp := range data.DataPoints {
					printAttrs(dp.Attributes, w)
					_, _ = fmt.Fprintf(w, "%s %s\n", renderStyled(styleLabel, "Value:"), renderStyled(styleValue, fmt.Sprintf("%.3f", dp.Value)))
				}
			case metricdata.Gauge[int64]:
				for _, dp := range data.DataPoints {
					printAttrs(dp.Attributes, w)
					_, _ = fmt.Fprintf(w, "%s %s\n", renderStyled(styleLabel, "Value:"), renderStyled(styleValue, fmt.Sprintf("%d", dp.Value)))
				}
			case metricdata.Gauge[float64]:
				for _, dp := range data.DataPoints {
					printAttrs(dp.Attributes, w)
					_, _ = fmt.Fprintf(w, "%s %s\n", renderStyled(styleLabel, "Value:"), renderStyled(styleValue, fmt.Sprintf("%.3f", dp.Value)))
				}
			case metricdata.Histogram[int64]:
				if e.cfg.histograms {
					for _, dp := range data.DataPoints {
						printAttrs(dp.Attributes, w)
						vMin, _ := dp.Min.Value()
						vMax, _ := dp.Max.Value()

						avg := 0.0
						if dp.Count > 0 {
							avg = float64(dp.Sum) / float64(dp.Count)
						}

						_, _ = fmt.Fprintf(w, "%s min=%s max=%s avg=%s count=%s\n",
							renderStyled(styleLabel, "Histogram:"),
							renderStyled(styleValue, fmt.Sprintf("%d", vMin)),
							renderStyled(styleValue, fmt.Sprintf("%d", vMax)),
							renderStyled(styleValue, fmt.Sprintf("%.3f", avg)),
							renderStyled(styleValue, fmt.Sprintf("%d", dp.Count)),
						)
					}
				}
			case metricdata.Histogram[float64]:
				if e.cfg.histograms {
					for _, dp := range data.DataPoints {
						printAttrs(dp.Attributes, w)
						vMin, _ := dp.Min.Value()
						vMax, _ := dp.Max.Value()

						avg := 0.0
						if dp.Count > 0 {
							avg = dp.Sum / float64(dp.Count)
						}

						_, _ = fmt.Fprintf(w, "%s min=%s max=%s avg=%s count=%s\n",
							renderStyled(styleLabel, "Histogram:"),
							renderStyled(styleValue, fmt.Sprintf("%.3f", vMin)),
							renderStyled(styleValue, fmt.Sprintf("%.3f", vMax)),
							renderStyled(styleValue, fmt.Sprintf("%.3f", avg)),
							renderStyled(styleValue, fmt.Sprintf("%d", dp.Count)),
						)
					}
				}
			}
		}

		_, _ = fmt.Fprintln(w, sep)
	}

	return nil
}
