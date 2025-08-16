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
	resumeWait bool
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resumes a paused session",
	Long: `Resumes the most recently paused Pomodoro or break session.

The session will continue from where it was paused, and the new end time
will be adjusted to account for the paused duration.

Example:
  pomodoro resume
  pomodoro resume --wait`,
	Run: func(cmd *cobra.Command, args []string) {
		database, err := db.NewDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer database.Close()

		// Get paused session
		session, err := database.GetPausedSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting paused session: %v\n", err)
			os.Exit(1)
		}

		if session == nil {
			fmt.Println("No paused session to resume.")
			return
		}

		// Calculate new end time
		now := time.Now()

		// Original duration minus already elapsed time when paused
		originalDuration := time.Duration(session.DurationSec) * time.Second
		elapsedWhenPaused := session.PausedAt.Sub(session.StartTime)
		remainingDuration := originalDuration - elapsedWhenPaused

		newEndTime := now.Add(remainingDuration)

		// Resume the session
		if err := database.ResumeSession(session.ID, newEndTime); err != nil {
			fmt.Fprintf(os.Stderr, "Error resuming session: %v\n", err)
			os.Exit(1)
		}

		if jsonOutput {
			fmt.Printf(`{"id":%d,"description":"%s","status":"resumed","new_end_time":"%s","remaining_duration":"%s"}`+"\n",
				session.ID, session.Description, newEndTime.Format(time.RFC3339), remainingDuration)
			return
		}

		fmt.Printf("▶️  Resumed session: %s\n", session.Description)
		fmt.Printf("Time remaining: %s\n", remainingDuration.Round(time.Second))

		// If wait flag is set, show the progress bar
		if resumeWait {
			p := model.NewPomodoroModel(session.ID, session.Description, now, remainingDuration, session.WasBreak)

			if _, err := tea.NewProgram(p).Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
				os.Exit(1)
			}

			// Send completion notification
			if session.WasBreak {
				if err := notify.NotifyBreakComplete(); err != nil {
					fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
				}
			} else {
				if err := notify.NotifyPomodoroComplete(session.Description); err != nil {
					fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
	resumeCmd.Flags().BoolVarP(&resumeWait, "wait", "w", false, "Wait and show progress bar after resuming")
	resumeCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}
