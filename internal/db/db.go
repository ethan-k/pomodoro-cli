package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var _ DB = (*InternalDB)(nil)

type InternalDB struct {
	db *sql.DB
}

type DB interface {
	CreateSession(startTime, endTime time.Time, description string, durationSec int64, tagsCSV string, wasBreak bool) (int64, error)
	GetActiveSession() (*PomodoroSession, error)
	GetPausedSession() (*PomodoroSession, error)
	GetLastSession() (*PomodoroSession, error)
	UpdateSessionEndTime(id int64, endTime time.Time) error
	PauseSession(id int64, pausedAt time.Time) error
	ResumeSession(id int64, newEndTime time.Time) error
	GetSessionsByDateRange(startDate, endDate time.Time) ([]PomodoroSession, error)
	GetTodaySessions() ([]PomodoroSession, error)
	Close() error
}

type PomodoroSession struct {
	ID                  int64
	StartTime           time.Time
	EndTime             time.Time
	Description         string
	DurationSec         int64
	Tags                []string
	TagsCSV             string
	WasBreak            bool
	PausedAt            *time.Time
	TotalPausedDuration int64
	IsPaused            bool
}

func NewDB() (*InternalDB, error) {
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

	// Create base table
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
		return nil, fmt.Errorf("error creating base table: %v", err)
	}

	// Add new columns if they don't exist (for database migration)
	migrations := []string{
		`ALTER TABLE pomodoros ADD COLUMN paused_at TIMESTAMP;`,
		`ALTER TABLE pomodoros ADD COLUMN total_paused_duration INTEGER DEFAULT 0;`,
		`ALTER TABLE pomodoros ADD COLUMN is_paused BOOLEAN DEFAULT 0;`,
		`CREATE INDEX IF NOT EXISTS idx_pomodoros_active ON pomodoros(is_paused, end_time);`,
	}

	for _, migration := range migrations {
		// Ignore errors for columns that already exist
		db.Exec(migration)
	}

	return &InternalDB{db: db}, nil
}

func (d *InternalDB) Close() error {
	return d.db.Close()
}

func (d *InternalDB) CreateSession(startTime, endTime time.Time, description string, durationSec int64, tagsCSV string, wasBreak bool) (int64, error) {
	res, err := d.db.Exec(
		`INSERT INTO pomodoros(start_time, end_time, description, duration_secs, tags_csv, was_break) VALUES(?, ?, ?, ?, ?, ?)`,
		startTime, endTime, description, durationSec, tagsCSV, wasBreak,
	)
	if err != nil {
		return 0, fmt.Errorf("error inserting record: %v", err)
	}

	return res.LastInsertId()
}

func (d *InternalDB) GetActiveSession() (*PomodoroSession, error) {
	now := time.Now()

	var session PomodoroSession
	err := d.db.QueryRow(
		`SELECT id, start_time, end_time, description, duration_secs, tags_csv, was_break, 
		        paused_at, total_paused_duration, is_paused 
		FROM pomodoros 
		WHERE (end_time > ? AND is_paused = 0) OR is_paused = 1
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
		&session.PausedAt,
		&session.TotalPausedDuration,
		&session.IsPaused,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying active session: %v", err)
	}

	return &session, nil
}

func (d *InternalDB) GetPausedSession() (*PomodoroSession, error) {
	var session PomodoroSession
	err := d.db.QueryRow(
		`SELECT id, start_time, end_time, description, duration_secs, tags_csv, was_break, 
		        paused_at, total_paused_duration, is_paused 
		FROM pomodoros 
		WHERE is_paused = 1
		ORDER BY start_time DESC LIMIT 1`,
	).Scan(
		&session.ID,
		&session.StartTime,
		&session.EndTime,
		&session.Description,
		&session.DurationSec,
		&session.TagsCSV,
		&session.WasBreak,
		&session.PausedAt,
		&session.TotalPausedDuration,
		&session.IsPaused,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying paused session: %v", err)
	}

	return &session, nil
}

func (d *InternalDB) GetLastSession() (*PomodoroSession, error) {
	var session PomodoroSession
	err := d.db.QueryRow(
		`SELECT id, start_time, end_time, description, duration_secs, tags_csv, was_break,
		        paused_at, total_paused_duration, is_paused
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
		&session.PausedAt,
		&session.TotalPausedDuration,
		&session.IsPaused,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying last session: %v", err)
	}

	return &session, nil
}

func (d *InternalDB) UpdateSessionEndTime(id int64, endTime time.Time) error {
	_, err := d.db.Exec(
		`UPDATE pomodoros SET end_time = ? WHERE id = ?`,
		endTime, id,
	)
	return err
}

func (d *InternalDB) PauseSession(id int64, pausedAt time.Time) error {
	_, err := d.db.Exec(
		`UPDATE pomodoros SET paused_at = ?, is_paused = 1 WHERE id = ?`,
		pausedAt, id,
	)
	return err
}

func (d *InternalDB) ResumeSession(id int64, newEndTime time.Time) error {
	// First, get the current paused duration
	var currentPausedAt time.Time
	var totalPausedDuration int64

	err := d.db.QueryRow(
		`SELECT paused_at, total_paused_duration FROM pomodoros WHERE id = ?`,
		id,
	).Scan(&currentPausedAt, &totalPausedDuration)

	if err != nil {
		return fmt.Errorf("error getting paused session data: %v", err)
	}

	// Calculate additional paused time
	now := time.Now()
	additionalPausedTime := now.Sub(currentPausedAt)
	newTotalPausedDuration := totalPausedDuration + int64(additionalPausedTime.Seconds())

	// Update the session
	_, err = d.db.Exec(
		`UPDATE pomodoros SET 
			end_time = ?, 
			paused_at = NULL, 
			total_paused_duration = ?, 
			is_paused = 0 
		WHERE id = ?`,
		newEndTime, newTotalPausedDuration, id,
	)
	return err
}

func (d *InternalDB) GetSessionsByDateRange(startDate, endDate time.Time) ([]PomodoroSession, error) {
	rows, err := d.db.Query(
		`SELECT id, start_time, end_time, description, duration_secs, tags_csv, was_break,
		        paused_at, total_paused_duration, is_paused
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
			&session.PausedAt,
			&session.TotalPausedDuration,
			&session.IsPaused,
		); err != nil {
			return nil, fmt.Errorf("error scanning session: %v", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (d *InternalDB) GetTodaySessions() ([]PomodoroSession, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	return d.GetSessionsByDateRange(today, tomorrow)
}
