# Session Templates

Session templates allow you to create and reuse predefined session configurations for common work types.

## Template Format

Templates are stored as YAML files in `~/.config/pomodoro/templates/` with the following structure:

```yaml
name: coding
description: Deep work coding session
duration: 50m
tags:
  - coding
  - focus
  - deep-work
audio:
  enabled: true
  volume: 0.7
created_at: 2025-08-22T13:15:00Z
updated_at: 2025-08-22T13:15:00Z
```

## Fields

- `name`: Template identifier (used as filename)
- `description`: Session description
- `duration`: Session length (e.g., "25m", "1h30m", "90s")
- `tags`: Array of session tags for categorization
- `audio`: Audio settings (optional)
  - `enabled`: Enable/disable audio notifications
  - `volume`: Audio volume (0.0 to 1.0)
- `created_at`/`updated_at`: Timestamps (automatically managed)

## CLI Commands

### Create Template
```bash
pomodoro template create coding \
  --description "Deep work coding session" \
  --duration 50m \
  --tags coding,focus,deep-work \
  --audio \
  --volume 0.7
```

### List Templates
```bash
pomodoro template list
```

### Show Template Details
```bash
pomodoro template show coding
```

### Update Template
```bash
pomodoro template update coding \
  --description "Updated description" \
  --duration 45m
```

### Start Session from Template
```bash
pomodoro template start coding
```

Override template values:
```bash
pomodoro template start coding \
  --duration 30m \
  --message "Custom session description"
```

### Delete Template
```bash
pomodoro template delete coding
```

### Export/Import Templates
```bash
# Export template
pomodoro template export coding ./my-coding-template.yml

# Import template
pomodoro template import ./my-coding-template.yml

# Import with overwrite
pomodoro template import ./my-coding-template.yml --overwrite
```

## Example Templates

This directory contains example templates you can import:

- `coding.yml` - Deep work coding sessions (50m)
- `meeting.yml` - Team meetings (25m, audio disabled)
- `research.yml` - Research and documentation (45m, quiet audio)