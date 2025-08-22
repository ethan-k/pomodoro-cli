// Package model contains the TUI models for the Pomodoro timer progress display
package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethan-k/pomodoro-cli/internal/utils"
)

const (
	padding  = 2
	maxWidth = 80
)

// TickMsg is sent when the timer ticks
type TickMsg time.Time

// PomodoroModel represents a Pomodoro timer model for bubbletea
type PomodoroModel struct {
	ID          int64
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	IsBreak     bool
	progress    progress.Model
	quitting    bool
}

// NewPomodoroModel creates a new Pomodoro timer model
func NewPomodoroModel(id int64, description string, startTime time.Time, duration time.Duration, isBreak bool) PomodoroModel {
	var p progress.Model

	if isBreak {
		// Green colors for break
		p = progress.New(
			progress.WithGradient("#5A8A20", "#98D44A"),
			progress.WithWidth(40),
			progress.WithoutPercentage(),
		)
	} else {
		// Default gradient for pomodoro (usually pinkish)
		p = progress.New(
			progress.WithDefaultGradient(),
			progress.WithWidth(40),
			progress.WithoutPercentage(),
		)
	}

	return PomodoroModel{
		ID:          id,
		Description: description,
		StartTime:   startTime,
		EndTime:     startTime.Add(duration),
		Duration:    duration,
		IsBreak:     isBreak,
		progress:    p,
	}
}

// Init initializes the model
func (m PomodoroModel) Init() tea.Cmd {
	return tea.Batch(
		tickEvery(time.Second),
	)
}

// Update handles messages and updates the model
func (m PomodoroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			m.quitting = true
			return m, tea.Quit
		}
	case TickMsg:
		if time.Now().After(m.EndTime) {
			m.quitting = true
			return m, tea.Quit
		}
		return m, tickEvery(time.Second)
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 20
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
	case progress.FrameMsg:
		// Handle animation frames
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	// Update progress percentage based on elapsed time
	cmd := m.updateProgress()
	return m, cmd
}

func (m *PomodoroModel) updateProgress() tea.Cmd {
	now := time.Now()
	elapsed := now.Sub(m.StartTime)

	// Ensure progress doesn't exceed 1.0
	var percent float64
	if elapsed >= m.Duration {
		percent = 1.0
	} else {
		percent = float64(elapsed) / float64(m.Duration)
	}

	// Set the progress percentage (this will animate smoothly)
	return m.progress.SetPercent(percent)
}

// View renders the model
func (m PomodoroModel) View() string {
	now := time.Now()

	if m.quitting || now.After(m.EndTime) {
		return "Completed!\n"
	}

	remaining := m.EndTime.Sub(now).Round(time.Second)
	remainingStr := utils.FormatDuration(remaining)

	emoji := "🍅"
	if m.IsBreak {
		emoji = "☕"
	}

	pad := strings.Repeat(" ", padding)
	progressBar := m.progress.View()

	return fmt.Sprintf("\n%s%s  %s %s  %s\n",
		pad,
		progressBar,
		remainingStr,
		emoji,
		m.Description)
}

// tickEvery returns a command that ticks at the specified interval
func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
