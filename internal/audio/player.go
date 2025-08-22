package audio

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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

		// Try to find built-in sounds directory (relative to executable or source)
		possiblePaths := []string{
			filepath.Join("internal", "audio", "sounds", filename),
			filepath.Join("audio", "sounds", filename),
			filepath.Join("sounds", filename),
		}

		found := false
		for _, builtinPath := range possiblePaths {
			if _, err := os.Stat(builtinPath); err == nil {
				p.soundPaths[soundType] = builtinPath
				found = true
				break
			}
		}

		if !found {
			// Use system beep as fallback
			p.soundPaths[soundType] = ""
		}
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
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("sound file not found: %s", path)
	}

	// Try platform-specific audio players
	// For macOS, use afplay
	if err := p.tryMacOSPlayer(path); err == nil {
		return nil
	}

	// For Linux, try common audio players
	if err := p.tryLinuxPlayer(path); err == nil {
		return nil
	}

	// Fallback to system beep if no audio player works
	return p.playSystemBeep()
}

// tryMacOSPlayer attempts to play audio using macOS afplay
func (p *SystemPlayer) tryMacOSPlayer(path string) error {
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("not on macOS")
	}
	
	cmd := exec.Command("afplay", path)
	return cmd.Run()
}

// tryLinuxPlayer attempts to play audio using common Linux audio players
func (p *SystemPlayer) tryLinuxPlayer(path string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("not on Linux")
	}
	
	// Try different Linux audio players in order of preference
	players := []string{"paplay", "aplay", "play"}
	
	for _, player := range players {
		// Check if player exists
		if _, err := exec.LookPath(player); err != nil {
			continue
		}
		
		cmd := exec.Command(player, path) // #nosec G204 - player is validated with exec.LookPath, path is embedded resource
		if err := cmd.Run(); err == nil {
			return nil
		}
	}
	
	return fmt.Errorf("no suitable audio player found")
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
