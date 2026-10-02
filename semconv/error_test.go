package semconv_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/foomo/opentelemetry-go/semconv"
)

func TestErrorType(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		err  error
		want string
	}{
		"canceled":         {err: context.Canceled, want: "context.Canceled"},
		"wrapped canceled": {err: fmt.Errorf("op: %w", context.Canceled), want: "context.Canceled"},
		"deadline":         {err: context.DeadlineExceeded, want: "context.DeadlineExceeded"},
		"other":            {err: errors.New("boom"), want: "*errors.errorString"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if got := semconv.ErrorType(tt.err).Value.AsString(); got != tt.want {
				t.Errorf("ErrorType() = %q, want %q", got, tt.want)
			}
		})
	}
}
