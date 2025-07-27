package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"hab/internal"
)

var (
	pruneForce bool
	pruneDryRun bool
)

// pruneCmd represents the prune command
var pruneCmd = &cobra.Command{
	Use:   "prune [habit-key]",
	Short: "Remove excess entries that exceed the target per day for habits",
	Long: `Remove excess entries that exceed the target per day for each habit.
	
For each day where a habit has more entries than its target_per_day setting,
this command will remove the excess entries, keeping only the number specified
by target_per_day.

Examples:
  hab prune                    # Prune all habits (with confirmation)
  hab prune exercise           # Prune only the 'exercise' habit
  hab prune --dry-run          # Show what would be pruned without making changes
  hab prune --force            # Prune without confirmation prompts`,
	Run: func(cmd *cobra.Command, args []string) {
		hm := internal.NewHabitManager()
		if err := hm.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Error loading habits: %v\n", err)
			os.Exit(1)
		}

		activities := hm.GetActivities()
		if len(activities) == 0 {
			fmt.Println("No habits found to prune.")
			return
		}

		// Determine which habits to prune
		var habitsToProcess []string
		if len(args) == 1 {
			// Specific habit provided
			habitKey := args[0]
			if _, exists := activities[habitKey]; !exists {
				fmt.Fprintf(os.Stderr, "Habit '%s' not found.\n", habitKey)
				os.Exit(1)
			}
			habitsToProcess = append(habitsToProcess, habitKey)
		} else {
			// All habits
			for key := range activities {
				habitsToProcess = append(habitsToProcess, key)
			}
		}

		totalPruned := 0
		for _, habitKey := range habitsToProcess {
			activity := activities[habitKey]
			pruned, err := pruneHabit(hm, habitKey, activity, pruneDryRun, pruneForce)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error pruning habit '%s': %v\n", habitKey, err)
				continue
			}
			totalPruned += pruned
		}

		if pruneDryRun {
			fmt.Printf("\nDry run complete. Would prune %d total entries.\n", totalPruned)
			fmt.Println("Run without --dry-run to actually remove entries.")
		} else if totalPruned > 0 {
			fmt.Printf("\nSuccessfully pruned %d total entries.\n", totalPruned)
		} else {
			fmt.Println("\nNo excess entries found to prune.")
		}
	},
	Args: cobra.MaximumNArgs(1),
}

func pruneHabit(hm *internal.HabitManager, habitKey string, activity internal.Activity, dryRun, force bool) (int, error) {
	// Count entries per date
	dateCounts := make(map[string]int)
	for _, date := range activity.Dates {
		dateCounts[date]++
	}

	// Determine target (default to 1 if not set)
	target := activity.TargetPerDay
	if target == 0 {
		target = 1
	}

	// Find dates with excess entries
	var excessDates []string
	totalExcess := 0
	for date, count := range dateCounts {
		if count > target {
			excessDates = append(excessDates, date)
			totalExcess += count - target
		}
	}

	if len(excessDates) == 0 {
		return 0, nil
	}

	// Show what will be pruned
	fmt.Printf("\nHabit: %s (target: %d per day)\n", activity.Name, target)
	for _, date := range excessDates {
		excess := dateCounts[date] - target
		fmt.Printf("  %s: %d entries → %d entries (removing %d)\n", 
			date, dateCounts[date], target, excess)
	}

	if dryRun {
		return totalExcess, nil
	}

	// Confirm unless force flag is used
	if !force {
		fmt.Printf("\nThis will remove %d excess entries for '%s'. Continue? (y/N): ", totalExcess, activity.Name)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" && response != "yes" {
			fmt.Println("Cancelled.")
			return 0, nil
		}
	}

	// Remove excess entries
	actualPruned := 0
	for _, date := range excessDates {
		excess := dateCounts[date] - target
		for i := 0; i < excess; i++ {
			if err := hm.RemoveEntry(habitKey, date); err != nil {
				return actualPruned, fmt.Errorf("failed to remove entry for %s: %w", date, err)
			}
			actualPruned++
		}
	}

	return actualPruned, nil
}

func init() {
	rootCmd.AddCommand(pruneCmd)

	pruneCmd.Flags().BoolVar(&pruneForce, "force", false, "Prune without confirmation prompts")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "Show what would be pruned without making changes")
}