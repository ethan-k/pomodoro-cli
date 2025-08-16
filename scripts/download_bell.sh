#!/bin/bash

# Download a simple, free bell notification sound for Pomodoro CLI
# This script downloads copyright-free bell sounds from free resources

SOUNDS_DIR="internal/audio/sounds"

echo "🔔 Downloading free bell notification sounds..."

# Create sounds directory if it doesn't exist
mkdir -p "$SOUNDS_DIR"

# Download a simple bell notification sound using curl
# Using a placeholder URL - you would replace this with actual Mixkit download links

# Method 1: Try to download from Mixkit (you need to get the actual download URLs)
echo "📥 Attempting to download bell notification from free sources..."

# For now, let's create a simple curl command template
# You would need to visit Mixkit, find the specific sound, and get the download URL

echo "⚠️  To download from Mixkit:"
echo "1. Visit: https://mixkit.co/free-sound-effects/bell/"
echo "2. Find 'Bell notification' or similar"
echo "3. Right-click download button and copy link"
echo "4. Replace URL below:"

echo ""
echo "Example download command:"
echo "curl -L 'MIXKIT_DOWNLOAD_URL' -o '$SOUNDS_DIR/pomodoro_complete.wav'"

echo ""
echo "Alternative: Use freesound.org with Creative Commons license"
echo "1. Visit: https://freesound.org/search/?q=bell+notification"
echo "2. Filter by: Creative Commons 0 (CC0) license"
echo "3. Download WAV format"

echo ""
echo "🎯 Recommended characteristics:"
echo "   - Duration: 1-3 seconds"
echo "   - Format: WAV or MP3"
echo "   - License: CC0 or Public Domain"
echo "   - Volume: Gentle/quiet"
echo "   - Type: Simple bell, chime, or ding"

echo ""
echo "💾 Save files as:"
echo "   - pomodoro_complete.wav (main notification)"
echo "   - break_complete.wav (break end)"
echo "   - session_start.wav (session start)"