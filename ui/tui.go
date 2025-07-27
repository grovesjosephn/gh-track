package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"hab/internal"
)

// RenderingLevel represents the terminal's font rendering capability
type RenderingLevel int

const (
	ASCII RenderingLevel = iota
	ASCIIExtended
	Unicode
)

// CharacterSet defines the characters to use for each completion level
type CharacterSet struct {
	None     string // 0% complete
	Low      string // < 50% complete  
	Partial  string // 50-99% complete
	Complete string // 100%+ complete
}

// Character sets for different rendering levels
var characterSets = map[RenderingLevel]CharacterSet{
	ASCII: {
		None:     ".",
		Low:      "-", 
		Partial:  "+",
		Complete: "#",
	},
	ASCIIExtended: {
		None:     "░",
		Low:      "▒",
		Partial:  "▓", 
		Complete: "█",
	},
	Unicode: {
		None:     "○",
		Low:      "◐",
		Partial:  "◑",
		Complete: "●",
	},
}

// ViewMode represents what to display
type ViewMode int

const (
	AllActivities ViewMode = iota
	SingleActivity
)

// ContributionGrid represents a grid cell for a specific date
type ContributionGrid struct {
	Date   time.Time
	Level  int
	Color  string
	Active bool
}

// Model represents the Bubble Tea model
type Model struct {
	habitManager   *internal.HabitManager
	activities     map[string]internal.Activity
	grid           [][]ContributionGrid
	width          int
	height         int
	ready          bool
	renderingLevel RenderingLevel
	viewMode       ViewMode
	activityKeys   []string
	selectedIndex  int
	timeline       TimelineDays
	showLegend     bool
}

// NewModel creates a new TUI model with default timeline
func NewModel() *Model {
	return NewModelWithOptions(Timeline12Months, true)
}

// NewModelWithTimeline creates a new TUI model with specified timeline
func NewModelWithTimeline(timeline TimelineDays) *Model {
	return NewModelWithOptions(timeline, true)
}

// NewModelWithOptions creates a new TUI model with specified timeline and legend visibility
func NewModelWithOptions(timeline TimelineDays, showLegend bool) *Model {
	hm := internal.NewHabitManager()
	if err := hm.Load(); err != nil {
		fmt.Printf("Error loading habits: %v\n", err)
		os.Exit(1)
	}

	activities := hm.GetActivities()
	grid := generateGrid(activities, timeline)
	renderingLevel := detectRenderingLevel()
	
	// Create sorted activity keys
	activityKeys := make([]string, 0, len(activities))
	for key := range activities {
		activityKeys = append(activityKeys, key)
	}
	sort.Strings(activityKeys)
	
	return &Model{
		habitManager:   hm,
		activities:     activities,
		grid:           grid,
		ready:          true,
		renderingLevel: renderingLevel,
		viewMode:       AllActivities,
		activityKeys:   activityKeys,
		selectedIndex:  0,
		timeline:       timeline,
		showLegend:     showLegend,
	}
}

// Detect terminal rendering capabilities
func detectRenderingLevel() RenderingLevel {
	// Allow manual override via environment variable
	if override := os.Getenv("HAB_RENDERING"); override != "" {
		switch strings.ToLower(override) {
		case "ascii":
			return ASCII
		case "extended", "ascii-extended":
			return ASCIIExtended
		case "unicode":
			return Unicode
		}
	}
	
	// Check environment variables for Unicode support
	term := strings.ToLower(os.Getenv("TERM"))
	lang := os.Getenv("LANG")
	lcAll := os.Getenv("LC_ALL")
	
	// Check for UTF-8 support in locale
	isUTF8 := strings.Contains(lang, "UTF-8") || strings.Contains(lcAll, "UTF-8")
	
	// Modern terminals that support Unicode
	modernTerms := []string{"xterm-256color", "screen-256color", "tmux-256color", "alacritty", "kitty", "iterm", "gnome-terminal"}
	for _, modernTerm := range modernTerms {
		if strings.Contains(term, modernTerm) && isUTF8 {
			return Unicode
		}
	}
	
	// ASCII-Extended support (most terminals support box drawing)
	extendedTerms := []string{"xterm", "screen", "tmux", "ansi", "vt100", "vt102", "vt220"}
	for _, extTerm := range extendedTerms {
		if strings.Contains(term, extTerm) {
			return ASCIIExtended
		}
	}
	
	// Fall back to basic ASCII for unknown or limited terminals
	return ASCII
}

// TimelineDays represents different timeline options
type TimelineDays int

const (
	Timeline3Months  TimelineDays = 90
	Timeline6Months  TimelineDays = 180
	Timeline12Months TimelineDays = 365
)

