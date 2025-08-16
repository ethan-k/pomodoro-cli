# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a minimalist macOS CLI Pomodoro timer built in Go using Cobra CLI framework and Bubble Tea TUI. The application tracks Pomodoro sessions in a SQLite database and sends native notifications.

## Common Development Commands

### Building and Testing
- `make` or `make all` - Run tests and build binary
- `make build` - Build for current platform (outputs to bin/pomodoro)
- `make build-all` - Build for all supported platforms (macOS/Linux, amd64/arm64)
- `make test` - Run all tests
- `make coverage` - Generate test coverage report
- `go test ./cmd` - Run tests for specific package
- `go test -run TestBreakCommand` - Run specific test

### Code Quality
- `make fmt` - Format code using go fmt
- `make lint` - Run golint (install with: go install golang.org/x/lint/golint@latest)
- `make vet` - Run go vet for static analysis

### Installation
- `make install` - Install binary to $GOPATH/bin
- `make uninstall` - Remove binary from $GOPATH/bin

## Architecture Overview

### Core Components

**CLI Layer (`cmd/`)**
- `root.go` - Main Cobra root command and app initialization
- `start.go` - Pomodoro timer start command with TUI integration
- `break.go` - Break timer command 
- `status.go`, `history.go`, `cancel.go` - Status and management commands
- `config.go` - Configuration management commands

**Database Layer (`internal/db/`)**
- `db.go` - SQLite database interface and operations
- Stores sessions in `~/.local/share/pomodoro/history.db`
- Interface-based design for testability

**TUI Layer (`internal/model/`)**
- `progress.go` - Bubble Tea model for animated progress bars
- Uses Charm Bracelet libraries (bubbletea, lipgloss, bubbles)
- Different color schemes for Pomodoro vs break timers

**Notification Layer (`internal/notify/`)**
- `notify.go` - Cross-platform desktop notifications using beeep

**Configuration (`internal/config/`)**
- `config.go` - YAML-based configuration management
- Goals tracking, hooks support, default durations
- Config stored in `~/.config/pomodoro/config.yml`

### Data Flow

1. CLI command parsing (Cobra) → Database session creation → TUI display (Bubble Tea) → Notification (beeep)
2. All sessions stored with start/end times, descriptions, tags, and break flags
3. Active sessions tracked by comparing current time with end times

### Key Dependencies

- **github.com/spf13/cobra** - CLI framework
- **github.com/charmbracelet/bubbletea** - TUI framework  
- **github.com/charmbracelet/bubbles** - TUI components (progress bars)
- **github.com/mattn/go-sqlite3** - SQLite driver
- **github.com/gen2brain/beeep** - Desktop notifications
- **gopkg.in/yaml.v3** - YAML configuration

### Testing Approach

- Tests located alongside source files (e.g., `cmd/break_test.go`)
- Database operations use interface for mocking
- TUI components tested through model updates
- Use `go test ./...` to run all tests

### Configuration and Data Storage

- Database: `~/.local/share/pomodoro/history.db`
- Config: `~/.config/pomodoro/config.yml` 
- Supports goals tracking, hooks, and custom durations
- Database uses WAL mode for concurrent access

## Development Notes

- Entry point is `main.go` which calls `cmd.Execute()`
- All commands support JSON output for scripting (`--json` flag)
- TUI uses different progress bar colors for Pomodoros (default gradient) vs breaks (green)
- Notifications are cross-platform but optimized for macOS
- Database schema includes indexing on date for performance