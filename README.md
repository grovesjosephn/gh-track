# hab - Terminal Habit Tracker

Track your daily habits with beautiful GitHub-style contribution grids, right in your terminal.

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white) ![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF69B4?style=flat) ![Terminal](https://img.shields.io/badge/terminal-ready-brightgreen)

## What is hab?

hab transforms habit tracking into a visual, engaging experience using familiar GitHub-style contribution grids. See your progress at a glance, track multiple habits simultaneously, and stay motivated with beautiful terminal visualizations.

### ✨ Key Features

- 🎯 **GitHub-style contribution grids** - Familiar visual progress tracking
- 🌈 **Color-coded habits** - Easy visual differentiation between different activities
- 📊 **Multi-frequency support** - Track daily, twice-daily, or custom frequency habits
- 📅 **Flexible timelines** - View 3, 6, or 12 months of history
- ⚡ **Lightning fast** - Compiled Go binary with zero dependencies
- 🎮 **Interactive TUI** - Modern, searchable interface with helpful shortcuts
- 📱 **Cross-platform** - Works on macOS, Linux, and Windows
- 📁 **Portable data** - Human-readable JSON format for easy backup/import

## Quick Start

### Installation

Install hab with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/grovesjosephn/hab/main/install.sh | bash
```

*Requirements: Git, Go 1.21+, and Make*

### Your First Habit

```bash
# Create a new habit
hab new exercise --color red

# Log an activity
hab exercise

# View your progress
hab
```

## Usage Guide

### Interactive Interface

Launch the beautiful terminal interface:

```bash
hab                    # Full 12-month view
hab -t 3m              # 3-month view
hab -t 6m              # 6-month view
hab --no-legend        # Hide the legend
```

**Navigation:**
- `Tab` - Open interactive habit selection
- `1-9` - Quick select habits by number
- `↑/↓` or `j/k` - Navigate between habits
- `Enter/Space` - Select habit or log today's activity
- `a` - Return to all habits view
- `ESC` - Go back
- `Ctrl+3/6/Y` - Switch timelines
- `L` - Toggle legend
- `?` - Show detailed help
- `q` or `Ctrl+C` - Quit

### Command Line Usage

**Creating Habits:**
```bash
hab new exercise                    # Basic habit
hab new exercise --color red        # With color
hab new meditation --target 2       # Twice-daily habit
```

**Tracking Activities:**
```bash
hab exercise                        # Log for today
hab add exercise 2025-01-15        # Log for specific date
hab exercise --date 2025-01-15     # Alternative syntax
```

**Managing Your Data:**
```bash
hab list                           # All habits with statistics
hab stats exercise                 # Detailed stats for one habit
hab prune                          # Clean up excess entries
hab prune --dry-run                # Preview cleanup
hab delete exercise                # Remove a habit
```

### Sample Output

**Interactive Grid:**
```
Activity Tracker - All Activities (12 months)

[1] Exercise (49 activities)
S  ●  ○  ●  ●  ○  ●  ◑  
M  ●  ●  ●  ○  ●  ●  ●  
T  ●  ●  ●  ●  ○  ●  ●  
W  ●  ●  ●  ●  ●  ●  ●  
T  ●  ●  ●  ●  ●  ○  ●  
F  ●  ●  ●  ●  ●  ●  ●  
S  ○  ●  ●  ●  ●  ●  ●  

[2] Reading (32 activities)
...

                             None  ○  ◐  ◑  ●  Complete
```

**Command Line:**
```bash
$ hab new exercise --color red
✓ Created habit 'Exercise' with color red
Add an entry with: hab exercise

$ hab exercise
✓ Added entry for 'Exercise' today
Current streak: 1 days 🔥

$ hab list
Your Habits:
============

[1] Exercise (red)
    Key: exercise
    Total entries: 1
    Unique days: 1
    Target per day: 1
    Current streak: 1 days 🔥
    Add entry: hab exercise
```

## Understanding the Visualization

### Completion Levels

hab automatically detects your terminal's capabilities and uses appropriate characters:

**Unicode (Modern terminals):**
- `○` No activity (0%)
- `◐` Low completion (< 50%)
- `◑` Partial completion (50-99%)
- `●` Target met/exceeded (100%+)

**ASCII-Extended (Most terminals):**
- `░` `▒` `▓` `█` (none to complete)

**Basic ASCII (All terminals):**
- `.` `-` `+` `#` (none to complete)

### Colors
- `red`, `blue`, `green`, `magenta`, `cyan`, `yellow`
- Each habit gets its own color for easy identification

### Multi-Frequency Habits

Track habits that occur multiple times per day:

```bash
hab new brushing --target 2        # Twice daily
hab brushing                       # Log first time
hab brushing                       # Log second time
```

The grid shows completion percentage based on your target.

## Data & Customization

### Data Location

Your habits are stored in standard locations:

- **macOS**: `~/Library/Application Support/hab/data/activities.json`
- **Linux**: `~/.config/hab/data/activities.json`
- **Windows**: `%APPDATA%/hab/data/activities.json`

**Custom location:**
```bash
export HAB_DATA_FILE="/path/to/my/habits.json"
```

### Data Format

Habits are stored in human-readable JSON:

```json
{
  "activities": {
    "exercise": {
      "name": "Exercise",
      "color": "red",
      "target_per_day": 1,
      "dates": ["2025-01-15", "2025-01-20", "2025-01-25"]
    },
    "meditation": {
      "name": "Meditation",
      "color": "blue",
      "target_per_day": 2,
      "dates": ["2025-01-15", "2025-01-15", "2025-01-16"]
    }
  }
}
```

### Terminal Customization

Force specific rendering modes:
```bash
HAB_RENDERING=unicode hab          # Force Unicode
HAB_RENDERING=extended hab         # Force ASCII-Extended  
HAB_RENDERING=ascii hab            # Force basic ASCII
```

Debug mode to see detected capabilities:
```bash
HAB_DEBUG=true hab
```

## Common Workflows

### Daily Check-in
```bash
hab                                # See today's status
hab exercise                       # Log workout
hab reading                        # Log reading session
```

### Weekly Review
```bash
hab stats exercise                 # See detailed progress
hab -t 3m                          # Check 3-month trends
```

### Data Management
```bash
hab prune --dry-run                # Check for cleanup opportunities
hab list                           # Review all habits
```

### Automation
```bash
# Add to your .bashrc/.zshrc
alias morning="hab reading && hab exercise && hab"

# Script integration
#!/bin/bash
echo "=== Daily Habit Check ==="
hab
```

---

## Installation Details

### Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/grovesjosephn/hab/main/install.sh | bash
```

Or with wget:
```bash
wget -qO- https://raw.githubusercontent.com/grovesjosephn/hab/main/install.sh | bash
```

**Custom install directory:**
```bash
INSTALL_DIR="$HOME/.local/bin" curl -fsSL https://raw.githubusercontent.com/grovesjosephn/hab/main/install.sh | bash
```

### Manual Installation

```bash
git clone https://github.com/grovesjosephn/hab.git
cd hab
make build
make install
```

## Development

### Project Structure
```
main.go              # Application entry point
cmd/                 # CLI commands (Cobra)
├── root.go          # Root command and TUI launcher
├── new.go           # Create new habits
├── add.go           # Add habit entries
├── list.go          # List all habits
├── stats.go         # Habit statistics
├── delete.go        # Delete habits
└── prune.go         # Clean up excess entries
internal/            # Data management
└── habit.go         # CRUD operations and data path logic
ui/                  # Terminal UI
└── tui.go           # Bubble Tea interface
Makefile            # Build and install targets
go.mod & go.sum     # Go module dependencies
```

### Tech Stack
- **Language**: Go 1.21+
- **CLI Framework**: Cobra for command parsing
- **TUI Framework**: Bubble Tea (Charm.sh)
- **UI Components**: Bubbles (Charm.sh)
- **Styling**: Lipgloss (Charm.sh)
- **Data Format**: JSON with standard library parsing

### Development Commands
```bash
make build          # Build binary
make run            # Run from source
make test           # Run tests
make clean          # Clean build artifacts
make fmt            # Format code
make vet            # Run go vet
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Run `make fmt && make vet && make test`
5. Commit: `git commit -m "Add feature"`
6. Push: `git push origin feature-name`
7. Open a Pull Request

---

## License

MIT License - feel free to use this project however you'd like!

## Inspiration

Inspired by GitHub's contribution graph, hab brings that familiar visualization to the terminal for tracking daily habits, with multi-frequency support and blazing fast Go performance.