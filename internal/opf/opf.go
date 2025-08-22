// Package opf provides Open Pomodoro Format (OPF) export functionality
package opf

import (
	"encoding/json"
	"time"

	"github.com/ethan-k/pomodoro-cli/internal/db"
)

// Pomodoro represents a single Pomodoro session in OPF format
type Pomodoro struct {
	ID          string   `json:"id"`
	StartedAt   string   `json:"started_at"`
	Duration    int      `json:"duration"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Type        string   `json:"type"` // "pomodoro" or "break"
}

// Export represents the root object for Open Pomodoro Format export
type Export struct {
	Pomodoros []Pomodoro `json:"pomodoros"`
}

// ConvertToOPF converts a PomodoroSession to OPF format
func ConvertToOPF(session *db.PomodoroSession) Pomodoro {
	pomType := "pomodoro"
	if session.WasBreak {
		pomType = "break"
	}

	// Convert tags CSV to slice
	var tags []string
	if session.TagsCSV != "" {
		tags = splitTags(session.TagsCSV)
	}

	return Pomodoro{
		ID:          formatID(session.ID),
		StartedAt:   formatTime(session.StartTime),
		Duration:    int(session.DurationSec / 60), // Convert to minutes
		Description: session.Description,
		Tags:        tags,
		Type:        pomType,
	}
}

// ConvertSessionsToOPF converts multiple PomodoroSessions to OPF format
func ConvertSessionsToOPF(sessions []db.PomodoroSession) Export {
	opfPomodoros := make([]Pomodoro, 0, len(sessions))

	for _, session := range sessions {
		opfPomodoros = append(opfPomodoros, ConvertToOPF(&session))
	}

	return Export{
		Pomodoros: opfPomodoros,
	}
}

// ExportToJSON exports sessions to OPF JSON format
func ExportToJSON(sessions []db.PomodoroSession) ([]byte, error) {
	opfExport := ConvertSessionsToOPF(sessions)
	return json.MarshalIndent(opfExport, "", "  ")
}

// Helper functions
func formatID(_ int64) string {
	return time.Now().Format("20060102") + "-" + time.Now().Format("150405") + "-" + time.Now().Format("000")
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func splitTags(tagsCSV string) []string {
	if tagsCSV == "" {
		return nil
	}

	// Use strings.Split to convert CSV to slice
	// This is a simple implementation; in a real app, you might want to handle
	// escaping commas in tag values, trimming whitespace, etc.
	tags := make([]string, 0)
	start := 0
	inQuote := false

	for i := 0; i < len(tagsCSV); i++ {
		if tagsCSV[i] == '"' {
			inQuote = !inQuote
		} else if tagsCSV[i] == ',' && !inQuote {
			tags = append(tags, tagsCSV[start:i])
			start = i + 1
		}
	}

	// Add the last tag
	if start < len(tagsCSV) {
		tags = append(tags, tagsCSV[start:])
	}

	return tags
}
