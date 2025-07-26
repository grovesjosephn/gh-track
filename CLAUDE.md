# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a terminal-based React application built with Ink, TypeScript, and Bun. It creates interactive CLI applications using React components that render to the terminal instead of the browser.

## Development Commands

- `bun dev` or `bun start` - Run the terminal application
- `bun run build` - Build TypeScript to JavaScript
- `bun run lint` - Run ESLint to check code quality

## Architecture

- **Entry Point**: `src/main.tsx` - Renders the React app to terminal using Ink
- **Main Component**: `src/App.tsx` - Main terminal application component
- **Runtime**: Node.js with tsx for TypeScript execution
- **Package Manager**: Bun for fast package management and execution
- **TypeScript**: Configured with separate configs for app (`tsconfig.app.json`) and Node (`tsconfig.node.json`)

## Key Technologies

- **Ink**: React renderer for interactive CLI applications
- **React 19.1.0**: Component-based architecture for terminal UI
- **TypeScript**: Type safety and development experience
- **tsx**: TypeScript execution for Node.js
- **Bun**: Fast package manager and JavaScript runtime

## Application Features

- **GitHub-style contributions grid** with ASCII block visualization
- **Multiple activity tracking** displayed in stacked grids
- **Custom color support** for different activity types
- **365-day activity visualization** with proper month alignment
- **JSON data format** for easy import/export

## Data Structure

Activities are stored in `data/activities.json` with this structure:
```json
{
  "activities": {
    "activityKey": {
      "name": "Activity Name",
      "color": "red|blue|green|magenta|cyan|yellow", 
      "dates": ["2025-01-15", "2025-01-20"]
    }
  }
}
```

## Ink Components and Patterns

- Use `<Box>` for layout and positioning (flexbox-like)
- Use `<Text>` for styled text output with colors and formatting
- ASCII blocks for activity levels: ▒ ░ ▓ █ (light to full)
- Terminal colors: red, blue, green, magenta, cyan, yellow, gray
- Text styling: bold, dimColor for visual hierarchy

## Project Structure

- `src/` - Source code
  - `App.tsx` - Main terminal application component
  - `main.tsx` - Application entry point with Ink renderer
- `dist/` - Build output for compiled JavaScript