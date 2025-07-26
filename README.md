# GitHub Activity Tracker

A terminal-based activity tracker that visualizes your daily activities using GitHub-style contribution grids with ASCII blocks and colors.

![Terminal Activity Tracker Demo](https://img.shields.io/badge/terminal-ready-brightgreen) ![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=flat&logo=typescript&logoColor=white) ![React](https://img.shields.io/badge/React-20232A?style=flat&logo=react&logoColor=61DAFB) ![Bun](https://img.shields.io/badge/Bun-282a36?style=flat&logo=bun&logoColor=fbf0df)

## Features

- 🎯 **GitHub-style contribution grids** with ASCII block visualization
- 🌈 **Color-coded activities** for easy differentiation
- 📊 **365-day activity tracking** with day-of-week labels
- 📁 **Human-readable JSON format** for easy import/export
- ⚡ **Fast terminal rendering** using React + Ink
- 🔄 **Auto-print and exit** - perfect for scripts and quick status checks

## Installation

### Prerequisites
- [Bun](https://bun.sh/) or Node.js 18+

### Clone and Setup
```bash
git clone https://github.com/grovesjosephn/gh-track.git
cd gh-track
bun install
```

## Usage

### Quick Start
```bash
# Run the activity tracker
bun dev
# or
bun start
```

### Sample Output
```
Activity Tracker - All Activities

Exercise (49 activities)
S  ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒ ▒░▒ ░▒ ░░▒ ▒▒ ░
M  ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒
T  ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒
W  ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒
T  ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒
F  ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒
S  ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒ ▒▒▒ ▒▒ ▒▒▒ ▒▒
   Less ▒░▓██ More

Reading (59 activities)
...
```

## Data Format

Activities are stored in `data/activities.json`:

```json
{
  "activities": {
    "exercise": {
      "name": "Exercise",
      "color": "red",
      "dates": ["2025-01-15", "2025-01-16", "2025-01-20"]
    },
    "reading": {
      "name": "Reading", 
      "color": "blue",
      "dates": ["2025-01-15", "2025-01-17", "2025-01-18"]
    }
  }
}
```

### Supported Colors
- `red`, `blue`, `green`, `magenta`, `cyan`, `yellow`
- Color field is optional (defaults to green)

### Adding Your Own Data
1. Edit `data/activities.json`
2. Add your activity with dates in `YYYY-MM-DD` format
3. Optionally specify a color
4. Run `bun dev` to see your updated grid

## ASCII Visualization

Activities are represented using ASCII blocks with different intensities:
- `▒` - No activity (gray)
- `░` - Low activity  
- `▓` - Medium activity
- `█` - High activity

Each activity type uses its specified color (red, blue, green, magenta, etc.)

## Use Cases

### Quick Status Check
```bash
bun dev
```

### Save to File
```bash
bun dev > activity-report.txt
```

### Filter Specific Activities
```bash
bun dev | grep "Exercise"
```

### Integrate into Scripts
```bash
#!/bin/bash
echo "=== Daily Activity Report ==="
cd ~/gh-track && bun dev
```

## Development

### Project Structure
```
src/
├── components/
│   ├── ContributionsGrid.tsx    # Main grid visualization
│   └── ActivitySelector.tsx     # Activity selection UI
├── utils/
│   ├── dateUtils.ts            # Date manipulation helpers
│   └── dataLoader.ts           # JSON data processing
├── App.tsx                     # Main application
└── main.tsx                    # Entry point
```

### Commands
```bash
# Development
bun dev

# Build
bun run build

# Lint
bun run lint
```

### Tech Stack
- **Runtime**: Bun
- **UI Framework**: React 19 + Ink (terminal renderer)
- **Language**: TypeScript
- **Linting**: ESLint

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

Inspired by GitHub's contribution graph, this project brings that familiar visualization to the terminal for tracking any type of daily activities.
