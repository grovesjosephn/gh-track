package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"hab/internal"
)

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats [habit]",
	Short: "Show detailed statistics for a habit",
	Long: `Show detailed statistics for a specific habit including total entries,
unique days, current streak, and other metrics.

Examples:
  hab stats exercise    # Show stats for exercise habit`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		habitKey := args[0]

		hm := internal.NewHabitManager()
		if err := hm.Load(); err != nil {
			fmt.Printf("Error loading habits: %v\n", err)
			os.Exit(1)
		}

		// Check if habit exists
		activity, exists := hm.GetActivity(habitKey)
		if !exists {
			fmt.Printf("Error: habit '%s' does not exist\n", habitKey)
			os.Exit(1)
		}

		// Get statistics
		stats, err := hm.GetStats(habitKey)
		if err != nil {
			fmt.Printf("Error getting statistics: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Statistics for '%s'\n", activity.Name)
		fmt.Println(strings.Repeat("=", len(activity.Name)+16))
		fmt.Printf("Key: %s\n", habitKey)
		fmt.Printf("Color: %s\n", activity.Color)
		fmt.Printf("Target per day: %d\n", stats["target_per_day"])
		fmt.Printf("Total entries: %d\n", stats["total_entries"])
		fmt.Printf("Unique days tracked: %d\n", stats["unique_days"])
		
		if streak, ok := stats["current_streak"].(int); ok {
			if streak > 0 {
				fmt.Printf("Current streak: %d days 🔥\n", streak)
			} else {
				fmt.Printf("Current streak: 0 days\n")
			}
		}

		// Calculate success rate for multi-frequency habits
		if targetPerDay, ok := stats["target_per_day"].(int); ok && targetPerDay > 1 {
			totalEntries := stats["total_entries"].(int)
			uniqueDays := stats["unique_days"].(int)
			expectedEntries := uniqueDays * targetPerDay
			if expectedEntries > 0 {
				successRate := float64(totalEntries) / float64(expectedEntries) * 100
				fmt.Printf("Completion rate: %.1f%% (%d/%d expected entries)\n", 
					successRate, totalEntries, expectedEntries)
			}
		}

		fmt.Printf("\nUse 'hab %s' to add an entry for today\n", habitKey)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}