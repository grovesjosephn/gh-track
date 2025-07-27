package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"hab/internal"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all habits with statistics",
	Long: `List all habits with their current statistics including total entries,
unique days tracked, and current streak.

Example:
  hab list`,
	Run: func(cmd *cobra.Command, args []string) {
		hm := internal.NewHabitManager()
		if err := hm.Load(); err != nil {
			fmt.Printf("Error loading habits: %v\n", err)
			os.Exit(1)
		}

		activities := hm.GetActivities()
		if len(activities) == 0 {
			fmt.Println("No habits found. Create one with: hab new [habit-name]")
			return
		}

		// Sort habit keys
		var keys []string
		for key := range activities {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		fmt.Println("Your Habits:")
		fmt.Println("============")

		for i, key := range keys {
			activity := activities[key]
			stats, err := hm.GetStats(key)
			if err != nil {
				fmt.Printf("Error getting stats for %s: %v\n", key, err)
				continue
			}

			fmt.Printf("\n[%d] %s (%s)\n", i+1, activity.Name, activity.Color)
			fmt.Printf("    Key: %s\n", key)
			fmt.Printf("    Total entries: %d\n", stats["total_entries"])
			fmt.Printf("    Unique days: %d\n", stats["unique_days"])
			fmt.Printf("    Target per day: %d\n", stats["target_per_day"])
			
			if streak, ok := stats["current_streak"].(int); ok {
				if streak > 0 {
					fmt.Printf("    Current streak: %d days 🔥\n", streak)
				} else {
					fmt.Printf("    Current streak: 0 days\n")
				}
			}

			fmt.Printf("    Add entry: hab %s\n", key)
		}

		fmt.Printf("\nTotal habits: %d\n", len(activities))
		fmt.Println("\nUse 'hab' to view the interactive grid, or 'hab [habit]' to add an entry.")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}