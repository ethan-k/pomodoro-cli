#!/bin/bash

# Deploy script for pomodoro CLI
# This script creates a new release by tagging and pushing to GitHub

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if version is provided
if [ $# -eq 0 ]; then
    print_error "Usage: $0 <version>"
    print_error "Example: $0 1.0.0"
    exit 1
fi

VERSION="$1"
TAG="v${VERSION}"

# Validate version format (basic semver check)
if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    print_error "Invalid version format. Please use semantic versioning (e.g., 1.0.0)"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    print_error "Not in a git repository"
    exit 1
fi

# Check if working directory is clean
if ! git diff-index --quiet HEAD --; then
    print_error "Working directory is not clean. Please commit or stash your changes."
    exit 1
fi

# Check if tag already exists
if git tag -l | grep -q "^${TAG}$"; then
    print_error "Tag ${TAG} already exists"
    exit 1
fi

# Fetch latest changes
print_status "Fetching latest changes..."
git fetch origin

# Check if we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    print_warning "You're not on the main branch (current: $CURRENT_BRANCH)"
    read -p "Do you want to continue? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_status "Deployment cancelled"
        exit 1
    fi
fi

# Run tests
print_status "Running tests..."
if ! make test; then
    print_error "Tests failed. Please fix them before releasing."
    exit 1
fi

# Build project
print_status "Building project..."
if ! make build; then
    print_error "Build failed. Please fix build issues before releasing."
    exit 1
fi

# Create and push tag
print_status "Creating tag ${TAG}..."
git tag -a "${TAG}" -m "Release ${TAG}"

print_status "Pushing tag to origin..."
git push origin "${TAG}"

print_status "✅ Release ${TAG} has been created!"
print_status "GitHub Actions will now build and publish the release."
print_status "Monitor progress at: https://github.com/$(git config remote.origin.url | sed 's/.*[:/]\([^/]*\/[^/]*\)\.git$/\1/')/actions"

# Optional: Open GitHub releases page
if command -v open >/dev/null 2>&1; then
    read -p "Open GitHub releases page? [y/N] " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        REPO_URL=$(git config remote.origin.url | sed 's/git@github.com:/https:\/\/github.com\//' | sed 's/\.git$//')
        open "${REPO_URL}/releases"
    fi
fi