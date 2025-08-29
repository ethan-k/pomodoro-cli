package model

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethan-k/pomodoro-cli/internal/goals"
)

// GoalDashboardModel represents the goal tracking dashboard
type GoalDashboardModel struct {
	width            int
	height           int
	goalManager      *goals.GoalManager
	dailyProgress    *goals.GoalProgress
	weeklyProgress   *goals.GoalProgress
	monthlyProgress  *goals.GoalProgress
	streak           *goals.StreakInfo
	history          []goals.DailyGoalResult
	dailyBar         progress.Model
	weeklyBar        progress.Model
	monthlyBar       progress.Model
	keys             keyMap
	help             help.Model
	showHistory      bool
	showAdjustment   bool
	newDailyTarget   string
	newWeeklyTarget  string
	loading          bool
	error            error
}

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	History  key.Binding
	Adjust   key.Binding
	Save     key.Binding
	Cancel   key.Binding
	Help     key.Binding
	Quit     key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.History, k.Adjust, k.Help, k.Quit}
}

// FullHelp returns keybindings to be shown in the full help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.History, k.Adjust, k.Save, k.Cancel},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	History: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "show history"),
	),
	Adjust: key.NewBinding(
		key.WithKeys("A"),
		key.WithHelp("A", "adjust goals"),
	),
	Save: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save changes"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// NewGoalDashboardModel creates a new goal dashboard model
