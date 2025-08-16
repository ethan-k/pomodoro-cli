package notify

import (
	"fmt"

	"github.com/ethan-k/pomodoro-cli/internal/audio"
	"github.com/ethan-k/pomodoro-cli/internal/config"
	"github.com/gen2brain/beeep"
)

// NotifyComplete sends a notification when a Pomodoro or break is complete
func NotifyComplete(title, message string) error {
	return beeep.Notify(title, message, "")
}

// NotifyWithAudio sends both visual and audio notifications
func NotifyWithAudio(title, message string, soundType audio.SoundType, silentMode bool) error {
	// Send visual notification
	if err := NotifyComplete(title, message); err != nil {
		return err
	}

	// Send audio notification if not in silent mode
	if !silentMode {
		cfg, err := config.LoadConfig()
		if err == nil && cfg.Audio != nil {
			player, err := audio.NewPlayer(cfg.Audio)
			if err == nil {
				audio.PlayAsync(player, soundType)
			}
		}
	}

	return nil
}

// NotifyPomodoroComplete sends a notification when a Pomodoro is complete
func NotifyPomodoroComplete(description string) error {
	title := "Pomodoro Complete"
	message := fmt.Sprintf("Task completed: %s", description)
	return NotifyWithAudio(title, message, audio.PomodoroComplete, false)
}

// NotifyPomodoroCompleteWithOptions sends a notification with audio options
func NotifyPomodoroCompleteWithOptions(description string, silentMode bool) error {
	title := "Pomodoro Complete"
	message := fmt.Sprintf("Task completed: %s", description)
	return NotifyWithAudio(title, message, audio.PomodoroComplete, silentMode)
}

// NotifyBreakComplete sends a notification when a break is complete
func NotifyBreakComplete() error {
	title := "Break Complete"
	message := "Break time is over. Resume work."
	return NotifyWithAudio(title, message, audio.BreakComplete, false)
}

// NotifyBreakCompleteWithOptions sends a notification with audio options
func NotifyBreakCompleteWithOptions(silentMode bool) error {
	title := "Break Complete"
	message := "Break time is over. Resume work."
	return NotifyWithAudio(title, message, audio.BreakComplete, silentMode)
}
