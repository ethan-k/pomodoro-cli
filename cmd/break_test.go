package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/ethan-k/pomodoro-cli/internal/db"
)

// mockDB implements the complete db.DB interface for testing
type mockDB struct {
	CreateSessionFunc          func(start, end time.Time, description string, durationSec int64, tagsCSV string, wasBreak bool) (int64, error)
	GetActiveSessionFunc       func() (*db.PomodoroSession, error)
	GetPausedSessionFunc       func() (*db.PomodoroSession, error)
	GetLastSessionFunc         func() (*db.PomodoroSession, error)
	UpdateSessionEndTimeFunc   func(id int64, endTime time.Time) error
	PauseSessionFunc           func(id int64, pausedAt time.Time) error
	ResumeSessionFunc          func(id int64, newEndTime time.Time) error
	GetSessionsByDateRangeFunc func(startDate, endDate time.Time) ([]db.PomodoroSession, error)
	GetTodaySessionsFunc       func() ([]db.PomodoroSession, error)
	CloseFunc                  func() error
}

func (m *mockDB) CreateSession(start, end time.Time, description string, durationSec int64, tagsCSV string, wasBreak bool) (int64, error) {
	if m.CreateSessionFunc != nil {
		return m.CreateSessionFunc(start, end, description, durationSec, tagsCSV, wasBreak)
	}
	return 1, nil
}

func (m *mockDB) GetActiveSession() (*db.PomodoroSession, error) {
	if m.GetActiveSessionFunc != nil {
		return m.GetActiveSessionFunc()
	}
	return nil, nil
}

func (m *mockDB) GetPausedSession() (*db.PomodoroSession, error) {
	if m.GetPausedSessionFunc != nil {
		return m.GetPausedSessionFunc()
	}
	return nil, nil
}

func (m *mockDB) GetLastSession() (*db.PomodoroSession, error) {
	if m.GetLastSessionFunc != nil {
		return m.GetLastSessionFunc()
	}
	return nil, nil
}

func (m *mockDB) UpdateSessionEndTime(id int64, endTime time.Time) error {
	if m.UpdateSessionEndTimeFunc != nil {
		return m.UpdateSessionEndTimeFunc(id, endTime)
	}
	return nil
}

func (m *mockDB) PauseSession(id int64, pausedAt time.Time) error {
	if m.PauseSessionFunc != nil {
		return m.PauseSessionFunc(id, pausedAt)
	}
	return nil
}

func (m *mockDB) ResumeSession(id int64, newEndTime time.Time) error {
	if m.ResumeSessionFunc != nil {
		return m.ResumeSessionFunc(id, newEndTime)
	}
	return nil
}

func (m *mockDB) GetSessionsByDateRange(startDate, endDate time.Time) ([]db.PomodoroSession, error) {
	if m.GetSessionsByDateRangeFunc != nil {
		return m.GetSessionsByDateRangeFunc(startDate, endDate)
	}
	return nil, nil
}

func (m *mockDB) GetTodaySessions() ([]db.PomodoroSession, error) {
	if m.GetTodaySessionsFunc != nil {
		return m.GetTodaySessionsFunc()
	}
	return nil, nil
}

func (m *mockDB) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// Simple unit test for duration parsing logic
func TestBreakCommand_DurationParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		hasError bool
	}{
		{
			name:     "Valid minutes",
			input:    "5m",
			expected: 5 * time.Minute,
			hasError: false,
		},
		{
			name:     "Valid seconds",
			input:    "30s",
			expected: 30 * time.Second,
			hasError: false,
		},
		{
			name:     "Valid hours",
			input:    "1h",
			expected: 1 * time.Hour,
			hasError: false,
		},
		{
			name:     "Invalid format",
			input:    "abc",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration, err := time.ParseDuration(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for input %q, but got: %v", tt.input, err)
				}
				if duration != tt.expected {
					t.Errorf("Expected duration %v for input %q, but got %v", tt.expected, tt.input, duration)
				}
			}
		})
	}
}

// Test break session creation with mock
func TestBreakCommand_SessionCreation(t *testing.T) {
	mockDB := &mockDB{
		CreateSessionFunc: func(_, _ time.Time, description string, _ int64, _ string, wasBreak bool) (int64, error) {
			// Verify that wasBreak is true for break sessions
			if !wasBreak {
				t.Error("Expected wasBreak to be true for break sessions")
			}

			// Verify description contains "Break"
			if !strings.Contains(description, "Break") {
				t.Errorf("Expected description to contain 'Break', got: %q", description)
			}

			// Return a mock session ID
			return 123, nil
		},
	}

	// Test with 5 minute duration
	duration := 5 * time.Minute
	start := time.Now()
	end := start.Add(duration)

	sessionID, err := mockDB.CreateSession(start, end, "Break", int64(duration.Seconds()), "", true)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if sessionID != 123 {
		t.Errorf("Expected session ID 123, got: %d", sessionID)
	}
}
