package glossytrace

import (
	"fmt"
	"time"
)

func formatDuration(d time.Duration) string {
	switch {
	case d < time.Microsecond:
		return fmt.Sprintf("%d ns", d.Nanoseconds())
	case d < time.Millisecond:
		return fmt.Sprintf("%d µs", d.Microseconds())
	case d < time.Second:
		return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000)
	default:
		return fmt.Sprintf("%.2f s", d.Seconds())
	}
}
