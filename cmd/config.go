package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ethan-k/pomodoro-cli/internal/config"
)

var (
	configInit  bool
	configList  bool
	configKey   string
	configValue string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage pomodoro configuration",
	Long: `Manage pomodoro configuration.

You can initialize the config file, list all settings, or set individual values.

Examples:
  pomodoro config --init
  pomodoro config --list
  pomodoro config goals.daily_count 10
  pomodoro config defaults.pomodoro_duration 30m`,
	Run: func(_ *cobra.Command, args []string) {
		// Initialize config file
		if configInit {
			cfg := config.DefaultConfig()
			if err := config.SaveConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Configuration initialized with default values.")
			return
		}

		// Load existing config
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}

		// List all settings
		if configList || (configKey == "" && configValue == "" && len(args) == 0) {
			fmt.Println("Current Configuration:")
			fmt.Println("======================")
			fmt.Println("Goals:")
			fmt.Printf("  Daily count: %d pomodoros\n", cfg.Goals.DailyCount)
			fmt.Printf("  Weekly count: %d pomodoros\n", cfg.Goals.WeeklyCount)
			fmt.Println("Hooks:")
			fmt.Printf("  Enabled: %v\n", cfg.Hooks.Enabled)
			fmt.Printf("  Path: %s\n", cfg.Hooks.Path)
			fmt.Println("Defaults:")
			fmt.Printf("  Pomodoro duration: %s\n", cfg.Defaults.PomodoroDuration)
			fmt.Printf("  Break duration: %s\n", cfg.Defaults.BreakDuration)
			fmt.Printf("  Long break duration: %s\n", cfg.Defaults.LongBreakDuration)
			fmt.Println("Paths:")
			fmt.Printf("  Database: %s\n", cfg.DataPaths.Database)
			fmt.Printf("  OPF export: %s\n", cfg.DataPaths.OPFExport)
			return
		}

		// Set a configuration value
		if len(args) == 2 {
			configKey = args[0]
			configValue = args[1]
		}

		if configKey != "" && configValue != "" {
			switch configKey {
			case "goals.daily_count":
				count, err := strconv.Atoi(configValue)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Invalid value for daily count: %v\n", err)
					os.Exit(1)
				}
				cfg.Goals.DailyCount = count
			case "goals.weekly_count":
				count, err := strconv.Atoi(configValue)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Invalid value for weekly count: %v\n", err)
					os.Exit(1)
				}
				cfg.Goals.WeeklyCount = count
			case "hooks.enabled":
				enabled, err := strconv.ParseBool(configValue)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Invalid value for hooks enabled: %v\n", err)
					os.Exit(1)
				}
				cfg.Hooks.Enabled = enabled
			case "hooks.path":
				cfg.Hooks.Path = configValue
			case "defaults.pomodoro_duration":
				cfg.Defaults.PomodoroDuration = configValue
			case "defaults.break_duration":
				cfg.Defaults.BreakDuration = configValue
			case "defaults.long_break_duration":
				cfg.Defaults.LongBreakDuration = configValue
			case "paths.database":
				cfg.DataPaths.Database = configValue
			case "paths.opf_export":
				cfg.DataPaths.OPFExport = configValue
			default:
				fmt.Fprintf(os.Stderr, "Unknown configuration key: %s\n", configKey)
				os.Exit(1)
			}

			if err := config.SaveConfig(cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Configuration updated: %s = %s\n", configKey, configValue)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Define flags for the config command
	configCmd.Flags().BoolVar(&configInit, "init", false, "Initialize config file with default values")
	configCmd.Flags().BoolVar(&configList, "list", false, "List all configuration values")
	configCmd.Flags().StringVar(&configKey, "key", "", "Configuration key to set")
	configCmd.Flags().StringVar(&configValue, "value", "", "Configuration value to set")
}
