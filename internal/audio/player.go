package audio

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gen2brain/beeep"
)

// SystemPlayer implements Player using system audio capabilities
type SystemPlayer struct {
	config     *Config
	soundPaths map[SoundType]string
}

// newSystemPlayer creates a new system audio player
func newSystemPlayer(config *Config) (*SystemPlayer, error) {
	player := &SystemPlayer{
		config:     config,
		soundPaths: make(map[SoundType]string),
	}

	// Resolve sound file paths
	if err := player.resolveSoundPaths(); err != nil {
		return nil, fmt.Errorf("failed to resolve sound paths: %w", err)
	}

	return player, nil
}

// resolveSoundPaths finds the actual file paths for configured sounds
func (p *SystemPlayer) resolveSoundPaths() error {
	for soundTypeStr, filename := range p.config.Sounds {
		soundType := SoundType(soundTypeStr)

		// Try custom sounds directory first
		customPath := filepath.Join(p.config.CustomSoundsDir, filename)
		if _, err := os.Stat(customPath); err == nil {
			p.soundPaths[soundType] = customPath
			continue
		}

		// Try built-in sounds directory
		builtinPath := filepath.Join("internal", "audio", "sounds", filename)
		if _, err := os.Stat(builtinPath); err == nil {
			p.soundPaths[soundType] = builtinPath
			continue
		}

		// For now, we'll use system beep as fallback
		// In a full implementation, you'd embed default sound files
		p.soundPaths[soundType] = ""
	}

	return nil
}

// Play plays the specified sound type
func (p *SystemPlayer) Play(soundType SoundType) error {
	if !p.config.Enabled {
		return nil
	}

	soundPath, exists := p.soundPaths[soundType]
	if !exists {
		return fmt.Errorf("sound type %s not configured", soundType)
	}

	// If we have a sound file, try to play it
	if soundPath != "" {
		return p.playFile(soundPath)
	}

	// Fallback to system beep
	return p.playSystemBeep()
}

// playFile attempts to play an audio file
func (p *SystemPlayer) playFile(path string) error {
	// For now, we'll use a simple approach
	// In a full implementation, you'd use a proper audio library like:
	// - github.com/hajimehoshi/oto for cross-platform audio
	// - github.com/faiface/beep for audio file format support

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("sound file not found: %s", path)
	}

	// For macOS, we can use afplay
	// For Linux, we could use aplay, paplay, or similar
	// For Windows, we could use a system call

	// This is a simplified implementation - in production you'd want
	// proper cross-platform audio support
	return p.playSystemBeep()
}

// playSystemBeep plays a system beep sound
func (p *SystemPlayer) playSystemBeep() error {
	// Use beeep library's Beep function for cross-platform system sound
	return beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration)
}

// SetVolume sets the playback volume (0.0 to 1.0)
func (p *SystemPlayer) SetVolume(volume float64) error {
	if volume < 0.0 || volume > 1.0 {
		return errors.New("volume must be between 0.0 and 1.0")
	}

	p.config.Volume = volume
	return nil
}

// IsEnabled returns whether audio is enabled
func (p *SystemPlayer) IsEnabled() bool {
	return p.config.Enabled
}

// Close cleans up any resources
func (p *SystemPlayer) Close() error {
	// Clean up any audio resources if needed
	return nil
}
