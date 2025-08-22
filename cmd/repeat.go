package cmd

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ethan-k/pomodoro-cli/internal/db"
	"github.com/ethan-k/pomodoro-cli/internal/model"
	"github.com/ethan-k/pomodoro-cli/internal/notify"
)

var (
	repeatWait bool
)

// repeatCmd represents the repeat command
var repeatCmd = &cobra.Command{
	Use:   "repeat",
	Short: "Repeats the last Pomodoro session",
	Long: `Repeats the most recently completed Pomodoro session with the same parameters.

This is useful when you want to continue working on the same task.
Use the --wait flag to keep the timer running in the terminal.

Example:
  pomodoro repeat --wait`,
	Aliases: []string{"r"},
	Run: func(_ *cobra.Command, _ []string) {
		// Connect to database
		database, err := db.NewDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer func() {
			if err := database.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing database: %v\n", err)
			}
		}()

		// Get last session
		lastSession, err := database.GetLastSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting last session: %v\n", err)
			os.Exit(1)
		}

		if lastSession == nil {
			fmt.Println("No previous Pomodoro session found to repeat.")
			return
		}

		// Start a new session with the same parameters
		duration := time.Duration(lastSession.DurationSec) * time.Second
		startTime := time.Now()
		endTime := startTime.Add(duration)

		// Create session in database
		id, err := database.CreateSession(
			startTime,
			endTime,
			lastSession.Description,
			lastSession.DurationSec,
			lastSession.TagsCSV,
			lastSession.WasBreak,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating session: %v\n", err)
			os.Exit(1)
		}

		// If JSON output is requested, just print the session info and exit
		if jsonOutput {
			fmt.Printf(`{"id":%d,"description":"%s","duration":"%s","end_time":"%s","repeated":true}`+"\n",
				id, lastSession.Description, duration, endTime.Format(time.RFC3339))
			return
		}

		// Print basic info if not waiting
		if !repeatWait {
			fmt.Printf("Started repeated Pomodoro ID %d: %s for %s\n",
				id, lastSession.Description, duration)
			return
		}

		// Create and run the TUI model if waiting
		p := model.NewPomodoroModel(
			id,
			lastSession.Description,
			startTime,
			duration,
			lastSession.WasBreak,
		)

		// Run the TUI program
		if _, err := tea.NewProgram(p).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
			os.Exit(1)
		}

		// Send notification when complete
		if lastSession.WasBreak {
			if err := notify.NotifyBreakComplete(); err != nil {
				fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
			}
		} else {
			if err := notify.NotifyPomodoroComplete(lastSession.Description); err != nil {
				fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(repeatCmd)

	// Define flags for the repeat command
	repeatCmd.Flags().BoolVarP(&repeatWait, "wait", "w", false, "Wait for the Pomodoro session to complete before exiting")
	repeatCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format (for non-TTY usage)")
}
