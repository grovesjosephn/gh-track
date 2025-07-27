# hab - Habit Tracker

A fast, terminal-based habit tracker with GitHub-style contribution grids, built with Go and the Charm.sh TUI library suite.

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white) ![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF69B4?style=flat) ![Terminal](https://img.shields.io/badge/terminal-ready-brightgreen)

## Features

- 🎯 **GitHub-style contribution grids** with circle character visualization (○ ◐ ◑ ●)
- 🌈 **Color-coded habits** for easy differentiation  
- 📊 **Multi-frequency habit support** with completion percentage tracking
- 📅 **365-day habit visualization** going backwards from today
- 📁 **Human-readable JSON format** for easy import/export
- ⚡ **Fast compiled binary** with no runtime dependencies
- 🔄 **Cross-platform terminal support** with automatic character fallbacks
- 🎮 **Interactive navigation** between all habits and individual views

## Installation

```bash
# Build the application
make build

# Or install to /usr/local/bin
make install

# Or run directly
make run
```

## Usage

### Interactive TUI (Default)
```bash
# Launch interactive visualization
hab

# Or explicitly
hab -i
```

**TUI Controls:**
- `1-9` - Select specific habit by number  
- `Tab` - Switch between all/single habit view
- `↑/↓` or `j/k` - Navigate between habits (single view)
- `a` - Return to all habits view
- `q`, `ESC`, or `Ctrl+C` - Quit

### Command Line Interface

**Create habits:**
```bash
hab new exercise                    # Create new habit
hab new exercise --color red        # With specific color
hab new brushing --target 2         # Multi-frequency habit
```

**Track habits:**
```bash
hab exercise                        # Add entry for today
hab add exercise 2025-01-15        # Add entry for specific date
hab exercise --date 2025-01-15     # Alternative syntax
```

**View and manage:**
```bash
hab list                           # List all habits with stats
hab stats exercise                 # Detailed statistics
hab delete exercise                # Delete habit (with confirmation)
hab delete exercise --force        # Delete without confirmation
```

### Sample Output

**Interactive TUI:**
```
Activity Tracker - All Activities

[1] Brushing Teeth (139 activities)
S  ●  ●  ◑  ●  ●  ●  ◑  
M  ●  ●  ●  ●  ●  ●  ●  
T  ●  ●  ●  ●  ●  ●  ●  
W  ●  ●  ●  ●  ●  ●  ●  
T  ●  ●  ●  ●  ●  ●  ●  
F  ●  ●  ●  ●  ●  ●  ●  
S  ●  ●  ●  ●  ●  ●  ●  

None  ○  ◐  ◑  ●  Complete

[2] Exercise (49 activities)
...

Controls: [1-9] Select habit • [Tab] Single view • [q/ESC] Quit
```

**CLI Commands:**
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

## Data Location

By default, habits are stored in platform-appropriate locations:

- **macOS**: `~/Library/Application Support/hab/data/activities.json`
- **Linux**: `~/.config/hab/data/activities.json`
- **Windows**: `%APPDATA%/hab/data/activities.json`
- **Fallback**: `./data/activities.json` (current directory)

**Custom Location**: Set `HAB_DATA_FILE` environment variable to override:
```bash
export HAB_DATA_FILE="/path/to/my/habits.json"
hab
```

## Data Structure

Habits are stored in JSON format:

```json
{
  "activities": {
    "exercise": {
      "name": "Exercise",
      "color": "red",
      "dates": ["2025-01-15", "2025-01-20", "2025-01-25"]
    },
    "brushing_teeth": {
      "name": "Brushing Teeth",
      "color": "blue",
      "target_per_day": 2,
      "dates": [
        "2025-01-15", "2025-01-15",
        "2025-01-16", 
        "2025-01-17", "2025-01-17"
      ]
    }
  }
}
```

**Multi-completion Habits**: For habits that should be done multiple times per day (like brushing teeth twice daily), set `target_per_day` and repeat the date in the array for each completion.

### Supported Colors
- `red`, `blue`, `green`, `magenta`, `cyan`, `yellow`
- Color field is optional (defaults to green)

### Adding Your Own Data
1. Create habits using `hab new [habit-name]` command
2. Add entries using `hab [habit-name]` or `hab add [habit-name] [date]`
3. Or manually edit the activities.json file in your data directory
4. Run `hab` to see your updated grid

## Character Visualization

The application automatically detects your terminal's rendering capabilities and uses appropriate characters:

### Unicode Mode (Modern terminals)
- `○` - No activity (0% complete)
- `◐` - Low completion (< 50% of target)
- `◑` - Partial completion (50-99% of target) 
- `●` - Target met or exceeded (100%+)

### ASCII-Extended Mode (Most terminals)  
- `░` - No activity (0% complete)
- `▒` - Low completion (< 50% of target)
- `▓` - Partial completion (50-99% of target)
- `█` - Target met or exceeded (100%+)

### ASCII Mode (Basic terminals)
- `.` - No activity (0% complete)
- `-` - Low completion (< 50% of target)
- `+` - Partial completion (50-99% of target)
- `#` - Target met or exceeded (100%+)

Each activity type uses its specified color (red, blue, green, magenta, etc.)

### Manual Override
You can force a specific rendering mode:
```bash
# Force Unicode characters
HAB_RENDERING=unicode hab

# Force ASCII-Extended characters  
HAB_RENDERING=extended hab

# Force basic ASCII characters
HAB_RENDERING=ascii hab
```

### Debug Mode
See which rendering mode was auto-detected:
```bash
HAB_DEBUG=true hab
```

## Use Cases

### Quick Status Check
```bash
hab
```

### Save to File
```bash
hab > habit-report.txt
```

### Integrate into Scripts
```bash
#!/bin/bash
echo "=== Daily Habit Report ==="
cd ~/hab && hab
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
└── delete.go        # Delete habits
internal/            # Data management
└── habit.go         # CRUD operations and data path logic
Makefile            # Build and install targets
go.mod & go.sum     # Go module dependencies
```

### Commands
```bash
# Development
go run main.go

# Build for production  
make build

# Test all functionality
make test

# Clean build artifacts
make clean
```

### Tech Stack
- **Language**: Go 1.21+
- **CLI Framework**: Cobra for command parsing
- **TUI Framework**: Bubble Tea (Charm.sh)
- **Styling**: Lipgloss (Charm.sh)
- **Data Format**: JSON with standard library parsing

## Sample Data

The application starts with no habits. Create your first habit with:
```bash
hab new exercise --color red
```

Then start tracking:
```bash
hab exercise  # Add entry for today
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Commit: `git commit -m "Add feature"`
5. Push: `git push origin feature-name`
6. Open a Pull Request

## License

MIT License - feel free to use this project however you'd like!

## Inspiration

Inspired by GitHub's contribution graph, this project brings that familiar visualization to the terminal for tracking daily habits, with multi-frequency support and blazing fast Go performance.