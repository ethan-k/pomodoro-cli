package notify

import (
	"fmt"

	"github.com/gen2brain/beeep"
)

// NotifyComplete sends a notification when a Pomodoro or break is complete
func NotifyComplete(title, message string) error {
	return beeep.Notify(title, message, "")
}

// NotifyPomodoroComplete sends a notification when a Pomodoro is complete
func NotifyPomodoroComplete(description string) error {
	title := "Pomodoro Complete"
	message := fmt.Sprintf("Task completed: %s", description)
	return NotifyComplete(title, message)
}

// NotifyBreakComplete sends a notification when a break is complete
func NotifyBreakComplete() error {
	title := "Break Complete"
	message := "Break time is over. Resume work."
	return NotifyComplete(title, message)
}
