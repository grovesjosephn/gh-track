package ui

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
	HabitSelection
)

// HabitItem represents an item in the habit list
type HabitItem struct {
	key      string
	activity internal.Activity
}

func (i HabitItem) FilterValue() string { return i.activity.Name }
func (i HabitItem) Title() string       { return i.activity.Name }
func (i HabitItem) Description() string {
	return fmt.Sprintf("Key: %s • Color: %s • Target: %d/day • Entries: %d", 
		i.key, i.activity.Color, max(1, i.activity.TargetPerDay), len(i.activity.Dates))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ContributionGrid represents a grid cell for a specific date
type ContributionGrid struct {
	Date   time.Time
	Level  int
	Color  string
	Active bool
}

// Key bindings
type keyMap struct {
	Up          key.Binding
	Down        key.Binding
	Left        key.Binding
	Right       key.Binding
	Enter       key.Binding
	Space       key.Binding
	Tab         key.Binding
	Quit        key.Binding
	Escape      key.Binding
	AllView     key.Binding
	Timeline3m  key.Binding
	Timeline6m  key.Binding
	Timeline12m key.Binding
	ToggleLegend key.Binding
	Help        key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Tab, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Space},
		{k.Tab, k.AllView, k.ToggleLegend},
		{k.Timeline3m, k.Timeline6m, k.Timeline12m},
		{k.Help, k.Quit, k.Escape},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select/log activity"),
	),
	Space: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "log activity"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch view"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	AllView: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "all activities"),
	),
	Timeline3m: key.NewBinding(
		key.WithKeys("ctrl+3"),
		key.WithHelp("ctrl+3", "3 months"),
	),
	Timeline6m: key.NewBinding(
		key.WithKeys("ctrl+6"),
		key.WithHelp("ctrl+6", "6 months"),
	),
	Timeline12m: key.NewBinding(
		key.WithKeys("ctrl+y"),
		key.WithHelp("ctrl+y", "12 months"),
	),
	ToggleLegend: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "toggle legend"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
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
	habitList      list.Model
	help           help.Model
	keys           keyMap
	showHelp       bool
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

	// Create list items for bubbles/list
	items := make([]list.Item, 0, len(activityKeys))
	for _, key := range activityKeys {
		items = append(items, HabitItem{
			key:      key,
			activity: activities[key],
		})
	}

	// Configure the list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a Habit"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))

	// Configure help
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	h.Styles.ShortDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	
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
		habitList:      l,
		help:           h,
		keys:           keys,
		showHelp:       false,
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle help toggle first
		if key.Matches(msg, m.keys.Help) {
			m.showHelp = !m.showHelp
			return m, nil
		}

		// Handle quit
		if key.Matches(msg, m.keys.Quit) {
			return m, tea.Quit
		}

		// Handle view-specific keys
		switch m.viewMode {
		case HabitSelection:
			// Update list
			m.habitList, cmd = m.habitList.Update(msg)
			
			// Handle selection
			if key.Matches(msg, m.keys.Enter) {
				if selectedItem, ok := m.habitList.SelectedItem().(HabitItem); ok {
					// Find the index of the selected habit
					for i, key := range m.activityKeys {
						if key == selectedItem.key {
							m.selectedIndex = i
							break
						}
					}
					m.viewMode = SingleActivity
				}
			}
			
			// Handle escape to go back
			if key.Matches(msg, m.keys.Escape) {
				m.viewMode = AllActivities
			}

		case AllActivities:
			// Handle tab to go to habit selection
			if key.Matches(msg, m.keys.Tab) {
				m.viewMode = HabitSelection
			}
			
			// Handle number keys for direct selection
			switch msg.String() {
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				num := int(msg.String()[0] - '1')
				if num >= 0 && num < len(m.activityKeys) {
					m.selectedIndex = num
					m.viewMode = SingleActivity
				}
			}

		case SingleActivity:
			// Handle navigation
			if key.Matches(msg, m.keys.Up) && len(m.activityKeys) > 0 {
				m.selectedIndex = (m.selectedIndex - 1 + len(m.activityKeys)) % len(m.activityKeys)
			}
			if key.Matches(msg, m.keys.Down) && len(m.activityKeys) > 0 {
				m.selectedIndex = (m.selectedIndex + 1) % len(m.activityKeys)
			}
			
			// Handle tab to go to habit selection  
			if key.Matches(msg, m.keys.Tab) {
				m.viewMode = HabitSelection
			}
			
			// Handle all activities view
			if key.Matches(msg, m.keys.AllView) {
				m.viewMode = AllActivities
			}
			
			// Handle logging activity
			if (key.Matches(msg, m.keys.Enter) || key.Matches(msg, m.keys.Space)) && len(m.activityKeys) > 0 {
				selectedKey := m.activityKeys[m.selectedIndex]
				today := time.Now().Format("2006-01-02")
				if err := m.habitManager.AddEntry(selectedKey, today); err == nil {
					// Reload activities and regenerate grid
					m.activities = m.habitManager.GetActivities()
					m.grid = generateGrid(m.activities, m.timeline)
					// Update list items
					m.updateListItems()
				}
			}
		}

		// Global keybindings (work in all views)
		if key.Matches(msg, m.keys.Timeline3m) {
			m.timeline = Timeline3Months
			m.grid = generateGrid(m.activities, m.timeline)
		}
		if key.Matches(msg, m.keys.Timeline6m) {
			m.timeline = Timeline6Months
			m.grid = generateGrid(m.activities, m.timeline)
		}
		if key.Matches(msg, m.keys.Timeline12m) {
			m.timeline = Timeline12Months
			m.grid = generateGrid(m.activities, m.timeline)
		}
		if key.Matches(msg, m.keys.ToggleLegend) {
			m.showLegend = !m.showLegend
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.habitList.SetWidth(msg.Width)
		m.habitList.SetHeight(msg.Height - 4) // Leave space for help
		return m, nil
	}
	return m, cmd
}

