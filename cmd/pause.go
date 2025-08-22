package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ethan-k/pomodoro-cli/internal/db"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pauses the current active session",
	Long: `Pauses the currently active Pomodoro or break session.

You can resume the session later using the 'resume' command.
The paused time will not count toward the session duration.

Example:
  pomodoro pause`,
	Run: func(_ *cobra.Command, _ []string) {
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

		// Get active session
		session, err := database.GetActiveSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting active session: %v\n", err)
			os.Exit(1)
		}

		if session == nil {
			fmt.Println("No active session to pause.")
			return
		}

		if session.IsPaused {
			fmt.Printf("Session '%s' is already paused.\n", session.Description)
			return
		}

		// Pause the session
		now := time.Now()
		if err := database.PauseSession(session.ID, now); err != nil {
			fmt.Fprintf(os.Stderr, "Error pausing session: %v\n", err)
			os.Exit(1)
		}

		// if jsonOutput {
		// fmt.Printf(`{"id":%d,"description":"%s","status":"paused","paused_at":"%s"}`+"\n",
		// session.ID, session.Description, now.Format(time.RFC3339))
		// return
		// }

		fmt.Printf("⏸️  Paused session: %s\n", session.Description)
		fmt.Println("Use 'pomodoro resume' to continue.")
	},
}

func init() {
	rootCmd.AddCommand(pauseCmd)
	pauseCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
}
