package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/ethan-k/pomodoro-cli/internal/config"
	"github.com/ethan-k/pomodoro-cli/internal/db"
	"github.com/ethan-k/pomodoro-cli/internal/goals"
	"github.com/ethan-k/pomodoro-cli/internal/model"
)

var (
	goalsOutputJSON bool
	goalsShowDaily  bool
	goalsShowWeekly bool
	goalsShowStreak bool
	goalsShowHistory bool
	goalsHistoryDays int
	goalsSetDaily   int
	goalsSetWeekly  int
)

var goalsCmd = &cobra.Command{
	Use:   "goals",
	Short: "View and manage pomodoro goals",
	Long: `View your pomodoro goal progress, streaks, and performance analytics.
	
The goals command provides an interactive dashboard showing:
- Daily, weekly, and monthly goal progress
- Visual progress bars and completion status
- Streak tracking and achievements
- Historical performance analysis
- Goal adjustment capabilities`,
	Example: `  # Show interactive goal dashboard
  pomodoro goals

  # Show goals as JSON
  pomodoro goals --json

  # Show only daily progress
  pomodoro goals --daily

  # Show goal history for last 30 days
  pomodoro goals --history --days 30

  # Set new daily goal target
  pomodoro goals --set-daily 10

  # Set new weekly goal target  
  pomodoro goals --set-weekly 50`,
	RunE: runGoalsCommand,
}

func init() {
	goalsCmd.Flags().BoolVar(&goalsOutputJSON, "json", false, "Output goals data as JSON")
	goalsCmd.Flags().BoolVar(&goalsShowDaily, "daily", false, "Show only daily goal progress")
	goalsCmd.Flags().BoolVar(&goalsShowWeekly, "weekly", false, "Show only weekly goal progress")
	goalsCmd.Flags().BoolVar(&goalsShowStreak, "streak", false, "Show only streak information")
	goalsCmd.Flags().BoolVar(&goalsShowHistory, "history", false, "Show goal history")
	goalsCmd.Flags().IntVar(&goalsHistoryDays, "days", 14, "Number of days for history (default: 14)")
	goalsCmd.Flags().IntVar(&goalsSetDaily, "set-daily", 0, "Set daily goal target")
	goalsCmd.Flags().IntVar(&goalsSetWeekly, "set-weekly", 0, "Set weekly goal target")

	rootCmd.AddCommand(goalsCmd)
}

func runGoalsCommand(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Initialize database
	database, err := db.NewDB()
	if err != nil {
		return fmt.Errorf("error initializing database: %w", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			fmt.Printf("Warning: error closing database: %v\n", err)
		}
	}()

	// Create goal manager
	goalManager := goals.NewGoalManager(database, cfg)

	// Handle goal setting
	if goalsSetDaily > 0 || goalsSetWeekly > 0 {
		return handleGoalSetting(goalManager, cfg)
	}

	// Handle specific view requests
	if goalsShowDaily {
		return showDailyGoals(goalManager)
	}
	if goalsShowWeekly {
		return showWeeklyGoals(goalManager)
	}
	if goalsShowStreak {
		return showStreakInfo(goalManager)
	}
	if goalsShowHistory {
		return showGoalHistory(goalManager)
	}

	// Show interactive dashboard or JSON output
	if goalsOutputJSON {
		return showGoalsJSON(goalManager)
	}

	return showInteractiveDashboard(goalManager)
}

func handleGoalSetting(goalManager *goals.GoalManager, cfg *config.Config) error {
	dailyTarget := cfg.Goals.DailyCount
	weeklyTarget := cfg.Goals.WeeklyCount

	if goalsSetDaily > 0 {
		dailyTarget = goalsSetDaily
	}
	if goalsSetWeekly > 0 {
		weeklyTarget = goalsSetWeekly
	}

	if err := goalManager.UpdateGoalTargets(dailyTarget, weeklyTarget); err != nil {
		return fmt.Errorf("error updating goals: %w", err)
	}

	fmt.Printf("✅ Goals updated successfully!\n")
	fmt.Printf("   Daily target: %d pomodoros\n", dailyTarget)
	fmt.Printf("   Weekly target: %d pomodoros\n", weeklyTarget)

	return nil
}

