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
)

var (
	statusFormat string
	statusWait   bool
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the status of the current Pomodoro session",
	Long: `Shows the status of the current Pomodoro session if one is active.

You can use the --format flag to customize the output using placeholders:
  %d  - Description
  %r  - Remaining time (MM:SS)
  %p  - Progress percentage
  %t  - Tags
  %e  - End time

Example:
  pomodoro status --format "%r remaining for %d"
  pomodoro status --wait (to show a live progress bar)`,
	Run: func(cmd *cobra.Command, args []string) {
		// Connect to database
		database, err := db.NewDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer database.Close()

		// Get active session
		session, err := database.GetActiveSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting active session: %v\n", err)
			os.Exit(1)
		}

		if session == nil {
			if jsonOutput {
				fmt.Println(`{"active":false}`)
			} else {
				fmt.Println("No active Pomodoro session.")
			}
			return
		}

		// If waiting, show progress bar
		if statusWait {
			duration := session.EndTime.Sub(session.StartTime)
			elapsed := time.Since(session.StartTime)
			remaining := duration - elapsed

			if remaining <= 0 {
				fmt.Println("Session already completed.")
				return
			}

			p := model.NewPomodoroModel(
				session.ID,
				session.Description,
				session.StartTime,
				duration,
				session.WasBreak,
			)

			if _, err := tea.NewProgram(p).Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error running UI: %v\n", err)
				os.Exit(1)
			}
			return
		}

		// JSON output
		if jsonOutput {
			now := time.Now()
			remaining := session.EndTime.Sub(now).Round(time.Second)
			totalDuration := session.EndTime.Sub(session.StartTime)
			progress := float64(time.Since(session.StartTime)) / float64(totalDuration) * 100

			fmt.Printf(`{"active":true,"id":%d,"description":"%s","remaining":"%s","progress":%.1f,"end_time":"%s","tags_csv":"%s","is_break":%t}`+"\n",
				session.ID,
				session.Description,
				remaining,
				progress,
				session.EndTime.Format(time.RFC3339),
				session.TagsCSV,
				session.WasBreak)
			return
		}

		// Format output
		now := time.Now()
		remaining := session.EndTime.Sub(now).Round(time.Second)
		totalDuration := session.EndTime.Sub(session.StartTime)
		progress := float64(time.Since(session.StartTime)) / float64(totalDuration) * 100

		output := statusFormat
		output = strings.ReplaceAll(output, "%d", session.Description)
		output = strings.ReplaceAll(output, "%r", formatDuration(remaining))
		output = strings.ReplaceAll(output, "%p", fmt.Sprintf("%.1f%%", progress))
		output = strings.ReplaceAll(output, "%t", session.TagsCSV)
		output = strings.ReplaceAll(output, "%e", session.EndTime.Format("15:04:05"))

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Define flags for the status command
	statusCmd.Flags().StringVarP(&statusFormat, "format", "f", "%r remaining for %d", "Format string for status output")
	statusCmd.Flags().BoolVarP(&statusWait, "wait", "w", false, "Wait and show live progress")
	statusCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format (for non-TTY usage)")
}

// formatDuration formats a duration in MM:SS format
func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
