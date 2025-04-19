package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// PomodoroSession represents a Pomodoro session in the database
type PomodoroSession struct {
	ID          int64
	StartTime   time.Time
	EndTime     time.Time
	Description string
	DurationSec int64
	Tags        []string
	TagsCSV     string
	WasBreak    bool
}

// DB handles database operations
type DB struct {
	db *sql.DB
}

// NewDB creates a new database connection
func NewDB() (*DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home dir: %v", err)
	}

	dbPath := filepath.Join(home, ".local", "share", "pomodoro", "history.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("error creating DB dir: %v", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("error opening DB: %v", err)
	}

	// Create tables if they don't exist
	ddl := `CREATE TABLE IF NOT EXISTS pomodoros (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_time TIMESTAMP NOT NULL,
		end_time TIMESTAMP NOT NULL,
		description TEXT,
		duration_secs INTEGER NOT NULL,
		tags_csv TEXT,
		was_break BOOLEAN NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_pomodoros_day ON pomodoros(date(start_time));`

	if _, err := db.Exec(ddl); err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &DB{db: db}, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// CreateSession creates a new Pomodoro session
func (d *DB) CreateSession(startTime, endTime time.Time, description string, durationSec int64, tagsCSV string, wasBreak bool) (int64, error) {
	res, err := d.db.Exec(
		`INSERT INTO pomodoros(start_time, end_time, description, duration_secs, tags_csv, was_break) VALUES(?, ?, ?, ?, ?, ?)`,
		startTime, endTime, description, durationSec, tagsCSV, wasBreak,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting record: %v", err)
	}

	return res.LastInsertId()
}

// GetActiveSession returns the most recent active session if any
func (d *DB) GetActiveSession() (*PomodoroSession, error) {
	now := time.Now()

	var session PomodoroSession
	err := d.db.QueryRow(
		`SELECT id, start_time, end_time, description, duration_secs, tags_csv, was_break 
		FROM pomodoros 
		WHERE end_time > ? 
		ORDER BY start_time DESC LIMIT 1`,
		now,
	).Scan(
		&session.ID,
		&session.StartTime,
		&session.EndTime,
		&session.Description,
		&session.DurationSec,
		&session.TagsCSV,
		&session.WasBreak,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying active session: %v", err)
	}

	return &session, nil
}

// GetLastSession returns the most recent completed session
func (d *DB) GetLastSession() (*PomodoroSession, error) {
	var session PomodoroSession
	err := d.db.QueryRow(
		`SELECT id, start_time, end_time, description, duration_secs, tags_csv, was_break 
		FROM pomodoros 
		ORDER BY start_time DESC LIMIT 1`,
	).Scan(
		&session.ID,
		&session.StartTime,
		&session.EndTime,
		&session.Description,
		&session.DurationSec,
		&session.TagsCSV,
		&session.WasBreak,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying last session: %v", err)
	}

	return &session, nil
}

// UpdateSessionEndTime updates the end time of a session
func (d *DB) UpdateSessionEndTime(id int64, endTime time.Time) error {
	_, err := d.db.Exec(
		`UPDATE pomodoros SET end_time = ? WHERE id = ?`,
		endTime, id,
	)
	return err
}

// GetSessionsByDateRange returns sessions within a date range
func (d *DB) GetSessionsByDateRange(startDate, endDate time.Time) ([]PomodoroSession, error) {
	rows, err := d.db.Query(
		`SELECT id, start_time, end_time, description, duration_secs, tags_csv, was_break 
		FROM pomodoros 
		WHERE date(start_time) >= date(?) AND date(start_time) <= date(?)
		ORDER BY start_time DESC`,
		startDate, endDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying sessions: %v", err)
	}
	defer rows.Close()

	var sessions []PomodoroSession
	for rows.Next() {
		var session PomodoroSession
		if err := rows.Scan(
			&session.ID,
			&session.StartTime,
			&session.EndTime,
			&session.Description,
			&session.DurationSec,
			&session.TagsCSV,
			&session.WasBreak,
		); err != nil {
			return nil, fmt.Errorf("error scanning session: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// GetTodaySessions returns all sessions for today
func (d *DB) GetTodaySessions() ([]PomodoroSession, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	return d.GetSessionsByDateRange(today, tomorrow)
}
