package tfprofile

import (
	"fmt"
	"time"
)

// Format a duration in seconds into "30s" or "2m30s"
func FormatDuration(seconds int) string {
	duration := time.Duration(seconds) * time.Second
	minutes := int(duration.Minutes())
	seconds = seconds - (minutes * 60)
	if minutes == 0 {
		return fmt.Sprintf("%ds", seconds)
	}
	return fmt.Sprintf("%dm%ds", minutes, seconds)
}
