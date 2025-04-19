package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ethan-k/pomodoro-cli/internal/db"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Goals     GoalConfig     `yaml:"goals"`
	Hooks     HooksConfig    `yaml:"hooks"`
	Defaults  DefaultsConfig `yaml:"defaults"`
	DataPaths DataPaths      `yaml:"paths"`
}

// GoalConfig represents the goals configuration
type GoalConfig struct {
	DailyCount  int `yaml:"daily_count"`  // Target number of Pomodoros per day
	WeeklyCount int `yaml:"weekly_count"` // Target number of Pomodoros per week
}

// HooksConfig represents the hooks configuration
type HooksConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"` // Path to hooks directory
}

// DefaultsConfig represents default values
type DefaultsConfig struct {
	PomodoroDuration  string `yaml:"pomodoro_duration"`
	BreakDuration     string `yaml:"break_duration"`
	LongBreakDuration string `yaml:"long_break_duration"`
}

// DataPaths represents paths for data storage
type DataPaths struct {
	Database  string `yaml:"database"`
	OPFExport string `yaml:"opf_export"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return &Config{
		Goals: GoalConfig{
			DailyCount:  8,
			WeeklyCount: 40,
		},
		Hooks: HooksConfig{
			Enabled: false,
			Path:    filepath.Join(home, ".config", "pomodoro", "hooks"),
		},
		Defaults: DefaultsConfig{
			PomodoroDuration:  "25m",
			BreakDuration:     "5m",
			LongBreakDuration: "15m",
		},
		DataPaths: DataPaths{
			Database:  filepath.Join(home, ".local", "share", "pomodoro", "history.db"),
			OPFExport: filepath.Join(home, ".local", "share", "pomodoro", "exports"),
		},
	}
}

// LoadConfig loads the configuration from the default path
func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home dir: %v", err)
	}

	configPath := filepath.Join(home, ".config", "pomodoro", "config.yml")

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	// Parse config
	config := DefaultConfig() // Start with defaults
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return config, nil
}

// SaveConfig saves the configuration to the default path
func SaveConfig(config *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting home dir: %v", err)
	}

	configDir := filepath.Join(home, ".config", "pomodoro")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yml")

	// Marshal config to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	// Write config file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// GetCurrentGoalStatus returns the current goal status
func GetCurrentGoalStatus() (*GoalStatus, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	database, err := db.NewDB()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	// Get today's sessions
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	todaySessions, err := database.GetSessionsByDateRange(today, tomorrow)
	if err != nil {
		return nil, err
	}

	// Get this week's sessions
	now := time.Now()
	// Start from the beginning of the week (Monday)
	daysToMonday := int(now.Weekday())
	if daysToMonday == 0 { // Sunday
		daysToMonday = 6
	} else {
		daysToMonday = daysToMonday - 1
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-daysToMonday, 0, 0, 0, 0, now.Location())
	weekSessions, err := database.GetSessionsByDateRange(weekStart, now)
	if err != nil {
		return nil, err
	}

	// Count non-break sessions
	dailyCount := 0
	weeklyCount := 0
	for _, session := range todaySessions {
		if !session.WasBreak {
			dailyCount++
		}
	}
	for _, session := range weekSessions {
		if !session.WasBreak {
			weeklyCount++
		}
	}

	return &GoalStatus{
		DailyGoal:       config.Goals.DailyCount,
		DailyCompleted:  dailyCount,
		WeeklyGoal:      config.Goals.WeeklyCount,
		WeeklyCompleted: weeklyCount,
	}, nil
}

// GoalStatus represents the current goal status
type GoalStatus struct {
	DailyGoal       int
	DailyCompleted  int
	WeeklyGoal      int
	WeeklyCompleted int
}
