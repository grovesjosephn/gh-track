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

```bash
# Run the application
./hab

# Or if installed globally
hab
```

**Controls:**

**All Activities View (Default):**
- `1-9` - Select specific activity by number
- `Tab` - Switch to single activity view
- `q`, `ESC`, or `Ctrl+C` - Quit

**Single Activity View:**
- `↑/↓` or `j/k` - Navigate between activities
- `a` - Return to all activities view
- `Tab` - Toggle back to all activities
- `q`, `ESC`, or `Ctrl+C` - Quit

### Sample Output
```
Activity Tracker - All Activities

[1] Brushing Teeth (128 activities)
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

## Data Structure

Habits are stored in `data/activities.json`:

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
1. Edit `data/activities.json`
2. Add your habit with dates in `YYYY-MM-DD` format
3. Optionally specify a color and `target_per_day`
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
main.go              # Main application with Bubble Tea TUI
data/activities.json # Habit data storage  
Makefile            # Build and install targets
go.mod & go.sum     # Go module dependencies
```

### Commands
```bash
# Development
go run main.go

# Build for production
go build -o hab main.go

# Clean build artifacts
make clean
```

### Tech Stack
- **Language**: Go 1.21+
- **TUI Framework**: Bubble Tea (Charm.sh)
- **Styling**: Lipgloss (Charm.sh)
- **Data Format**: JSON with standard library parsing

## Sample Data

The repository includes sample data for five habits:
- Exercise (red)
- Reading (blue) 
- Coding (green)
- Meditation (magenta)
- Brushing Teeth (cyan, 2x daily target)

You can modify `data/activities.json` to track your own habits.

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