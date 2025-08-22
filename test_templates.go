// +build ignore

// This is a simple test file to verify template functionality
package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/ethan-k/pomodoro-cli/internal/template"
)

func main() {
	// Test template creation
	tm, err := template.NewTemplateManager()
	if err != nil {
		log.Fatal("Error creating template manager:", err)
	}

	fmt.Println("Template manager created successfully")
	fmt.Println("Templates directory:", tm.GetTemplatesDir())

	// Test template creation
	err = tm.Create("coding", "Deep work coding session", "50m", []string{"coding", "focus"}, nil)
	if err != nil {
		log.Fatal("Error creating template:", err)
	}
	fmt.Println("Template 'coding' created successfully")

	// Test template retrieval
	tmpl, err := tm.Get("coding")
	if err != nil {
		log.Fatal("Error getting template:", err)
	}
	fmt.Printf("Template: %+v\n", tmpl)

	// Test template listing
	templates, err := tm.List()
	if err != nil {
		log.Fatal("Error listing templates:", err)
	}
	fmt.Printf("Found %d templates\n", len(templates))

	// Clean up
	err = tm.Delete("coding")
	if err != nil {
		log.Fatal("Error deleting template:", err)
	}
	fmt.Println("Template 'coding' deleted successfully")

	fmt.Println("All template tests passed!")
}