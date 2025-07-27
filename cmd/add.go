package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"hab/internal"
)

var dateFlag string

// addEntry is the shared function for adding entries
func addEntry(habitKey, date string) {
	hm := internal.NewHabitManager()
	if err := hm.Load(); err != nil {
		fmt.Printf("Error loading habits: %v\n", err)
		os.Exit(1)
	}

	// Check if habit exists
	if _, exists := hm.GetActivity(habitKey); !exists {
		fmt.Printf("Error: habit '%s' does not exist\n", habitKey)
		fmt.Println("Create it first with: hab new " + habitKey)
		os.Exit(1)
	}

	// Use provided date or default to today
	entryDate := date
	if entryDate == "" {
		entryDate = time.Now().Format("2006-01-02")
	}

	// Add the entry
	if err := hm.AddEntry(habitKey, entryDate); err != nil {
		fmt.Printf("Error adding entry: %v\n", err)
		os.Exit(1)
	}

	// Get habit info for confirmation
	activity, _ := hm.GetActivity(habitKey)
	
	if entryDate == time.Now().Format("2006-01-02") {
		fmt.Printf("✓ Added entry for '%s' today\n", activity.Name)
	} else {
		fmt.Printf("✓ Added entry for '%s' on %s\n", activity.Name, entryDate)
	}

	// Show current streak if available
	if stats, err := hm.GetStats(habitKey); err == nil {
		if streak, ok := stats["current_streak"].(int); ok && streak > 0 {
			fmt.Printf("Current streak: %d days 🔥\n", streak)
		}
	}
}

// addCmd represents the add command (this is actually handled by the root command for convenience)
var addCmd = &cobra.Command{
	Use:   "add [habit] [date]",
	Short: "Add an entry for a habit",
	Long: `Add an entry for a habit. If no date is specified, today's date is used.

Examples:
  hab add exercise           # Add entry for today
  hab add exercise 2025-01-15  # Add entry for specific date`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		habitKey := args[0]
		date := ""
		if len(args) > 1 {
			date = args[1]
		}
		if dateFlag != "" {
			date = dateFlag
		}
		addEntry(habitKey, date)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&dateFlag, "date", "d", "", "Date to add entry for (YYYY-MM-DD)")
}