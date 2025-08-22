package audio

import (
	"fmt"
	"os"
	"path/filepath"
)

// SoundType represents different types of audio notifications
type SoundType string

const (
	PomodoroComplete SoundType = "pomodoro_complete"
	BreakComplete    SoundType = "break_complete"
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
		if err := player.Play(soundType); err != nil {
			// Intentionally ignoring audio playback errors to prevent disrupting user workflow
			// Audio is a nice-to-have feature, not critical functionality
		}
	}()
}

// NoOpPlayer is a no-operation player for when audio is disabled
type NoOpPlayer struct{}

func (p *NoOpPlayer) Play(soundType SoundType) error { return nil }
func (p *NoOpPlayer) SetVolume(volume float64) error { return nil }
func (p *NoOpPlayer) IsEnabled() bool                { return false }
func (p *NoOpPlayer) Close() error                   { return nil }
