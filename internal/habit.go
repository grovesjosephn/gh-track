package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Activity represents a single activity with its metadata
type Activity struct {
	Name         string   `json:"name"`
	Color        string   `json:"color"`
	Dates        []string `json:"dates"`
	TargetPerDay int      `json:"target_per_day,omitempty"` // Optional: defaults to 1
}

// ActivitiesData represents the root JSON structure
type ActivitiesData struct {
	Activities map[string]Activity `json:"activities"`
}

// HabitManager handles all habit-related operations
type HabitManager struct {
	dataFile string
	data     *ActivitiesData
}

// getDefaultDataPath returns the default data file path based on OS
func getDefaultDataPath() string {
	// Check for HAB_DATA_FILE environment variable first
	if dataFile := os.Getenv("HAB_DATA_FILE"); dataFile != "" {
		return dataFile
	}

	var configDir string
	
	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming")
		}
	case "darwin":
		configDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	default: // Linux and other Unix-like systems
		configDir = os.Getenv("XDG_CONFIG_HOME")
		if configDir == "" {
			configDir = filepath.Join(os.Getenv("HOME"), ".config")
		}
	}

	// If we can't determine a config directory, fall back to current directory
	if configDir == "" || os.Getenv("HOME") == "" {
		return "data/activities.json"
	}

	return filepath.Join(configDir, "hab", "data", "activities.json")
}

// NewHabitManager creates a new habit manager with default data file
func NewHabitManager() *HabitManager {
	return &HabitManager{
		dataFile: getDefaultDataPath(),
		data:     &ActivitiesData{Activities: make(map[string]Activity)},
	}
}

// Load reads the activities data from the JSON file
func (hm *HabitManager) Load() error {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(hm.dataFile), 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create file if it doesn't exist
	if _, err := os.Stat(hm.dataFile); os.IsNotExist(err) {
		return hm.Save() // Create empty file with proper structure
	}

	data, err := ioutil.ReadFile(hm.dataFile)
	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}

	if err := json.Unmarshal(data, hm.data); err != nil {
		return fmt.Errorf("failed to parse data file: %w", err)
	}

	return nil
}

// Save writes the activities data to the JSON file
func (hm *HabitManager) Save() error {
	data, err := json.MarshalIndent(hm.data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := ioutil.WriteFile(hm.dataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
}

// GetActivities returns all activities
func (hm *HabitManager) GetActivities() map[string]Activity {
	return hm.data.Activities
}

// GetActivity returns a specific activity by key
func (hm *HabitManager) GetActivity(key string) (Activity, bool) {
	activity, exists := hm.data.Activities[key]
	return activity, exists
}

// CreateActivity creates a new activity
func (hm *HabitManager) CreateActivity(key, name, color string, targetPerDay int) error {
	if _, exists := hm.data.Activities[key]; exists {
		return fmt.Errorf("activity '%s' already exists", key)
	}

	if targetPerDay <= 0 {
		targetPerDay = 1
	}

	hm.data.Activities[key] = Activity{
		Name:         name,
		Color:        color,
		Dates:        []string{},
		TargetPerDay: targetPerDay,
	}

	return hm.Save()
}

// AddEntry adds a date entry to an activity
func (hm *HabitManager) AddEntry(key, dateStr string) error {
	activity, exists := hm.data.Activities[key]
	if !exists {
		return fmt.Errorf("activity '%s' does not exist", key)
	}

	// Validate date format
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return fmt.Errorf("invalid date format '%s', use YYYY-MM-DD", dateStr)
	}

	// Add the date
	activity.Dates = append(activity.Dates, dateStr)
	hm.data.Activities[key] = activity

	return hm.Save()
}

// RemoveEntry removes a date entry from an activity
func (hm *HabitManager) RemoveEntry(key, dateStr string) error {
	activity, exists := hm.data.Activities[key]
	if !exists {
		return fmt.Errorf("activity '%s' does not exist", key)
	}

	// Find and remove the date (only first occurrence)
	for i, date := range activity.Dates {
		if date == dateStr {
			activity.Dates = append(activity.Dates[:i], activity.Dates[i+1:]...)
			hm.data.Activities[key] = activity
			return hm.Save()
		}
	}

	return fmt.Errorf("date '%s' not found in activity '%s'", dateStr, key)
}

// DeleteActivity removes an activity entirely
func (hm *HabitManager) DeleteActivity(key string) error {
	if _, exists := hm.data.Activities[key]; !exists {
		return fmt.Errorf("activity '%s' does not exist", key)
	}

	delete(hm.data.Activities, key)
	return hm.Save()
}

// UpdateActivity updates activity metadata
func (hm *HabitManager) UpdateActivity(key string, name, color string, targetPerDay int) error {
	activity, exists := hm.data.Activities[key]
	if !exists {
		return fmt.Errorf("activity '%s' does not exist", key)
	}

	if name != "" {
		activity.Name = name
	}
	if color != "" {
		activity.Color = color
	}
	if targetPerDay > 0 {
		activity.TargetPerDay = targetPerDay
	}

	hm.data.Activities[key] = activity
	return hm.Save()
}

// GetStats returns statistics for an activity
func (hm *HabitManager) GetStats(key string) (map[string]interface{}, error) {
	activity, exists := hm.data.Activities[key]
	if !exists {
		return nil, fmt.Errorf("activity '%s' does not exist", key)
	}

	stats := make(map[string]interface{})
	stats["name"] = activity.Name
	stats["total_entries"] = len(activity.Dates)
	stats["target_per_day"] = activity.TargetPerDay

	// Calculate unique days (for multi-frequency habits)
	uniqueDays := make(map[string]bool)
	for _, date := range activity.Dates {
		uniqueDays[date] = true
	}
	stats["unique_days"] = len(uniqueDays)

	// Calculate current streak
	stats["current_streak"] = hm.calculateStreak(activity)

	return stats, nil
}

// calculateStreak calculates the current streak for an activity
func (hm *HabitManager) calculateStreak(activity Activity) int {
	if len(activity.Dates) == 0 {
		return 0
	}

	// Get unique days and sort them
	uniqueDays := make(map[string]bool)
	for _, date := range activity.Dates {
		uniqueDays[date] = true
	}

	var sortedDates []string
	for date := range uniqueDays {
		sortedDates = append(sortedDates, date)
	}

	if len(sortedDates) == 0 {
		return 0
	}

	// Sort dates in descending order
	for i := 0; i < len(sortedDates)-1; i++ {
		for j := i + 1; j < len(sortedDates); j++ {
			if sortedDates[i] < sortedDates[j] {
				sortedDates[i], sortedDates[j] = sortedDates[j], sortedDates[i]
			}
		}
	}

	// Calculate streak from the most recent date
	streak := 0
	today := time.Now().Format("2006-01-02")
	currentDate := today

	for _, date := range sortedDates {
		if date == currentDate {
			streak++
			// Move to previous day
			t, _ := time.Parse("2006-01-02", currentDate)
			currentDate = t.AddDate(0, 0, -1).Format("2006-01-02")
		} else {
			break
		}
	}

	return streak
}