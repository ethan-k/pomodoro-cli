package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats a duration in MM:SS format
func FormatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// FormatDurationLong formats a duration in a more human-readable format
func FormatDurationLong(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %02dm %02ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %02ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// ParseDurationWithDefaults parses a duration string with sensible defaults
func ParseDurationWithDefaults(s string, defaultDuration time.Duration) time.Duration {
	if s == "" {
		return defaultDuration
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return defaultDuration
	}

	return duration
}
