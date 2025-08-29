package utils

import (
	"testing"
	"time"
)

func TestValidateDuration(t *testing.T) {
	cases := []struct {
		name    string
		d       time.Duration
		wantErr bool
	}{
		{"negative", -1 * time.Second, true},
		{"zero", 0, true},
		{"too_small", 500 * time.Millisecond, true},
		{"ok_one_second", 1 * time.Second, false},
		{"ok_minute", 1 * time.Minute, false},
		{"too_large", 25 * time.Hour, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateDuration(c.d)
			if c.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateDurationString(t *testing.T) {
	cases := []struct {
		name    string
		s       string
		wantErr bool
	}{
		{"empty", "", true},
		{"invalid_format", "twenty", true},
		{"too_small", "0.5s", true},
		{"valid_seconds", "5s", false},
		{"valid_combo", "1h30m", false},
		{"too_large", "25h", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidateDurationString(c.s)
			if c.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	cases := []struct {
		name     string
		desc     string
		required bool
		wantErr  bool
	}{
		{"required_empty", "   ", true, true},
		{"optional_empty", "", false, false},
		{"within_limit", "Focus on docs", true, false},
		{"too_long", string(make([]byte, 201)), true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// make too_long description explicitly
			if c.name == "too_long" {
				long := make([]byte, 201)
				for i := range long {
					long[i] = 'a'
				}
				c.desc = string(long)
			}
			err := ValidateDescription(c.desc, c.required)
			if c.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !c.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateTagsAndSanitize(t *testing.T) {
	// SanitizeTags removes empties, trims, lowercases and de-dupes
	in := []string{" Go ", "go", "Productivity", " ", "Focus"}
	got := SanitizeTags(in)
	want := []string{"go", "productivity", "focus"}
	if len(got) != len(want) {
		t.Fatalf("SanitizeTags length = %d; want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SanitizeTags[%d] = %q; want %q", i, got[i], want[i])
		}
	}

	// ValidateTags error cases
	if err := ValidateTags([]string{""}); err == nil {
		t.Fatalf("expected error for empty tag")
	}
	if err := ValidateTags([]string{"a,b"}); err == nil {
		t.Fatalf("expected error for comma in tag")
	}
	longTag := make([]byte, 51)
	for i := range longTag {
		longTag[i] = 'x'
	}
	if err := ValidateTags([]string{string(longTag)}); err == nil {
		t.Fatalf("expected error for tag too long")
	}

	// Max tags
	tooMany := make([]string, 11)
	for i := range tooMany {
		tooMany[i] = "t"
	}
	if err := ValidateTags(tooMany); err == nil {
		t.Fatalf("expected error for too many tags")
	}

	// Valid
	if err := ValidateTags([]string{"go", "cli"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateVolumeAndSanitizeDescription(t *testing.T) {
	// Volume
	if err := ValidateVolume(-0.1); err == nil {
		t.Fatalf("expected error for volume < 0")
	}
	if err := ValidateVolume(1.1); err == nil {
		t.Fatalf("expected error for volume > 1")
	}
	if err := ValidateVolume(0.5); err != nil {
		t.Fatalf("unexpected error for valid volume: %v", err)
	}

	// SanitizeDescription condenses spaces and trims
	if got := SanitizeDescription("  Hello   world  "); got != "Hello world" {
		t.Fatalf("SanitizeDescription = %q; want %q", got, "Hello world")
	}
}
