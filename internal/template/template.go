package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethan-k/pomodoro-cli/internal/audio"
	"github.com/ethan-k/pomodoro-cli/internal/utils"
	"gopkg.in/yaml.v3"
)

// Template represents a session template
type Template struct {
	Name        string        `yaml:"name"`
	Description string        `yaml:"description"`
	Duration    string        `yaml:"duration"`
	Tags        []string      `yaml:"tags,omitempty"`
	Audio       *audio.Config `yaml:"audio,omitempty"`
	CreatedAt   time.Time     `yaml:"created_at"`
	UpdatedAt   time.Time     `yaml:"updated_at"`
}

// TemplateManager handles template operations
type TemplateManager struct {
	templatesDir string
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() (*TemplateManager, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home dir: %v", err)
	}

	templatesDir := filepath.Join(home, ".config", "pomodoro", "templates")
	if err := os.MkdirAll(templatesDir, 0750); err != nil {
		return nil, fmt.Errorf("error creating templates directory: %v", err)
	}

	return &TemplateManager{
		templatesDir: templatesDir,
	}, nil
}

// Create creates a new template
func (tm *TemplateManager) Create(name, description, duration string, tags []string, audioConfig *audio.Config) error {
	if err := tm.validateTemplateName(name); err != nil {
		return err
	}

	if err := tm.validateDuration(duration); err != nil {
		return err
	}

	// Check if template already exists
	if tm.Exists(name) {
		return fmt.Errorf("template '%s' already exists", name)
	}

	template := &Template{
		Name:        name,
		Description: description,
		Duration:    duration,
		Tags:        tags,
		Audio:       audioConfig,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return tm.save(template)
}

// Get retrieves a template by name
func (tm *TemplateManager) Get(name string) (*Template, error) {
	if err := tm.validateTemplateName(name); err != nil {
		return nil, err
	}

	templatePath := filepath.Join(tm.templatesDir, name+".yml")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("template '%s' not found", name)
		}
		return nil, fmt.Errorf("error reading template '%s': %v", name, err)
	}

	var template Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("error parsing template '%s': %v", name, err)
	}

	return &template, nil
}

// List returns all available templates
func (tm *TemplateManager) List() ([]*Template, error) {
	files, err := os.ReadDir(tm.templatesDir)
	if err != nil {
		return nil, fmt.Errorf("error reading templates directory: %v", err)
	}

	var templates []*Template
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}

		name := strings.TrimSuffix(file.Name(), ".yml")
		template, err := tm.Get(name)
		if err != nil {
			// Skip invalid templates but log the error
			fmt.Printf("Warning: skipping invalid template '%s': %v\n", name, err)
			continue
		}
		templates = append(templates, template)
	}

	return templates, nil
}

// Update updates an existing template
func (tm *TemplateManager) Update(name, description, duration string, tags []string, audioConfig *audio.Config) error {
	if err := tm.validateTemplateName(name); err != nil {
		return err
	}

	if err := tm.validateDuration(duration); err != nil {
		return err
	}

	// Check if template exists
	existing, err := tm.Get(name)
	if err != nil {
		return err
	}

	// Update fields
	existing.Description = description
	existing.Duration = duration
	existing.Tags = tags
	existing.Audio = audioConfig
	existing.UpdatedAt = time.Now()

	return tm.save(existing)
}

// Delete removes a template
func (tm *TemplateManager) Delete(name string) error {
	if err := tm.validateTemplateName(name); err != nil {
		return err
	}

	templatePath := filepath.Join(tm.templatesDir, name+".yml")
	if err := os.Remove(templatePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("template '%s' not found", name)
		}
		return fmt.Errorf("error deleting template '%s': %v", name, err)
	}

	return nil
}

// Exists checks if a template exists
func (tm *TemplateManager) Exists(name string) bool {
	templatePath := filepath.Join(tm.templatesDir, name+".yml")
	_, err := os.Stat(templatePath)
	return err == nil
}

// Export exports a template to a file
func (tm *TemplateManager) Export(name, outputPath string) error {
	template, err := tm.Get(name)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("error marshaling template: %v", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("error writing template file: %v", err)
	}

	return nil
}

// Import imports a template from a file
func (tm *TemplateManager) Import(templatePath string, overwrite bool) error {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template file: %v", err)
	}

	var template Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return fmt.Errorf("error parsing template file: %v", err)
	}

	// Validate the imported template
	if err := tm.validateTemplateName(template.Name); err != nil {
		return fmt.Errorf("invalid template name in import: %v", err)
	}

	if err := tm.validateDuration(template.Duration); err != nil {
		return fmt.Errorf("invalid duration in template: %v", err)
	}

	// Check if template already exists
	if tm.Exists(template.Name) && !overwrite {
		return fmt.Errorf("template '%s' already exists (use --overwrite to replace)", template.Name)
	}

	// Update timestamps
	if template.CreatedAt.IsZero() {
		template.CreatedAt = time.Now()
	}
	template.UpdatedAt = time.Now()

	return tm.save(&template)
}

// save saves a template to disk
func (tm *TemplateManager) save(template *Template) error {
	templatePath := filepath.Join(tm.templatesDir, template.Name+".yml")

	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("error marshaling template: %v", err)
	}

	if err := os.WriteFile(templatePath, data, 0644); err != nil {
		return fmt.Errorf("error writing template file: %v", err)
	}

	return nil
}

// validateTemplateName validates the template name
func (tm *TemplateManager) validateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("template name contains invalid characters")
	}

	if len(name) > 100 {
		return fmt.Errorf("template name too long (max 100 characters)")
	}

	return nil
}

// validateDuration validates the duration string
func (tm *TemplateManager) validateDuration(duration string) error {
	return utils.ValidateDurationString(duration)
}

// GetTemplatesDir returns the templates directory path
func (tm *TemplateManager) GetTemplatesDir() string {
	return tm.templatesDir
}