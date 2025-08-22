// Package cmd contains the CLI commands for the Pomodoro timer application
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
	"github.com/ethan-k/pomodoro-cli/internal/utils"
)

var (
	breakDuration time.Duration
	breakWait     bool
	breakJSON     bool
	breakSilent   bool
)

// breakCmd represents the break command
var breakCmd = &cobra.Command{
	Use:   "break [duration]",
	Short: "Starts a break timer",
	Long: `Starts a break timer.

You can specify the duration for the break. If not provided, a default of 5 minutes will be used.
Use the --wait flag to keep the timer running in the terminal.

Example:
  pomodoro break 10m --wait`,
	Aliases: []string{"b"},
	Run: func(_ *cobra.Command, args []string) {
		// If duration is provided as argument, override flag
		if len(args) > 0 {
			var err error
			breakDuration, err = time.ParseDuration(args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing duration: %v\n", err)
				os.Exit(1)
			}
		}

		// Validate duration
		if err := utils.ValidateDuration(breakDuration); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid break duration: %v\n", err)
			os.Exit(1)
		}

		startTime := time.Now()
		endTime := startTime.Add(breakDuration)

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

		// Create break session in database
		id, err := database.CreateSession(
			startTime,
			endTime,
			"Break",
			int64(breakDuration.Seconds()),
			"",
			true, // isBreak = true
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating break session: %v\n", err)
			os.Exit(1)
		}

		// If JSON output is requested, just print the session info and exit
		if breakJSON {
			fmt.Printf(`{"id":%d,"type":"break","duration":"%s","end_time":"%s"}`+"\n",
				id, breakDuration, endTime.Format(time.RFC3339))
			return
		}

		// Print basic info if not waiting
		if !breakWait {
			fmt.Printf("Started break for %s\n", breakDuration)
			return
		}

		// Create and run the TUI model if waiting
		p := model.NewPomodoroModel(id, "Break Time", startTime, breakDuration, true)

		// Run the TUI program
		if _, err := tea.NewProgram(p).Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
			os.Exit(1)
		}

		// Send notification when complete
		if err := notify.NotifyBreakCompleteWithOptions(breakSilent); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)

	// Define flags for the break command
	breakCmd.Flags().DurationVarP(&breakDuration, "duration", "d", 5*time.Minute, "Duration of the break (e.g., 5m, 10m)")
	breakCmd.Flags().BoolVarP(&breakWait, "wait", "w", false, "Wait for the break to complete before exiting")
	breakCmd.Flags().BoolVar(&breakJSON, "json", false, "Output in JSON format (for non-TTY usage)")
	breakCmd.Flags().BoolVar(&breakSilent, "silent", false, "Disable audio notifications for this break")
}
