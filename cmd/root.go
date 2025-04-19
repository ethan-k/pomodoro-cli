package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pomodoro",
	Short: "A minimalist macOS CLI Pomodoro timer",
	Long: `pomodoro is a friction-free terminal tool that starts a Pomodoro timer,
shows progress, saves sessions, and sends notifications.

It aims to be fast, scriptable, and visually informative.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
