// Package main is the entry point for the Pomodoro CLI application
package main

import (
	"github.com/ethan-k/pomodoro-cli/cmd"
)

// Build information set by linker flags
var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, buildDate)
	cmd.Execute()
}
