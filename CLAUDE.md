# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `hab` - a terminal-based habit tracker with GitHub-style contribution grids, built as a fast, compiled TUI application using Go and the Charm.sh library suite (Bubble Tea + Lipgloss).

## Development Commands

- `go run main.go` or `make run` - Run the hab TUI application
- `make build` - Build the hab binary
- `make install` - Install to /usr/local/bin
- `make clean` - Clean build artifacts

## Architecture

- **Entry Point**: `main.go` - Main application with Bubble Tea TUI
- **Framework**: Bubble Tea for interactive terminal UI
- **Styling**: Lipgloss for colors and layout
- **Runtime**: Compiled Go binary, no external dependencies

## Key Technologies

- **Go 1.21+**: Modern compiled language for fast execution
- **Bubble Tea**: The fun, functional TUI framework
- **Lipgloss**: Style definitions for nice terminal layouts  
- **Standard Library**: JSON parsing, file I/O, time handling

## Application Features

- **GitHub-style contributions grid** with circle character visualization (○ ◐ ◑ ●)
- **Multiple habit tracking** with navigation (shows all by default, individual targeting)
- **Multi-frequency habit support** with completion percentage tracking
- **Custom color support** for different habit types
- **365-day habit visualization** going backwards from current date
- **JSON data format** for easy import/export

## Data Structure

Habits are stored in `data/activities.json` with this structure:
```json
{
  "activities": {
    "activityKey": {
      "name": "Activity Name",
      "color": "red|blue|green|magenta|cyan|yellow",
      "target_per_day": 2,
      "dates": ["2025-01-15", "2025-01-15", "2025-01-20"]
    }
  }
}
```

For multi-frequency habits (like brushing teeth twice daily), set `target_per_day` and repeat dates for multiple completions.

## Visual Elements

The application supports three rendering levels with automatic detection:

### Unicode Level (Modern terminals)
- Circle characters for completion levels: ○ ◐ ◑ ● (none to complete)

### ASCII-Extended Level (Most terminals)  
- Box drawing characters for completion levels: ░ ▒ ▓ █ (none to complete)

### ASCII Level (Basic terminals)
- Basic ASCII characters for completion levels: . - + # (none to complete)

**Environment Override**: Set `HAB_RENDERING=ascii|extended|unicode` to force a specific level

- Terminal colors: red, blue, green, magenta, cyan, yellow, gray
- Lipgloss styling: bold text, color foregrounds, spacing and margins

## Project Structure

- `main.go` - Main application with Bubble Tea TUI and data models
- `data/activities.json` - Habit data storage
- `Makefile` - Build and install targets
- `go.mod` & `go.sum` - Go module dependencies