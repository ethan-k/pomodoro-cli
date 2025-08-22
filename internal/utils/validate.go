//nolint:revive // utils is a common and acceptable package name for this purpose
package utils

import (
	"errors"
	"strings"
	"time"
)

// ValidateDuration validates a duration value
func ValidateDuration(d time.Duration) error {
	if d <= 0 {
		return errors.New("duration must be positive")
	}
	if d > 24*time.Hour {
		return errors.New("duration cannot exceed 24 hours")
	}
	if d < time.Second {
		return errors.New("duration must be at least 1 second")
	}
	return nil
}

// ValidateDescription validates a session description
func ValidateDescription(desc string, required bool) error {
	trimmed := strings.TrimSpace(desc)
	if required && trimmed == "" {
		return errors.New("description cannot be empty")
	}
	if len(trimmed) > 200 {
		return errors.New("description cannot exceed 200 characters")
	}
	return nil
}

// ValidateTags validates session tags
func ValidateTags(tags []string) error {
	if len(tags) > 10 {
		return errors.New("cannot have more than 10 tags")
	}

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return errors.New("tags cannot be empty")
		}
		if len(tag) > 50 {
			return errors.New("individual tags cannot exceed 50 characters")
		}
		if strings.Contains(tag, ",") {
			return errors.New("tags cannot contain commas")
		}
	}

	return nil
}

// ValidateVolume validates audio volume level
func ValidateVolume(volume float64) error {
	if volume < 0.0 || volume > 1.0 {
		return errors.New("volume must be between 0.0 and 1.0")
	}
	return nil
}

// SanitizeDescription cleans up a description string
func SanitizeDescription(desc string) string {
	// Trim whitespace
	desc = strings.TrimSpace(desc)

	// Replace multiple spaces with single space
	for strings.Contains(desc, "  ") {
		desc = strings.ReplaceAll(desc, "  ", " ")
	}

	return desc
}

// SanitizeTags cleans up tag strings
func SanitizeTags(tags []string) []string {
	var cleaned []string
	seen := make(map[string]bool)

	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		tag = strings.ToLower(tag) // Normalize to lowercase

		// Skip empty tags and duplicates
		if tag != "" && !seen[tag] {
			cleaned = append(cleaned, tag)
			seen[tag] = true
		}
	}

	return cleaned
}
