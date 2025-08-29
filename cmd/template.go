package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/ethan-k/pomodoro-cli/internal/audio"
	"github.com/ethan-k/pomodoro-cli/internal/db"
	"github.com/ethan-k/pomodoro-cli/internal/model"
	"github.com/ethan-k/pomodoro-cli/internal/notify"
	"github.com/ethan-k/pomodoro-cli/internal/template"
	"github.com/ethan-k/pomodoro-cli/internal/utils"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage session templates",
	Long: `Create, manage, and use session templates to quickly start predefined pomodoro sessions.

Templates allow you to save common session configurations including duration,
description, tags, and audio settings for easy reuse.`,
}

// templateCreateCmd creates a new template
var templateCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new session template",
	Long: `Create a new session template with the specified name and configuration.

You can specify duration, description, tags, and audio settings that will be
saved and can be reused when starting sessions.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		description, _ := cmd.Flags().GetString("description")
		duration, _ := cmd.Flags().GetString("duration")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		audioEnabled, _ := cmd.Flags().GetBool("audio")
		volume, _ := cmd.Flags().GetFloat64("volume")

		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		// Create audio config if specified
		var audioConfig *audio.Config
		if cmd.Flags().Changed("audio") || cmd.Flags().Changed("volume") {
			audioConfig = audio.DefaultConfig()
			audioConfig.Enabled = audioEnabled
			audioConfig.Volume = volume
		}

		if err := tm.Create(name, description, duration, tags, audioConfig); err != nil {
			return err
		}

		fmt.Printf("Template '%s' created successfully\n", name)
		return nil
	},
}

// templateListCmd lists all templates
var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all session templates",
	Long:  `List all available session templates with their details.`,
	RunE: func(_ *cobra.Command, _ []string) error {
		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		templates, err := tm.List()
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			fmt.Println("No templates found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if _, err := fmt.Fprintf(w, "NAME\tDURATION\tDESCRIPTION\tTAGS\n"); err != nil {
			return err
		}
		for _, t := range templates {
			tags := strings.Join(t.Tags, ", ")
			if len(tags) > 30 {
				tags = tags[:27] + "..."
			}
			description := t.Description
			if len(description) > 40 {
				description = description[:37] + "..."
			}
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.Name, t.Duration, description, tags); err != nil {
				return err
			}
		}
		if err := w.Flush(); err != nil {
			return err
		}

		return nil
	},
}

// templateShowCmd shows details of a specific template
var templateShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show details of a session template",
	Long:  `Display detailed information about a specific session template.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		name := args[0]

		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		template, err := tm.Get(name)
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\n", template.Name)
		fmt.Printf("Description: %s\n", template.Description)
		fmt.Printf("Duration: %s\n", template.Duration)
		if len(template.Tags) > 0 {
			fmt.Printf("Tags: %s\n", strings.Join(template.Tags, ", "))
		}
		if template.Audio != nil {
			fmt.Printf("Audio Enabled: %t\n", template.Audio.Enabled)
			if template.Audio.Enabled {
				fmt.Printf("Volume: %.1f\n", template.Audio.Volume)
			}
		}
		fmt.Printf("Created: %s\n", template.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", template.UpdatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// templateUpdateCmd updates an existing template
var templateUpdateCmd = &cobra.Command{
	Use:   "update <name>",
	Short: "Update an existing session template",
	Long:  `Update an existing session template with new configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		// Get existing template first
		existing, err := tm.Get(name)
		if err != nil {
			return err
		}

		// Use existing values if flags not provided
		description := existing.Description
		duration := existing.Duration
		tags := existing.Tags
		audioConfig := existing.Audio

		// Update with flag values if provided
		if cmd.Flags().Changed("description") {
			description, _ = cmd.Flags().GetString("description")
		}
		if cmd.Flags().Changed("duration") {
			duration, _ = cmd.Flags().GetString("duration")
		}
		if cmd.Flags().Changed("tags") {
			tags, _ = cmd.Flags().GetStringSlice("tags")
		}
		if cmd.Flags().Changed("audio") || cmd.Flags().Changed("volume") {
			if audioConfig == nil {
				audioConfig = audio.DefaultConfig()
			}
			if cmd.Flags().Changed("audio") {
				audioEnabled, _ := cmd.Flags().GetBool("audio")
				audioConfig.Enabled = audioEnabled
			}
			if cmd.Flags().Changed("volume") {
				volume, _ := cmd.Flags().GetFloat64("volume")
				audioConfig.Volume = volume
			}
		}

		if err := tm.Update(name, description, duration, tags, audioConfig); err != nil {
			return err
		}

		fmt.Printf("Template '%s' updated successfully\n", name)
		return nil
	},
}

// templateDeleteCmd deletes a template
var templateDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a session template",
	Long:  `Delete an existing session template. This action cannot be undone.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Printf("Are you sure you want to delete template '%s'? (y/N): ", name)
			var response string
			if _, err := fmt.Scanln(&response); err != nil {
				fmt.Println("Input error; deletion cancelled")
				return nil
			}
			if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
				fmt.Println("Deletion cancelled")
				return nil
			}
		}

		if err := tm.Delete(name); err != nil {
			return err
		}

		fmt.Printf("Template '%s' deleted successfully\n", name)
		return nil
	},
}

