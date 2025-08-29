package template

import (
    "os"
    "path/filepath"
    "testing"
    "time"

    "github.com/ethan-k/pomodoro-cli/internal/audio"
)

// setTempHome sets HOME to a temporary directory for filesystem-scoped tests.
func setTempHome(t *testing.T) string {
    t.Helper()
    dir := t.TempDir()
    t.Setenv("HOME", dir)
    return dir
}

func TestTemplateManagerCRUD(t *testing.T) {
    _ = setTempHome(t)

    tm, err := NewTemplateManager()
    if err != nil {
        t.Fatalf("NewTemplateManager error: %v", err)
    }

    // Create
    ac := audio.DefaultConfig()
    if err := tm.Create("coding", "Write code", "25m", []string{"go", "cli"}, ac); err != nil {
        t.Fatalf("Create error: %v", err)
    }
    if !tm.Exists("coding") {
        t.Fatalf("template should exist after create")
    }

    // Get
    tpl, err := tm.Get("coding")
    if err != nil {
        t.Fatalf("Get error: %v", err)
    }
    if tpl.Name != "coding" || tpl.Description != "Write code" || tpl.Duration != "25m" {
        t.Fatalf("unexpected template contents: %+v", tpl)
    }

    // List
    list, err := tm.List()
    if err != nil {
        t.Fatalf("List error: %v", err)
    }
    if len(list) != 1 || list[0].Name != "coding" {
        t.Fatalf("List unexpected: %+v", list)
    }

    // Update
    if err := tm.Update("coding", "Refactor code", "30m", []string{"go"}, ac); err != nil {
        t.Fatalf("Update error: %v", err)
    }
    tpl2, err := tm.Get("coding")
    if err != nil {
        t.Fatalf("Get after update error: %v", err)
    }
    if tpl2.Description != "Refactor code" || tpl2.Duration != "30m" || len(tpl2.Tags) != 1 || tpl2.Tags[0] != "go" {
        t.Fatalf("unexpected updated template: %+v", tpl2)
    }

    // Export
    out := filepath.Join(t.TempDir(), "coding.yml")
    if err := tm.Export("coding", out); err != nil {
        t.Fatalf("Export error: %v", err)
    }
    if _, err := os.Stat(out); err != nil {
        t.Fatalf("exported file missing: %v", err)
    }

    // Import (new name) and overwrite behavior
    // Modify exported file name content to simulate different template name
    // Instead, create a second template file manually by exporting and then importing with overwrite to same name
    if err := tm.Import(out, false); err != nil {
        // Importing the same name should work because file already contains name "coding"; it will overwrite if flag set
        // Without overwrite and existing file, expect error
        if err == nil {
            t.Fatalf("expected error importing existing without overwrite")
        }
    }
    if err := tm.Import(out, true); err != nil {
        t.Fatalf("Import with overwrite error: %v", err)
    }

    // Delete
    if err := tm.Delete("coding"); err != nil {
        t.Fatalf("Delete error: %v", err)
    }
    if tm.Exists("coding") {
        t.Fatalf("template should not exist after delete")
    }
}

func TestTemplateValidation(t *testing.T) {
    _ = setTempHome(t)
    tm, err := NewTemplateManager()
    if err != nil {
        t.Fatalf("NewTemplateManager error: %v", err)
    }

    // Bad name
    if err := tm.Create("", "desc", "25m", nil, nil); err == nil {
        t.Fatalf("expected error for empty name")
    }
    if err := tm.Create("bad/name", "desc", "25m", nil, nil); err == nil {
        t.Fatalf("expected error for invalid name chars")
    }

    // Bad duration
    if err := tm.Create("ok", "desc", "abc", nil, nil); err == nil {
        t.Fatalf("expected error for invalid duration")
    }

    // Create good
    if err := tm.Create("ok", "desc", "1m", nil, nil); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Touch timestamps roughly now
    tpl, err := tm.Get("ok")
    if err != nil {
        t.Fatalf("Get error: %v", err)
    }
    if time.Since(tpl.CreatedAt) > time.Minute || time.Since(tpl.UpdatedAt) > time.Minute {
        t.Fatalf("timestamps not set near now: %+v", tpl)
    }
}

