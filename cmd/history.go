package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ethan-k/pomodoro-cli/internal/db"
	"github.com/ethan-k/pomodoro-cli/internal/opf"
)

var (
	historyToday  bool
	historyWeek   bool
	historyFrom   string
	historyTo     string
	historyLimit  int
	historyFormat string
	historyOutput string
	historyTags   []string
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Shows your Pomodoro session history",
	Long: `Shows your Pomodoro session history.

You can filter by date range, limit the number of results, and specify the output format.

Examples:
  pomodoro history --today
  pomodoro history --week
  pomodoro history --from 2025-04-01 --to 2025-04-19
  pomodoro history --tags coding,writing
  pomodoro history --output opf > pomodoros.json
  pomodoro history --output json --limit 10`,
	Aliases: []string{"h"},
	Run: func(cmd *cobra.Command, args []string) {
		// Connect to database
		database, err := db.NewDB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer database.Close()

		var sessions []db.PomodoroSession

		// Determine date range
		now := time.Now()
		var startDate, endDate time.Time

		if historyToday {
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.Add(24 * time.Hour)
		} else if historyWeek {
			// Start from the beginning of the week (Monday)
			daysToMonday := int(now.Weekday())
			if daysToMonday == 0 { // Sunday
				daysToMonday = 6
			} else {
				daysToMonday = daysToMonday - 1
			}
			startDate = time.Date(now.Year(), now.Month(), now.Day()-daysToMonday, 0, 0, 0, 0, now.Location())
			endDate = now
		} else if historyFrom != "" || historyTo != "" {
			if historyFrom != "" {
				var parseErr error
				startDate, parseErr = time.Parse("2006-01-02", historyFrom)
				if parseErr != nil {
					fmt.Fprintf(os.Stderr, "Error parsing from date: %v\n", parseErr)
					os.Exit(1)
				}
			} else {
				// Default to 30 days ago if not specified
				startDate = now.AddDate(0, 0, -30)
			}

			if historyTo != "" {
				var parseErr error
				endDate, parseErr = time.Parse("2006-01-02", historyTo)
				if parseErr != nil {
					fmt.Fprintf(os.Stderr, "Error parsing to date: %v\n", parseErr)
					os.Exit(1)
				}
				// Include the full day
				endDate = endDate.Add(24 * time.Hour)
			} else {
				endDate = now
			}
		} else {
			// Default to today if no date range specified
			startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			endDate = startDate.Add(24 * time.Hour)
		}

		// Get sessions
		sessions, err = database.GetSessionsByDateRange(startDate, endDate)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting sessions: %v\n", err)
			os.Exit(1)
		}

		// Filter by tags if specified
		if len(historyTags) > 0 {
			var filteredSessions []db.PomodoroSession
			for _, session := range sessions {
				// Check if session has any of the specified tags
				for _, tag := range historyTags {
					if strings.Contains(session.TagsCSV, tag) {
						filteredSessions = append(filteredSessions, session)
						break
					}
				}
			}
			sessions = filteredSessions
		}

		// Limit the number of results
		if historyLimit > 0 && historyLimit < len(sessions) {
			sessions = sessions[:historyLimit]
		}

		// Handle different output formats
		switch historyOutput {
		case "opf":
			data, err := opf.ExportToJSON(sessions)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error exporting to OPF: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(data))

		case "json":
			// Convert sessions to a simple JSON format
			type jsonSession struct {
				ID          int64  `json:"id"`
				StartTime   string `json:"start_time"`
				EndTime     string `json:"end_time"`
				Description string `json:"description"`
				Duration    string `json:"duration"`
				Tags        string `json:"tags"`
				WasBreak    bool   `json:"was_break"`
			}

			jsonSessions := make([]jsonSession, 0, len(sessions))
			for _, s := range sessions {
				duration := s.EndTime.Sub(s.StartTime)
				jsonSessions = append(jsonSessions, jsonSession{
					ID:          s.ID,
					StartTime:   s.StartTime.Format(time.RFC3339),
					EndTime:     s.EndTime.Format(time.RFC3339),
					Description: s.Description,
					Duration:    duration.String(),
					Tags:        s.TagsCSV,
					WasBreak:    s.WasBreak,
				})
			}

			data, err := json.MarshalIndent(jsonSessions, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling to JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(data))

		default: // text or unspecified
			if len(sessions) == 0 {
				fmt.Println("No sessions found.")
				return
			}

			// Calculate statistics
			var totalDuration time.Duration
			pomodoroCount := 0
			breakCount := 0

			fmt.Println("Recent Pomodoro Sessions:")
			fmt.Println("-------------------------")

			for _, s := range sessions {
				duration := s.EndTime.Sub(s.StartTime)
				totalDuration += duration

				if s.WasBreak {
					breakCount++
				} else {
					pomodoroCount++
				}

				sessionType := "🍅"
				if s.WasBreak {
					sessionType = "☕"
				}

				fmt.Printf("%s %s: %s (%s) %s\n",
					s.StartTime.Format("2006-01-02 15:04"),
					sessionType,
					s.Description,
					duration.Round(time.Second),
					s.TagsCSV)
			}

			fmt.Println("\nSummary:")
			fmt.Printf("Total sessions: %d (%d pomodoros, %d breaks)\n",
				len(sessions),
				pomodoroCount,
				breakCount)
			fmt.Printf("Total time: %s\n", totalDuration.Round(time.Minute))
		}
	},
}

func init() {
	rootCmd.AddCommand(historyCmd)

	// Define flags for the history command
	historyCmd.Flags().BoolVar(&historyToday, "today", false, "Show sessions from today")
	historyCmd.Flags().BoolVar(&historyWeek, "week", false, "Show sessions from this week")
	historyCmd.Flags().StringVar(&historyFrom, "from", "", "Start date (YYYY-MM-DD)")
	historyCmd.Flags().StringVar(&historyTo, "to", "", "End date (YYYY-MM-DD)")
	historyCmd.Flags().IntVar(&historyLimit, "limit", 0, "Limit number of results")
	historyCmd.Flags().StringVar(&historyFormat, "format", "", "Format string for session output")
	historyCmd.Flags().StringVar(&historyOutput, "output", "text", "Output format (text, json, opf)")
	historyCmd.Flags().StringSliceVarP(&historyTags, "tags", "t", []string{}, "Filter by tags")
}