func NewGoalDashboardModel(goalManager *goals.GoalManager) GoalDashboardModel {
	dailyBar := progress.New(
		progress.WithGradient("#FF6B6B", "#4ECDC4"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	
	weeklyBar := progress.New(
		progress.WithGradient("#4ECDC4", "#45B7D1"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	monthlyBar := progress.New(
		progress.WithGradient("#45B7D1", "#96CEB4"),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	return GoalDashboardModel{
		goalManager:     goalManager,
		dailyBar:        dailyBar,
		weeklyBar:       weeklyBar,
		monthlyBar:      monthlyBar,
		keys:            keys,
		help:            help.New(),
		loading:         true,
		newDailyTarget:  "",
		newWeeklyTarget: "",
	}
}

// LoadDataMsg represents a message to load goal data
type LoadDataMsg struct{}

// DataLoadedMsg represents loaded goal data
type DataLoadedMsg struct {
	Daily    *goals.GoalProgress
	Weekly   *goals.GoalProgress
	Monthly  *goals.GoalProgress
	Streak   *goals.StreakInfo
	History  []goals.DailyGoalResult
	Error    error
}

// Init initializes the goal dashboard model
func (m GoalDashboardModel) Init() tea.Cmd {
	return func() tea.Msg {
		return LoadDataMsg{}
	}
}

// Update handles messages and updates the goal dashboard model
func (m GoalDashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showAdjustment {
			switch msg.Type {
			case tea.KeyEnter:
				// Save goal adjustments
				m.showAdjustment = false
				return m, func() tea.Msg { return LoadDataMsg{} }
			case tea.KeyEsc:
				m.showAdjustment = false
				m.newDailyTarget = ""
				m.newWeeklyTarget = ""
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.History):
			m.showHistory = !m.showHistory
		case key.Matches(msg, m.keys.Adjust):
			m.showAdjustment = !m.showAdjustment
			if m.dailyProgress != nil {
				m.newDailyTarget = fmt.Sprintf("%d", m.dailyProgress.Target)
			}
			if m.weeklyProgress != nil {
				m.newWeeklyTarget = fmt.Sprintf("%d", m.weeklyProgress.Target)
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update progress bar widths
		barWidth := min(40, msg.Width-20)
		m.dailyBar.Width = barWidth
		m.weeklyBar.Width = barWidth
		m.monthlyBar.Width = barWidth

	case LoadDataMsg:
		return m, m.loadGoalData()

	case DataLoadedMsg:
		m.loading = false
		m.error = msg.Error
		if msg.Error == nil {
			m.dailyProgress = msg.Daily
			m.weeklyProgress = msg.Weekly
			m.monthlyProgress = msg.Monthly
			m.streak = msg.Streak
			m.history = msg.History

			// Update progress bars
			if m.dailyProgress != nil {
				m.dailyBar.SetPercent(m.dailyProgress.Percentage / 100.0)
			}
			if m.weeklyProgress != nil {
				m.weeklyBar.SetPercent(m.weeklyProgress.Percentage / 100.0)
			}
			if m.monthlyProgress != nil {
				m.monthlyBar.SetPercent(m.monthlyProgress.Percentage / 100.0)
			}
		}

	case progress.FrameMsg:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		m.dailyBar, cmd = m.dailyBar.Update(msg)
		cmds = append(cmds, cmd)

		m.weeklyBar, cmd = m.weeklyBar.Update(msg)
		cmds = append(cmds, cmd)

		m.monthlyBar, cmd = m.monthlyBar.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View renders the goal dashboard
func (m GoalDashboardModel) View() string {
	if m.loading {
		return "\n🍅 Loading goal data...\n"
	}

	if m.error != nil {
		return fmt.Sprintf("\n❌ Error loading goals: %v\n", m.error)
	}

	var b strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		MarginBottom(1)
	
	b.WriteString(headerStyle.Render("🎯 Pomodoro Goal Dashboard"))
	b.WriteString("\n\n")

	// Goal progress section
	b.WriteString(m.renderGoalProgress())

	// Streak section
	b.WriteString(m.renderStreak())

	if m.showHistory {
		b.WriteString(m.renderHistory())
	}

	if m.showAdjustment {
		b.WriteString(m.renderAdjustment())
	}

	// Help
	helpView := m.help.View(m.keys)
	b.WriteString("\n")
	b.WriteString(helpView)

	return b.String()
}

func (m GoalDashboardModel) renderGoalProgress() string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39"))

	progressStyle := lipgloss.NewStyle().
		MarginBottom(1)

	// Daily Progress
	if m.dailyProgress != nil {
		b.WriteString(titleStyle.Render("📅 Daily Goal"))
		b.WriteString("\n")
		b.WriteString(m.renderProgressBar(m.dailyProgress, m.dailyBar))
		b.WriteString(progressStyle.Render(""))
	}

	// Weekly Progress
	if m.weeklyProgress != nil {
		b.WriteString(titleStyle.Render("📊 Weekly Goal"))
		b.WriteString("\n")
		b.WriteString(m.renderProgressBar(m.weeklyProgress, m.weeklyBar))
		b.WriteString(progressStyle.Render(""))
	}

	// Monthly Progress
	if m.monthlyProgress != nil {
		b.WriteString(titleStyle.Render("📈 Monthly Goal"))
		b.WriteString("\n")
		b.WriteString(m.renderProgressBar(m.monthlyProgress, m.monthlyBar))
		b.WriteString(progressStyle.Render(""))
	}

	return b.String()
}

func (m GoalDashboardModel) renderProgressBar(progress *goals.GoalProgress, bar progress.Model) string {
	var b strings.Builder

	// Progress bar
	b.WriteString("  ")
	b.WriteString(bar.View())
	b.WriteString(fmt.Sprintf(" %.1f%%", progress.Percentage))
	
	// Status indicator
	if progress.IsComplete {
		if progress.IsOverAchieved {
			b.WriteString(" 🌟")
		} else {
			b.WriteString(" ✅")
		}
	}
	b.WriteString("\n")

	// Details
	b.WriteString(fmt.Sprintf("  Progress: %d/%d pomodoros", progress.Current, progress.Target))
	if progress.Remaining > 0 {
		b.WriteString(fmt.Sprintf(" (%d remaining)", progress.Remaining))
	}
	b.WriteString("\n")

	if progress.Type != goals.GoalTypeDaily && progress.RequiredPerDay > 0 {
		b.WriteString(fmt.Sprintf("  Required per day: %.1f", progress.RequiredPerDay))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	return b.String()
}

func (m GoalDashboardModel) renderStreak() string {
	if m.streak == nil {
		return ""
	}

	var b strings.Builder
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("208"))

	b.WriteString(titleStyle.Render("🔥 Streak Tracking"))
	b.WriteString("\n")

	streakStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("202")).
		Bold(true)

	if m.streak.Current > 0 {
		b.WriteString("  Current streak: ")
		b.WriteString(streakStyle.Render(fmt.Sprintf("%d days", m.streak.Current)))
		if m.streak.IsActive {
			b.WriteString(" 🔥")
		}
		b.WriteString("\n")
	} else {
		b.WriteString("  No active streak - start one today! 💪\n")
	}

	b.WriteString(fmt.Sprintf("  Best streak: %d days", m.streak.Best))
	if m.streak.Best > 0 {
		b.WriteString(" 🏆")
	}
	b.WriteString("\n\n")

	return b.String()
}

func (m GoalDashboardModel) renderHistory() string {
	if len(m.history) == 0 {
		return ""
	}

	var b strings.Builder
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("99"))

	b.WriteString(titleStyle.Render("📊 Goal History (Last 14 Days)"))
	b.WriteString("\n")

	// Show last 14 days
	start := len(m.history) - 14
	if start < 0 {
		start = 0
	}

	for i := start; i < len(m.history); i++ {
		day := m.history[i]
		indicator := "❌"
		if day.GoalMet {
			indicator = "✅"
		}

		b.WriteString(fmt.Sprintf("  %s %s %d/%d",
			indicator,
			day.Date.Format("Jan 02"),
			day.PomodoroCount,
			day.GoalTarget))

		if day.PomodoroCount > day.GoalTarget {
			b.WriteString(" 🌟")
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	return b.String()
}

func (m GoalDashboardModel) renderAdjustment() string {
	var b strings.Builder
	
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214"))

	b.WriteString(titleStyle.Render("⚙️  Adjust Goals"))
	b.WriteString("\n")
	b.WriteString("  Daily target: " + m.newDailyTarget + "\n")
	b.WriteString("  Weekly target: " + m.newWeeklyTarget + "\n")
	b.WriteString("  Press Enter to save, Esc to cancel\n\n")

	return b.String()
}

func (m GoalDashboardModel) loadGoalData() tea.Cmd {
	return func() tea.Msg {
		daily, err := m.goalManager.GetDailyGoalProgress()
		if err != nil {
			return DataLoadedMsg{Error: err}
		}

		weekly, err := m.goalManager.GetWeeklyGoalProgress()
		if err != nil {
			return DataLoadedMsg{Error: err}
		}

		monthly, err := m.goalManager.GetMonthlyGoalProgress()
		if err != nil {
			return DataLoadedMsg{Error: err}
		}

		streak, err := m.goalManager.GetStreak()
		if err != nil {
			return DataLoadedMsg{Error: err}
		}

		history, err := m.goalManager.GetGoalHistory(14)
		if err != nil {
			return DataLoadedMsg{Error: err}
		}

		return DataLoadedMsg{
			Daily:   daily,
			Weekly:  weekly,
			Monthly: monthly,
			Streak:  streak,
			History: history,
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}