// updateListItems refreshes the list items with current activity data
func (m *Model) updateListItems() {
	items := make([]list.Item, 0, len(m.activityKeys))
	for _, key := range m.activityKeys {
		items = append(items, HabitItem{
			key:      key,
			activity: m.activities[key],
		})
	}
	m.habitList.SetItems(items)
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Handle habit selection view
	if m.viewMode == HabitSelection {
		var s strings.Builder
		s.WriteString(m.habitList.View())
		
		// Add help at bottom
		s.WriteString("\n")
		if m.showHelp {
			s.WriteString(m.help.View(m.keys))
		} else {
			s.WriteString(m.help.ShortHelpView([]key.Binding{m.keys.Enter, m.keys.Escape, m.keys.Help, m.keys.Quit}))
		}
		return s.String()
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
		if len(m.activityKeys) > 0 {
			selectedActivity := m.activities[m.activityKeys[m.selectedIndex]]
			titleText = fmt.Sprintf("Activity Tracker - %s (%s)", selectedActivity.Name, timelineText)
		} else {
			titleText = fmt.Sprintf("Activity Tracker (%s)", timelineText)
		}
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
	} else if m.viewMode == SingleActivity {
		// Show single selected activity
		if len(m.activityKeys) > 0 {
			key := m.activityKeys[m.selectedIndex]
			activity := m.activities[key]
			s.WriteString(m.renderActivityGrid(activity, key, -1)) // -1 means no number
		}
	}

	// Legend (right-aligned to grid end) - only show if enabled
	if m.showLegend && m.viewMode != HabitSelection {
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

	// Footer with help/controls
	s.WriteString("\n\n")
	if m.showHelp {
		s.WriteString(m.help.View(m.keys))
	} else {
		// Show appropriate short help based on view mode
		var helpKeys []key.Binding
		if m.viewMode == AllActivities {
			helpKeys = []key.Binding{m.keys.Tab, m.keys.Timeline3m, m.keys.ToggleLegend, m.keys.Help, m.keys.Quit}
		} else {
			helpKeys = []key.Binding{m.keys.Up, m.keys.Enter, m.keys.AllView, m.keys.Tab, m.keys.Help, m.keys.Quit}
		}
		s.WriteString(m.help.ShortHelpView(helpKeys))
	}

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