func showDailyGoals(goalManager *goals.GoalManager) error {
	progress, err := goalManager.GetDailyGoalProgress()
	if err != nil {
		return fmt.Errorf("error getting daily progress: %w", err)
	}

	if goalsOutputJSON {
		data, err := json.MarshalIndent(progress, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println("📅 Daily Goal Progress")
	fmt.Println(strings.Repeat("─", 40))
	fmt.Printf("Progress: %d/%d pomodoros (%.1f%%)\n",
		progress.Current, progress.Target, progress.Percentage)
	
	if progress.IsComplete {
		if progress.IsOverAchieved {
			fmt.Println("Status: Overachieved! 🌟")
		} else {
			fmt.Println("Status: Complete! ✅")
		}
	} else {
		fmt.Printf("Remaining: %d pomodoros\n", progress.Remaining)
	}

	return nil
}

func showWeeklyGoals(goalManager *goals.GoalManager) error {
	progress, err := goalManager.GetWeeklyGoalProgress()
	if err != nil {
		return fmt.Errorf("error getting weekly progress: %w", err)
	}

	if goalsOutputJSON {
		data, err := json.MarshalIndent(progress, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println("📊 Weekly Goal Progress")
	fmt.Println(strings.Repeat("─", 40))
	fmt.Printf("Progress: %d/%d pomodoros (%.1f%%)\n",
		progress.Current, progress.Target, progress.Percentage)
	
	if progress.IsComplete {
		if progress.IsOverAchieved {
			fmt.Println("Status: Overachieved! 🌟")
		} else {
			fmt.Println("Status: Complete! ✅")
		}
	} else {
		fmt.Printf("Remaining: %d pomodoros\n", progress.Remaining)
		if progress.RequiredPerDay > 0 {
			fmt.Printf("Required per day: %.1f\n", progress.RequiredPerDay)
		}
	}

	return nil
}

func showStreakInfo(goalManager *goals.GoalManager) error {
	streak, err := goalManager.GetStreak()
	if err != nil {
		return fmt.Errorf("error getting streak: %w", err)
	}

	if goalsOutputJSON {
		data, err := json.MarshalIndent(streak, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Println("🔥 Streak Information")
	fmt.Println(strings.Repeat("─", 40))
	
	if streak.Current > 0 {
		fmt.Printf("Current streak: %d days", streak.Current)
		if streak.IsActive {
			fmt.Print(" 🔥")
		}
		fmt.Println()
	} else {
		fmt.Println("No active streak - start one today! 💪")
	}

	fmt.Printf("Best streak: %d days", streak.Best)
	if streak.Best > 0 {
		fmt.Print(" 🏆")
	}
	fmt.Println()

	if !streak.LastActive.IsZero() {
		fmt.Printf("Last active: %s\n", streak.LastActive.Format("2006-01-02"))
	}

	return nil
}

func showGoalHistory(goalManager *goals.GoalManager) error {
	history, err := goalManager.GetGoalHistory(goalsHistoryDays)
	if err != nil {
		return fmt.Errorf("error getting history: %w", err)
	}

	if goalsOutputJSON {
		data, err := json.MarshalIndent(history, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling JSON: %w", err)
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("📊 Goal History (Last %d Days)\n", goalsHistoryDays)
	fmt.Println(strings.Repeat("─", 50))

	totalMet := 0
	for _, day := range history {
		indicator := "❌"
		if day.GoalMet {
			indicator = "✅"
			totalMet++
		}

		fmt.Printf("%s %s %d/%d pomodoros",
			indicator,
			day.Date.Format("Jan 02"),
			day.PomodoroCount,
			day.GoalTarget)

		if day.PomodoroCount > day.GoalTarget {
			fmt.Print(" 🌟")
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("─", 50))
	successRate := float64(totalMet) / float64(len(history)) * 100
	fmt.Printf("Goal success rate: %.1f%% (%d/%d days)\n",
		successRate, totalMet, len(history))

	return nil
}

func showGoalsJSON(goalManager *goals.GoalManager) error {
	// Gather all goal data
	daily, err := goalManager.GetDailyGoalProgress()
	if err != nil {
		return fmt.Errorf("error getting daily progress: %w", err)
	}

	weekly, err := goalManager.GetWeeklyGoalProgress()
	if err != nil {
		return fmt.Errorf("error getting weekly progress: %w", err)
	}

	monthly, err := goalManager.GetMonthlyGoalProgress()
	if err != nil {
		return fmt.Errorf("error getting monthly progress: %w", err)
	}

	streak, err := goalManager.GetStreak()
	if err != nil {
		return fmt.Errorf("error getting streak: %w", err)
	}

	history, err := goalManager.GetGoalHistory(goalsHistoryDays)
	if err != nil {
		return fmt.Errorf("error getting history: %w", err)
	}

	// Create combined output
	output := map[string]interface{}{
		"daily":   daily,
		"weekly":  weekly,
		"monthly": monthly,
		"streak":  streak,
		"history": history,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}

func showInteractiveDashboard(goalManager *goals.GoalManager) error {
	dashboardModel := model.NewGoalDashboardModel(goalManager)

	p := tea.NewProgram(dashboardModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running dashboard: %w", err)
	}

	return nil
}