package glossymetric

import (
	"fmt"
	"io"
	"strings"

	"go.opentelemetry.io/otel/attribute"
)

func printAttrs(attrs attribute.Set, w io.Writer) {
	kvs := attrs.ToSlice()
	if len(kvs) == 0 {
		return
	}

	parts := make([]string, len(kvs))
	for i, a := range kvs {
		key := renderStyled(styleAttrKey, string(a.Key))
		val := renderStyled(styleAttrVal, a.Value.Emit())
		parts[i] = fmt.Sprintf("%s=%s", key, val)
	}

	_, _ = fmt.Fprintf(w, "  %s %s\n", renderStyled(styleLabel, "attrs:"), strings.Join(parts, ", "))
}
