package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"hab/internal"
)

var forceDelete bool

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [habit]",
	Short: "Delete a habit",
	Long: `Delete a habit permanently. This will remove all data for the habit.
By default, you'll be prompted for confirmation unless --force is used.

Examples:
  hab delete exercise      # Delete with confirmation
  hab delete exercise -f   # Delete without confirmation`,
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

		// Confirmation prompt unless --force is used
		if !forceDelete {
			fmt.Printf("Are you sure you want to delete habit '%s'? This will remove all %d entries. (y/N): ", 
				activity.Name, len(activity.Dates))
			
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				os.Exit(1)
			}

			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Operation cancelled")
				return
			}
		}

		// Delete the habit
		if err := hm.DeleteActivity(habitKey); err != nil {
			fmt.Printf("Error deleting habit: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Deleted habit '%s'\n", activity.Name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Delete without confirmation")
}