// Generate a grid going backwards from current date with specified timeline
func generateGrid(activities map[string]internal.Activity, timeline TimelineDays) [][]ContributionGrid {
	now := time.Now()
	
	// Start from today and go back specified number of days
	endDate := now
	startDate := endDate.AddDate(0, 0, -int(timeline-1)) // timeline days total (including today)
	
	// Create a map for quick date lookups
	activityDates := make(map[string]map[string]bool)
	for key, activity := range activities {
		activityDates[key] = make(map[string]bool)
		for _, dateStr := range activity.Dates {
			activityDates[key][dateStr] = true
		}
	}

	var weeks [][]ContributionGrid

	// Find the Sunday that contains or comes before our start date
	current := startDate
	for current.Weekday() != time.Sunday {
		current = current.AddDate(0, 0, -1)
	}

	// Generate weeks until we cover all dates up to today
	for current.Before(endDate) || current.Equal(endDate) {
		currentWeek := make([]ContributionGrid, 7)
		
		for day := 0; day < 7; day++ {
			dateStr := current.Format("2006-01-02")
			cell := ContributionGrid{
				Date:   current,
				Level:  0,
				Color:  "gray",
				Active: false,
			}

			// Only show data for dates within our range and not in the future
			if !current.Before(startDate) && !current.After(endDate) {
				// Check if this date has activities
				for key, activity := range activities {
					if activityDates[key][dateStr] {
						cell.Level = 1
						cell.Color = activity.Color
						cell.Active = true
						break // Use first matching activity color
					}
				}
			}

			currentWeek[day] = cell
			current = current.AddDate(0, 0, 1)
		}
		
		weeks = append(weeks, currentWeek)
		
		// Safety check to prevent infinite loop
		if len(weeks) > 60 { // Max ~14 months of weeks
			break
		}
	}

	return weeks
}

// Bubble Tea methods
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "a":
			// Show all activities
			m.viewMode = AllActivities
		case "tab":
			// Toggle between view modes
			if m.viewMode == AllActivities {
				m.viewMode = SingleActivity
			} else {
				m.viewMode = AllActivities
			}
		case "up", "k":
			// Navigate up in single activity mode
			if m.viewMode == SingleActivity && len(m.activityKeys) > 0 {
				m.selectedIndex = (m.selectedIndex - 1 + len(m.activityKeys)) % len(m.activityKeys)
			}
		case "down", "j":
			// Navigate down in single activity mode
			if m.viewMode == SingleActivity && len(m.activityKeys) > 0 {
				m.selectedIndex = (m.selectedIndex + 1) % len(m.activityKeys)
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			// Select activity by number
			num := int(msg.String()[0] - '1') // Convert '1' to 0, '2' to 1, etc.
			if num >= 0 && num < len(m.activityKeys) {
				m.selectedIndex = num
				m.viewMode = SingleActivity
			}
		case "ctrl+3", "alt+3":
			// Switch to 3 month timeline
			m.timeline = Timeline3Months
			m.grid = generateGrid(m.activities, m.timeline)
		case "ctrl+6", "alt+6":
			// Switch to 6 month timeline
			m.timeline = Timeline6Months
			m.grid = generateGrid(m.activities, m.timeline)
		case "ctrl+y", "alt+y":
			// Switch to 12 month (year) timeline
			m.timeline = Timeline12Months
			m.grid = generateGrid(m.activities, m.timeline)
		case "l":
			// Toggle legend visibility
			m.showLegend = !m.showLegend
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	var s strings.Builder
	
	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("6")).
		Margin(1, 0)
	
	// Title changes based on view mode and timeline
	var timelineText string
	switch m.timeline {
	case Timeline3Months:
		timelineText = "3 months"
	case Timeline6Months:
		timelineText = "6 months"
	case Timeline12Months:
		timelineText = "12 months"
	}
	
	var titleText string
	if m.viewMode == AllActivities {
		titleText = fmt.Sprintf("Activity Tracker - All Activities (%s)", timelineText)
	} else {
		selectedActivity := m.activities[m.activityKeys[m.selectedIndex]]
		titleText = fmt.Sprintf("Activity Tracker - %s (%s)", selectedActivity.Name, timelineText)
	}
	s.WriteString(titleStyle.Render(titleText))
	
	// Show rendering mode in debug mode
	if os.Getenv("HAB_DEBUG") == "true" {
		var modeText string
		switch m.renderingLevel {
		case ASCII:
			modeText = "ASCII"
		case ASCIIExtended:
			modeText = "ASCII-Extended"
		case Unicode:
			modeText = "Unicode"
		}
		debugStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
		s.WriteString("\n")
		s.WriteString(debugStyle.Render(fmt.Sprintf("Rendering Mode: %s", modeText)))
	}
	
	s.WriteString("\n\n")

	// Render activities based on view mode
	if m.viewMode == AllActivities {
		// Show all activities
		for i, key := range m.activityKeys {
			activity := m.activities[key]
			s.WriteString(m.renderActivityGrid(activity, key, i+1))
			if i < len(m.activityKeys)-1 {
				s.WriteString("\n\n")
			}
		}
	} else {
		// Show single selected activity
		if len(m.activityKeys) > 0 {
			key := m.activityKeys[m.selectedIndex]
			activity := m.activities[key]
			s.WriteString(m.renderActivityGrid(activity, key, -1)) // -1 means no number
		}
	}

	// Legend (right-aligned to grid end) - only show if enabled
	if m.showLegend {
		s.WriteString("\n\n")
		legendStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
		charSet := characterSets[m.renderingLevel]
		legendText := fmt.Sprintf("None  %s  %s  %s  %s  Complete", 
			charSet.None, charSet.Low, charSet.Partial, charSet.Complete)
		
		// Calculate grid width: day labels (3 chars) + weeks * 3 chars per week (1 char + 2 spaces)
		// But subtract the trailing 2 spaces from the last week
		gridWidth := 3 + len(m.grid)*3 - 2
		
		// Create the actual legend text with character substitutions
		actualLegendText := fmt.Sprintf("None  %s  %s  %s  %s  Complete", 
			charSet.None, charSet.Low, charSet.Partial, charSet.Complete)
		legendWidth := len(actualLegendText)
		
		// Create padding to align legend to the right edge of the grid
		if gridWidth > legendWidth {
			padding := strings.Repeat(" ", gridWidth-legendWidth)
			s.WriteString(padding)
		}
		s.WriteString(legendStyle.Render(legendText))
	}

	// Footer with controls
	s.WriteString("\n\n")
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8"))
	
	var controls string
	if m.viewMode == AllActivities {
		controls = "Controls: [1-9] Select activity • [Tab] Single view • [Ctrl+3/6/Y] Timeline • [L] Legend • [q/ESC] Quit"
	} else {
		controls = "Controls: [↑/↓] or [j/k] Navigate • [a] All activities • [Ctrl+3/6/Y] Timeline • [L] Legend • [q/ESC] Quit"
	}
	s.WriteString(footerStyle.Render(controls))

	return s.String()
}

