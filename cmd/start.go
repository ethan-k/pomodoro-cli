package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ethan-k/pomodoro-cli/internal/db"
	"github.com/ethan-k/pomodoro-cli/internal/model"
	"github.com/ethan-k/pomodoro-cli/internal/notify"
	"github.com/ethan-k/pomodoro-cli/internal/utils"
)

var (
	description    string
	tags           []string
	duration       time.Duration
	wait           bool
	noWait         bool
	ago            time.Duration
	jsonOutput     bool
	silentMode     bool
	continuousMode bool
)

var startCmd = &cobra.Command{
	Use:   "start [description]",
	Short: "Starts a new Pomodoro session",
	Long: `Starts a new Pomodoro timer.

You can optionally provide a description for the session.
Use flags to specify tags, duration, or if the timer should block.

Example:
  pomodoro start "Refactor API" -t coding,backend --duration 50m`,
	Aliases: []string{"s"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			description = args[0]
		}

		// Validate and sanitize inputs
		description = utils.SanitizeDescription(description)
		if err := utils.ValidateDescription(description, false); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid description: %v\n", err)
			os.Exit(1)
		}

		if err := utils.ValidateDuration(duration); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid duration: %v\n", err)
			os.Exit(1)
		}

		tags = utils.SanitizeTags(tags)
		if err := utils.ValidateTags(tags); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid tags: %v\n", err)
			os.Exit(1)
		}
		startTime := time.Now().Add(-ago)
		endTime := startTime.Add(duration)

		database, err := db.NewDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer database.Close()

		tagsCSV := strings.Join(tags, ",")
		id, err := database.CreateSession(
			startTime,
			endTime,
			description,
			int64(duration.Seconds()),
			tagsCSV,
			false,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating session: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			fmt.Printf(`{"id":%d,"description":"%s","duration":"%s","end_time":"%s"}`+"\n",
				id, description, duration, endTime.Format(time.RFC3339))
			return
		}

		if noWait {
			fmt.Printf("Started Pomodoro ID %d: %s for %s (running in background)\n", id, description, duration)
			return
		}

		p := model.NewPomodoroModel(id, description, startTime, duration, false)

		if _, err := tea.NewProgram(p).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
			os.Exit(1)
		}

		if err := notify.NotifyPomodoroCompleteWithOptions(description, silentMode); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
		}

		// Continuous mode: prompt for next action
		if continuousMode {
			handleContinuousMode()
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringSliceVarP(&tags, "tags", "t", []string{}, "Comma-separated tags for the session (e.g., coding,backend)")
	startCmd.Flags().DurationVarP(&duration, "duration", "d", 25*time.Minute, "Duration of the Pomodoro session (e.g., 25m, 1h)")
	startCmd.Flags().BoolVar(&noWait, "no-wait", false, "Run in background without showing progress bar")
	startCmd.Flags().DurationVar(&ago, "ago", 0, "Start the Pomodoro as if it began some time ago (e.g., 5m)")
	startCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format (for non-TTY usage)")
	startCmd.Flags().BoolVar(&silentMode, "silent", false, "Disable audio notifications for this session")
	startCmd.Flags().BoolVar(&continuousMode, "continuous", false, "Stay in program after session completion for next action")
}

// handleContinuousMode prompts user for next action after session completion
func handleContinuousMode() {
	fmt.Println("\n🍅 Session completed! What would you like to do next?")
	fmt.Println("1. Start a break (b)")
	fmt.Println("2. Start another pomodoro (p)")
	fmt.Println("3. View status (s)")
	fmt.Println("4. Quit (q)")
	fmt.Print("\nChoose an option: ")

	var choice string
	fmt.Scanln(&choice)

	switch strings.ToLower(choice) {
	case "1", "b", "break":
		fmt.Println("Starting break...")
		// Use the break command implementation
		runBreakSession(5*time.Minute, false)
	case "2", "p", "pomodoro":
		fmt.Println("Starting another pomodoro...")
		// Restart with same settings
		runPomodoroSession()
	case "3", "s", "status":
		fmt.Println("Checking status...")
		// Could show today's stats or goals
		showQuickStatus()
	case "4", "q", "quit":
		fmt.Println("Goodbye! 👋")
		return
	default:
		fmt.Println("Invalid option. Goodbye! 👋")
	}
}

// runBreakSession runs a break session with specified duration
func runBreakSession(duration time.Duration, wait bool) {
	startTime := time.Now()
	endTime := startTime.Add(duration)

	database, err := db.NewDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer database.Close()

	id, err := database.CreateSession(startTime, endTime, "Break", int64(duration.Seconds()), "", true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating break session: %v\n", err)
		return
	}

	if !wait {
		fmt.Printf("Started break for %s\n", duration)
		return
	}

	p := model.NewPomodoroModel(id, "Break Time", startTime, duration, true)
	if _, err := tea.NewProgram(p).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
		return
	}

	if err := notify.NotifyBreakCompleteWithOptions(silentMode); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
	}

	// Continue the continuous mode loop
	if continuousMode {
		handleContinuousMode()
	}
}

// runPomodoroSession runs another pomodoro with the same settings
func runPomodoroSession() {
	startTime := time.Now().Add(-ago)
	endTime := startTime.Add(duration)

	database, err := db.NewDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer database.Close()

	tagsCSV := strings.Join(tags, ",")
	id, err := database.CreateSession(startTime, endTime, description, int64(duration.Seconds()), tagsCSV, false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating session: %v\n", err)
		return
	}

	p := model.NewPomodoroModel(id, description, startTime, duration, false)
	if _, err := tea.NewProgram(p).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
		return
	}

	if err := notify.NotifyPomodoroCompleteWithOptions(description, silentMode); err != nil {
		fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
	}

	// Continue the continuous mode loop
	if continuousMode {
		handleContinuousMode()
	}
}

// showQuickStatus shows a quick overview of today's progress
func showQuickStatus() {
	database, err := db.NewDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer database.Close()

	sessions, err := database.GetTodaySessions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting today's sessions: %v\n", err)
		return
	}

	pomodoroCount := 0
	breakCount := 0
	for _, session := range sessions {
		if session.WasBreak {
			breakCount++
		} else {
			pomodoroCount++
		}
	}

	fmt.Printf("\n📊 Today's Progress:\n")
	fmt.Printf("🍅 Pomodoros: %d\n", pomodoroCount)
	fmt.Printf("☕ Breaks: %d\n", breakCount)
	fmt.Printf("📈 Total sessions: %d\n", len(sessions))
	fmt.Println()
}
