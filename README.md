# 🍅 Pomodoro CLI

A minimalist macOS CLI Pomodoro timer built in Go with advanced features for productivity tracking.

## ✨ Features

### Core Timer Functionality
- **Start Pomodoro Sessions** - Focus work sessions with customizable durations and descriptions
- **Break Timers** - Short and long break management
- **Pause/Resume** - Pause sessions temporarily and resume where you left off
- **Session Cancellation** - End sessions early when needed

### Audio & Notifications
- **Audio Alerts** 🔊 - Sound notifications when sessions complete (configurable)
- **Desktop Notifications** - Native OS notifications with session details
- **Silent Mode** - Disable audio for individual sessions with `--silent` flag

### Session Management
- **Continuous Mode** - Stay in the program after session completion for seamless workflow
- **Session History** - Track all your completed pomodoros and breaks
- **Tags & Organization** - Organize sessions with custom tags
- **Goal Tracking** - Set and monitor daily/weekly pomodoro targets

### Advanced Features
- **Progress Visualization** - Beautiful terminal UI with animated progress bars
- **Multiple Export Formats** - JSON and Open Pomodoro Format (OPF) support
- **Flexible Configuration** - YAML-based configuration with hooks support
- **Input Validation** - Smart validation and sanitization of user inputs

## 🚀 Quick Start

### Installation

#### Pre-built Binaries (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/ethan-k/pomodoro-cli/releases).

**macOS:**
```bash
# Intel Macs
curl -L -o pomodoro https://github.com/ethan-k/pomodoro-cli/releases/latest/download/pomodoro-darwin-amd64
# Apple Silicon Macs  
curl -L -o pomodoro https://github.com/ethan-k/pomodoro-cli/releases/latest/download/pomodoro-darwin-arm64

chmod +x pomodoro
sudo mv pomodoro /usr/local/bin/
```

**Linux:**
```bash
# x86_64
curl -L -o pomodoro https://github.com/ethan-k/pomodoro-cli/releases/latest/download/pomodoro-linux-amd64
# ARM64
curl -L -o pomodoro https://github.com/ethan-k/pomodoro-cli/releases/latest/download/pomodoro-linux-arm64

chmod +x pomodoro
sudo mv pomodoro /usr/local/bin/
```

**Windows:**
Download the appropriate `.exe` file from the releases page and add it to your PATH.

#### Build from Source

```bash
git clone https://github.com/ethan-k/pomodoro-cli
cd pomodoro-cli
make build
make install
```

#### Verify Installation

```bash
pomodoro --version
```

### Basic Usage

```bash
# Start a 25-minute pomodoro
pomodoro start "Fix authentication bug"

# Start with custom duration and tags
pomodoro start "Code review" --duration 50m --tags coding,review

# Start with continuous mode (stay in program after completion)
pomodoro start "Deep work" --continuous

# Start in silent mode (no audio alerts)
pomodoro start "Meeting focus" --silent

# Take a break
pomodoro break 5m --wait

# Pause current session
pomodoro pause

# Resume paused session
pomodoro resume --wait

# Check status
pomodoro status

# View session history
pomodoro history --today
pomodoro history --week
pomodoro history --output json
```

## 📖 Commands Reference

### Session Commands

| Command | Description | Examples |
|---------|-------------|----------|
| `start` | Start a pomodoro session | `pomodoro start "Task name"` |
| `break` | Start a break timer | `pomodoro break 10m` |
| `pause` | Pause active session | `pomodoro pause` |
| `resume` | Resume paused session | `pomodoro resume --wait` |
| `cancel` | Cancel active session | `pomodoro cancel` |
| `status` | Show current session status | `pomodoro status` |

### Data & Analysis

| Command | Description | Examples |
|---------|-------------|----------|
| `history` | View session history | `pomodoro history --today` |
| `config` | Manage configuration | `pomodoro config show` |

### Global Flags

| Flag | Description | Available Commands |
|------|-------------|-------------------|
| `--json` | JSON output format | All commands |
| `--silent` | Disable audio alerts | `start`, `break` |
| `--continuous` | Continuous mode | `start` |
| `--wait` | Show progress bar | `break`, `resume`, `status` |

## ⚙️ Configuration

Configuration is stored in `~/.config/pomodoro/config.yml`:

```yaml
# Goal settings
goals:
  daily_count: 8      # Target pomodoros per day
  weekly_count: 40    # Target pomodoros per week

# Default durations
defaults:
  pomodoro_duration: "25m"
  break_duration: "5m"
  long_break_duration: "15m"

# Audio settings
audio:
  enabled: true
  volume: 0.7                    # 0.0 to 1.0
  custom_sounds_dir: "~/.config/pomodoro/sounds"
  sounds:
    pomodoro_complete: "pomodoro_complete.wav"
    break_complete: "break_complete.wav"
    session_start: "session_start.wav"

# Data storage paths
paths:
  database: "~/.local/share/pomodoro/history.db"
  opf_export: "~/.local/share/pomodoro/exports"

# Hooks for automation
hooks:
  enabled: false
  path: "~/.config/pomodoro/hooks"
```

### Audio Configuration

#### Built-in Sounds

The Pomodoro CLI includes high-quality, copyright-free notification sounds:

- **Pomodoro Complete**: Gentle bell chime with natural decay
- **Break Complete**: Soft piano chord (C major) 
- **Session Start**: Light two-tone notification

