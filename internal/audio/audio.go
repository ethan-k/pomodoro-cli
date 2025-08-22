// Package audio provides sound notification functionality for the Pomodoro timer
package audio

import (
	"fmt"
	"os"
	"path/filepath"
)

// SoundType represents different types of audio notifications
type SoundType string

const (
	// PomodoroComplete represents the sound played when a Pomodoro session completes
	PomodoroComplete SoundType = "pomodoro_complete"
	// BreakComplete represents the sound played when a break session completes
	BreakComplete    SoundType = "break_complete"
	// SessionStart represents the sound played when starting a session
	SessionStart     SoundType = "session_start"
)

// Player interface for audio playback
type Player interface {
	Play(soundType SoundType) error
	SetVolume(volume float64) error
	IsEnabled() bool
	Close() error
}

// Config represents audio configuration
type Config struct {
	Enabled         bool              `yaml:"enabled"`
	Volume          float64           `yaml:"volume"`
	Sounds          map[string]string `yaml:"sounds"`
	CustomSoundsDir string            `yaml:"custom_sounds_dir"`
}

// DefaultConfig returns default audio configuration
func DefaultConfig() *Config {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	return &Config{
		Enabled: true,
		Volume:  0.5,  // Reduced default volume
		Sounds: map[string]string{
			string(PomodoroComplete): "pomodoro_complete.wav",
			string(BreakComplete):    "break_complete.wav",
			string(SessionStart):     "session_start.wav",
		},
		CustomSoundsDir: filepath.Join(home, ".config", "pomodoro", "sounds"),
	}
}

// NewPlayer creates a new audio player
func NewPlayer(config *Config) (Player, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if !config.Enabled {
		return &NoOpPlayer{}, nil
	}

	// Try to create platform-specific player
	player, err := newSystemPlayer(config)
	if err != nil {
		// Fallback to no-op player if audio system unavailable
		fmt.Printf("Warning: Audio unavailable, continuing without sound: %v\n", err)
		return &NoOpPlayer{}, nil
	}

	return player, nil
}

// PlayAsync plays sound without blocking
func PlayAsync(player Player, soundType SoundType) {
	if player == nil || !player.IsEnabled() {
		return
	}

	go func() {
		_ = player.Play(soundType) // Ignore audio playback errors in production
	}()
}

// NoOpPlayer is a no-operation player for when audio is disabled
type NoOpPlayer struct{}

// Play does nothing and returns no error for the no-op player
func (p *NoOpPlayer) Play(_ SoundType) error { return nil }

// SetVolume does nothing and returns no error for the no-op player  
func (p *NoOpPlayer) SetVolume(_ float64) error { return nil }

// IsEnabled always returns false for the no-op player
func (p *NoOpPlayer) IsEnabled() bool { return false }

// Close does nothing and returns no error for the no-op player
func (p *NoOpPlayer) Close() error { return nil }