// templateStartCmd starts a session from a template
var templateStartCmd = &cobra.Command{
	Use:   "start <name>",
	Short: "Start a session from a template",
	Long: `Start a pomodoro session using the configuration from the specified template.
	
This will load the template's duration, tags, and audio settings and start
a new session with those parameters.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		template, err := tm.Get(name)
		if err != nil {
			return err
		}

		fmt.Printf("Starting session from template '%s'...\n", name)
		return runTemplateStart(cmd, template)
	},
}

// templateExportCmd exports a template to a file
var templateExportCmd = &cobra.Command{
	Use:   "export <name> <output-file>",
	Short: "Export a template to a file",
	Long:  `Export a session template to a YAML file that can be shared or imported later.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		name := args[0]
		outputPath := args[1]

		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		if err := tm.Export(name, outputPath); err != nil {
			return err
		}

		fmt.Printf("Template '%s' exported to '%s'\n", name, outputPath)
		return nil
	},
}

// templateImportCmd imports a template from a file
var templateImportCmd = &cobra.Command{
	Use:   "import <template-file>",
	Short: "Import a template from a file",
	Long:  `Import a session template from a YAML file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templatePath := args[0]
		overwrite, _ := cmd.Flags().GetBool("overwrite")

		tm, err := template.NewTemplateManager()
		if err != nil {
			return fmt.Errorf("error initializing template manager: %v", err)
		}

		if err := tm.Import(templatePath, overwrite); err != nil {
			return err
		}

		fmt.Printf("Template imported from '%s'\n", templatePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	// Add subcommands
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateUpdateCmd)
	templateCmd.AddCommand(templateDeleteCmd)
	templateCmd.AddCommand(templateStartCmd)
	templateCmd.AddCommand(templateExportCmd)
	templateCmd.AddCommand(templateImportCmd)

	// Flags for create command
	templateCreateCmd.Flags().StringP("description", "d", "", "Template description")
	templateCreateCmd.Flags().String("duration", "25m", "Session duration")
	templateCreateCmd.Flags().StringSliceP("tags", "t", nil, "Session tags")
	templateCreateCmd.Flags().Bool("audio", true, "Enable audio notifications")
	templateCreateCmd.Flags().Float64("volume", 0.5, "Audio volume (0.0-1.0)")

	// Flags for update command (same as create)
	templateUpdateCmd.Flags().StringP("description", "d", "", "Template description")
	templateUpdateCmd.Flags().String("duration", "", "Session duration")
	templateUpdateCmd.Flags().StringSliceP("tags", "t", nil, "Session tags")
	templateUpdateCmd.Flags().Bool("audio", false, "Enable audio notifications")
	templateUpdateCmd.Flags().Float64("volume", 0.0, "Audio volume (0.0-1.0)")

	// Flags for delete command
	templateDeleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")

	// Flags for start command
	templateStartCmd.Flags().String("duration", "", "Override template duration")
	templateStartCmd.Flags().StringSliceP("tags", "t", nil, "Override template tags")
	templateStartCmd.Flags().StringP("message", "m", "", "Override template description")

	// Flags for import command
	templateImportCmd.Flags().Bool("overwrite", false, "Overwrite existing template")
}

// runTemplateStart runs a pomodoro session from a template
func runTemplateStart(cmd *cobra.Command, tmpl *template.Template) error {
	// Parse template duration
	templateDuration, err := time.ParseDuration(tmpl.Duration)
	if err != nil {
		return fmt.Errorf("invalid duration in template: %v", err)
	}

	// Get template values
	desc := tmpl.Description
	templateTags := tmpl.Tags

	// Override with command line flags if provided
	if cmd.Flags().Changed("duration") {
		durationStr, _ := cmd.Flags().GetString("duration")
		templateDuration, err = time.ParseDuration(durationStr)
		if err != nil {
			return fmt.Errorf("invalid duration: %v", err)
		}
	}
	if cmd.Flags().Changed("tags") {
		templateTags, _ = cmd.Flags().GetStringSlice("tags")
	}
	if cmd.Flags().Changed("message") {
		desc, _ = cmd.Flags().GetString("message")
	}

	// Local variables for this session (no global state dependencies)
	sessionDesc := desc
	sessionTags := templateTags
	sessionDuration := templateDuration
	sessionNoWait := false
	sessionAgo := time.Duration(0)
	sessionJSONOutput := false
	sessionSilentMode := false

	// Validate and sanitize inputs (same logic as start.go)
	sessionDesc = utils.SanitizeDescription(sessionDesc)
	if err := utils.ValidateDescription(sessionDesc, false); err != nil {
		return fmt.Errorf("invalid description: %v", err)
	}

	if err := utils.ValidateDuration(sessionDuration); err != nil {
		return fmt.Errorf("invalid duration: %v", err)
	}

	sessionTags = utils.SanitizeTags(sessionTags)
	if err := utils.ValidateTags(sessionTags); err != nil {
		return fmt.Errorf("invalid tags: %v", err)
	}

	startTime := time.Now().Add(-sessionAgo)
	endTime := startTime.Add(sessionDuration)

	database, err := db.NewDB()
	if err != nil {
		return err
	}
	defer func() {
		if err := database.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: closing DB: %v\n", err)
		}
	}()

	tagsCSV := strings.Join(sessionTags, ",")
	id, err := database.CreateSession(
		startTime,
		endTime,
		sessionDesc,
		int64(sessionDuration.Seconds()),
		tagsCSV,
		false,
	)
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}

	if sessionJSONOutput {
		fmt.Printf(`{"id":%d,"description":"%s","duration":"%s","end_time":"%s"}`+"\n",
			id, sessionDesc, sessionDuration, endTime.Format(time.RFC3339))
		return nil
	}

	if sessionNoWait {
		fmt.Printf("Started Pomodoro ID %d: %s for %s (running in background)\n", id, sessionDesc, sessionDuration)
		return nil
	}

	p := model.NewPomodoroModel(id, sessionDesc, startTime, sessionDuration, false)

	if _, err := tea.NewProgram(p).Run(); err != nil {
		return fmt.Errorf("error running UI: %v", err)
	}

	if err := notify.NotifyPomodoroCompleteWithOptions(sessionDesc, sessionSilentMode); err != nil {
		return fmt.Errorf("error sending notification: %v", err)
	}

	return nil
}
