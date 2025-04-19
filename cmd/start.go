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
)

var (
	description string
	tags        []string
	duration    time.Duration
	wait        bool
	noWait      bool
	ago         time.Duration
	jsonOutput  bool
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

		if err := notify.NotifyPomodoroComplete(description); err != nil {
			fmt.Fprintf(os.Stderr, "Error sending notification: %v\n", err)
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
}