All sounds are:
- ✅ **Copyright-free** and safe for commercial use
- 🎵 **Musically pleasant** - designed to be non-jarring
- ⚡ **Optimized** for notification purposes (2-3 seconds)
- 🔊 **Cross-platform** compatible (WAV format)

#### Custom Sounds
Replace with your own sounds by placing files in `~/.config/pomodoro/sounds/`:

```yaml
audio:
  sounds:
    pomodoro_complete: "my-custom-bell.wav"
    break_complete: "my-custom-chime.mp3"
```

#### Volume Control
```bash
# Set volume in config file
pomodoro config audio --volume 0.8
```

## 📊 Session History & Analytics

### View History
```bash
# Today's sessions
pomodoro history --today

# This week's sessions
pomodoro history --week

# Custom date range
pomodoro history --from 2025-01-01 --to 2025-01-31

# Filter by tags
pomodoro history --tags coding,review

# Export formats
pomodoro history --output json > sessions.json
pomodoro history --output opf > sessions-opf.json
```

### Session Data Structure

Each session includes:
- **ID** - Unique session identifier
- **Start/End Times** - Precise timing data
- **Description** - Session description
- **Duration** - Planned vs actual duration
- **Tags** - Organization labels
- **Type** - Pomodoro or break
- **Pause Data** - Pause/resume tracking

## 🔄 Workflow Examples

### Classic Pomodoro Technique
```bash
# Start 25min work session
pomodoro start "Feature development" --continuous

# After completion, program prompts for next action:
# 1. Start a break (b)
# 2. Start another pomodoro (p) 
# 3. View status (s)
# 4. Quit (q)
```

### Extended Work Sessions
```bash
# 50-minute deep work session
pomodoro start "Complex analysis" --duration 50m --tags analysis,research

# Custom break
pomodoro break 15m --wait
```

### Meeting Focus
```bash
# Silent mode for meetings
pomodoro start "Team meeting" --duration 1h --silent --tags meeting
```

## 🛠️ Development

### Build Commands
```bash
make build          # Build for current platform
make build-all      # Build for all platforms
make test           # Run tests
make coverage       # Generate coverage report
make fmt            # Format code
make lint           # Run linter
make install        # Install to $GOPATH/bin
```

### Project Structure
```
pomodoro-cli/
├── cmd/                    # CLI commands
│   ├── start.go           # Pomodoro start command
│   ├── break.go           # Break timer command
│   ├── pause.go           # Pause session command
│   ├── resume.go          # Resume session command
│   ├── status.go          # Session status command
│   └── history.go         # History management
├── internal/
│   ├── audio/             # Audio notification system
│   ├── config/            # Configuration management
│   ├── db/                # SQLite database layer
│   ├── model/             # Bubble Tea UI models
│   ├── notify/            # Notification system
│   └── utils/             # Shared utilities
├── specs/                 # Feature specifications
└── Makefile              # Build automation
```

## 📈 Advanced Usage

### Automation & Scripting

The JSON output format makes scripting easy:

```bash
# Check if session is active
if pomodoro status --json | jq -r '.active' == "true"; then
    echo "Session in progress"
fi

# Get today's pomodoro count
count=$(pomodoro history --today --output json | jq 'map(select(.was_break == false)) | length')
echo "Completed $count pomodoros today"
```

### Integration Examples

#### Git Hooks
```bash
# .git/hooks/pre-commit
#!/bin/bash
if pomodoro status --json | jq -r '.active' == "true"; then
    echo "⚠️  Pomodoro session active - consider finishing before committing"
fi
```

#### Shell Prompt
```bash
# Add to .bashrc/.zshrc
function pomodoro_prompt() {
    status=$(pomodoro status --json 2>/dev/null)
    if echo "$status" | jq -r '.active' 2>/dev/null | grep -q "true"; then
        if echo "$status" | jq -r '.status' 2>/dev/null | grep -q "paused"; then
            echo "⏸️"
        else
            echo "🍅"
        fi
    fi
}

# Use in prompt
PS1="$(pomodoro_prompt) $PS1"
```

## 🤝 Contributing

We welcome contributions! Please see our [feature specifications](specs/) for planned enhancements.

### Development Setup

```bash
# Clone and setup
git clone https://github.com/ethan-k/pomodoro-cli
cd pomodoro-cli
go mod download

# Run tests
make test

# Build and test locally
make build
./bin/pomodoro --version
```

### CI/CD Pipeline

The project uses GitHub Actions for:
- **Continuous Integration**: Tests on Linux, macOS, and Windows
- **Automated Releases**: Multi-platform binary builds on git tags
- **Security Scanning**: Automated security analysis with Gosec
- **Code Quality**: Linting with golangci-lint

### Creating a Release

For maintainers:

```bash
# Create a new release (replace 1.0.0 with actual version)
./scripts/deploy.sh 1.0.0
```

This will:
1. Run tests and builds
2. Create and push a git tag
3. Trigger GitHub Actions to build and publish release binaries

### High Priority Features
- Session templates and presets
- Enhanced goal tracking UI
- Comprehensive test coverage
- Additional export formats

## 📄 License

MIT License - see LICENSE file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI framework
- UI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Cross-platform notifications via [beeep](https://github.com/gen2brain/beeep)
- Data storage with SQLite

---

*Stay focused, stay productive! 🍅*