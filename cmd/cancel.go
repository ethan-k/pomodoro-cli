package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/ethan-k/pomodoro-cli/internal/db"
)

// cancelCmd represents the cancel command
var cancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancels the active Pomodoro session",
	Long: `Cancels the currently active Pomodoro session.

This will update the session in the database with the current time as the end time.

Example:
  pomodoro cancel`,
	Aliases: []string{"c"},
	Run: func(cmd *cobra.Command, args []string) {
		// Connect to database
		database, err := db.NewDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer func() {
		if err := database.Close(); err != nil {
			// Log error but don't override the main error
		}
	}()

		// Get active session
		session, err := database.GetActiveSession()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting active session: %v\n", err)
			os.Exit(1)
		}

		if session == nil {
			fmt.Println("No active Pomodoro session to cancel.")
			return
		}

		// Update session end time to now
		now := time.Now()
		if err := database.UpdateSessionEndTime(session.ID, now); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating session: %v\n", err)
			os.Exit(1)
		}

		// Calculate actual duration
		actualDuration := now.Sub(session.StartTime).Round(time.Second)

		if jsonOutput {
			fmt.Printf(`{"id":%d,"description":"%s","status":"canceled","actual_duration":"%s"}`+"\n",
				session.ID, session.Description, actualDuration)
			return
		}

		// Output result
		fmt.Printf("Cancelled Pomodoro session: %s (ran for %s)\n",
			session.Description,
			actualDuration)
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)

	// Define flags for the cancel command
	cancelCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output in JSON format (for non-TTY usage)")
}
