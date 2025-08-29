package db

import (
    "os"
    "testing"
    "time"
)

// setTempHome sets HOME to a temporary directory for DB tests to avoid touching real user data.
func setTempHome(t *testing.T) string {
    t.Helper()
    dir := t.TempDir()
    t.Setenv("HOME", dir)
    return dir
}

func TestDB_CreateAndQuerySessions(t *testing.T) {
    _ = setTempHome(t)
    d, err := NewDB()
    if err != nil {
        t.Fatalf("NewDB error: %v", err)
    }
    t.Cleanup(func() { _ = d.Close() })

    // Create two sessions: one active pomodoro, one break in the past
    now := time.Now()
    // Active session ends in future
    activeEnd := now.Add(10 * time.Minute)
    id1, err := d.CreateSession(now.Add(-5*time.Minute), activeEnd, "Work", int64((activeEnd.Sub(now.Add(-5*time.Minute))).Seconds()), "go,cli", false)
    if err != nil || id1 == 0 {
        t.Fatalf("CreateSession active error: %v id=%d", err, id1)
    }

    // Past break
    pastStart := now.Add(-2 * time.Hour)
    pastEnd := pastStart.Add(5 * time.Minute)
    id2, err := d.CreateSession(pastStart, pastEnd, "Break", int64(5*60), "", true)
    if err != nil || id2 == 0 {
        t.Fatalf("CreateSession break error: %v id=%d", err, id2)
    }

    // GetActiveSession should find the first session
    s, err := d.GetActiveSession()
    if err != nil {
        t.Fatalf("GetActiveSession error: %v", err)
    }
    if s == nil || s.ID != id1 || s.WasBreak {
        t.Fatalf("unexpected active session: %+v", s)
    }

    // Update end time
    newEnd := activeEnd.Add(5 * time.Minute)
    if err := d.UpdateSessionEndTime(id1, newEnd); err != nil {
        t.Fatalf("UpdateSessionEndTime error: %v", err)
    }

    // Pause and Resume flow increases total_paused_duration and clears paused_at/is_paused
    pausedAt := time.Now().Add(-2 * time.Minute)
    if err := d.PauseSession(id1, pausedAt); err != nil {
        t.Fatalf("PauseSession error: %v", err)
    }
    // Ensure paused session is now considered active (is_paused = 1)
    ps, err := d.GetPausedSession()
    if err != nil || ps == nil || ps.ID != id1 || !ps.IsPaused {
        t.Fatalf("GetPausedSession unexpected: s=%+v err=%v", ps, err)
    }

    // Resume should accumulate paused duration
    if err := d.ResumeSession(id1, newEnd.Add(10*time.Minute)); err != nil {
        t.Fatalf("ResumeSession error: %v", err)
    }

    // Verify no paused session remains
    ps2, err := d.GetPausedSession()
    if err != nil {
        t.Fatalf("GetPausedSession after resume error: %v", err)
    }
    if ps2 != nil {
        t.Fatalf("expected no paused session after resume, got %+v", ps2)
    }

    // Date range query should include both sessions for today and earlier range
    today := time.Now().Truncate(24 * time.Hour)
    tomorrow := today.Add(24 * time.Hour)
    todaySessions, err := d.GetSessionsByDateRange(today, tomorrow)
    if err != nil {
        t.Fatalf("GetSessionsByDateRange today error: %v", err)
    }
    if len(todaySessions) == 0 {
        t.Fatalf("expected at least one session today")
    }

    // GetTodaySessions should be equivalent
    ts, err := d.GetTodaySessions()
    if err != nil {
        t.Fatalf("GetTodaySessions error: %v", err)
    }
    if len(ts) != len(todaySessions) {
        t.Fatalf("GetTodaySessions count mismatch: %d vs %d", len(ts), len(todaySessions))
    }

    // GetLastSession should return the most recent by start_time (id1)
    last, err := d.GetLastSession()
    if err != nil {
        t.Fatalf("GetLastSession error: %v", err)
    }
    if last == nil || last.ID != id1 {
        t.Fatalf("unexpected last session: %+v", last)
    }

    // Ensure DB file created in temp HOME and not elsewhere
    home, _ := os.UserHomeDir()
    if _, err := os.Stat(home + "/.local/share/pomodoro/history.db"); err != nil {
        t.Fatalf("db file not found in temp home: %v", err)
    }
}

