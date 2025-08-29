package utils

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		in   time.Duration
		want string
	}{
		{"zero", 0, "00:00"},
		{"seconds", 42 * time.Second, "00:42"},
		{"minutes", 3 * time.Minute, "03:00"},
		{"minutes+seconds", 2*time.Minute + 5*time.Second, "02:05"},
		{"over_hour", 1*time.Hour + 2*time.Minute + 3*time.Second, "62:03"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuration(tt.in); got != tt.want {
				t.Fatalf("FormatDuration(%v) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFormatDurationLong(t *testing.T) {
	tests := []struct {
		name string
		in   time.Duration
		want string
	}{
		{"seconds_only", 9 * time.Second, "9s"},
		{"minutes_seconds", 1*time.Minute + 2*time.Second, "1m 02s"},
		{"hours_minutes_seconds", 2*time.Hour + 3*time.Minute + 4*time.Second, "2h 03m 04s"},
		{"hours_zero_minutes", 1 * time.Hour, "1h 00m 00s"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDurationLong(tt.in); got != tt.want {
				t.Fatalf("FormatDurationLong(%v) = %q; want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseDurationWithDefaults(t *testing.T) {
	def := 25 * time.Minute
	tests := []struct {
		name string
		in   string
		want time.Duration
	}{
		{"empty_uses_default", "", def},
		{"invalid_uses_default", "abc", def},
		{"valid_overrides_default", "5m", 5 * time.Minute},
		{"valid_seconds", "90s", 90 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDurationWithDefaults(tt.in, def); got != tt.want {
				t.Fatalf("ParseDurationWithDefaults(%q) = %v; want %v", tt.in, got, tt.want)
			}
		})
	}
}