// Render a single activity grid
func (m Model) renderActivityGrid(activity internal.Activity, activityKey string, activityNumber int) string {
	var s strings.Builder
	
	// Activity title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(getColorCode(activity.Color)))
	
	totalDates := len(activity.Dates)
	var titleText string
	if activityNumber > 0 {
		titleText = fmt.Sprintf("[%d] %s (%d activities)", activityNumber, activity.Name, totalDates)
	} else {
		titleText = fmt.Sprintf("%s (%d activities)", activity.Name, totalDates)
	}
	s.WriteString(titleStyle.Render(titleText))
	s.WriteString("\n")

	// Render grid rows (7 days per row)  
	dayLabels := []string{"S", "M", "T", "W", "T", "F", "S"}
	for row := 0; row < 7; row++ {
		s.WriteString(fmt.Sprintf("%-3s", dayLabels[row]))
		
		for week := 0; week < len(m.grid); week++ {
			if week < len(m.grid) && row < len(m.grid[week]) {
				cell := m.grid[week][row]
				char := m.getCellChar(cell, activity, activityKey)
				color := m.getCellColor(cell, activity, activityKey)
				
				cellStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
				s.WriteString(cellStyle.Render(char))
				s.WriteString("  ") // Two spaces for better week separation
			}
		}
		s.WriteString("\n")
	}

	return s.String()
}

// Get character for cell based on activity level
func (m Model) getCellChar(cell ContributionGrid, activity internal.Activity, activityKey string) string {
	dateStr := cell.Date.Format("2006-01-02")
	
	// Count how many times this activity was completed on this date
	completions := 0
	for _, activityDate := range activity.Dates {
		if activityDate == dateStr {
			completions++
		}
	}
	
	// Default target is 1 if not specified
	target := activity.TargetPerDay
	if target == 0 {
		target = 1
	}
	
	// Calculate completion percentage
	completionRate := float64(completions) / float64(target)
	
	// Get the appropriate character set for this terminal
	charSet := characterSets[m.renderingLevel]
	
	// Return character based on completion rate
	switch {
	case completionRate == 0:
		return charSet.None // No activity
	case completionRate < 0.5:
		return charSet.Low // Low completion (< 50%)
	case completionRate < 1.0:
		return charSet.Partial // Partial completion (50-99%)
	default:
		return charSet.Complete // Target met or exceeded
	}
}

// Get color for cell based on activity
func (m Model) getCellColor(cell ContributionGrid, activity internal.Activity, activityKey string) string {
	dateStr := cell.Date.Format("2006-01-02")
	for _, activityDate := range activity.Dates {
		if activityDate == dateStr {
			return getColorCode(activity.Color)
		}
	}
	return "8" // Dim gray for inactive
}

// Convert color names to terminal color codes
func getColorCode(colorName string) string {
	colorMap := map[string]string{
		"red":     "1",
		"green":   "2", 
		"yellow":  "3",
		"blue":    "4",
		"magenta": "5",
		"cyan":    "6",
		"gray":    "8",
	}
	
	if code, exists := colorMap[colorName]; exists {
		return code
	}
	return "7" // Default to white
}

// RunTUI starts the interactive TUI
func RunTUI() {
	RunTUIWithTimeline(Timeline12Months)
}

func RunTUIWithTimeline(timeline TimelineDays) {
	RunTUIWithOptions(timeline, true)
}

func RunTUIWithOptions(timeline TimelineDays, showLegend bool) {
	m := NewModelWithOptions(timeline, showLegend)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}