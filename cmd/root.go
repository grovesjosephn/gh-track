package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"hab/ui"
)

var (
	interactiveMode bool
	timelineFlag    string
	hideLegend      bool
	version         = "dev"
)

// SetVersion sets the version for the application
func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hab",
	Short: "A terminal-based habit tracker with GitHub-style contribution grids",
	Long: `hab is a fast, terminal-based habit tracker that visualizes your daily habits 
using GitHub-style contribution grids. Track multiple habits with different frequencies
and view your progress over time.

Examples:
  hab                    # Launch interactive TUI (default: 12 months)
  hab -i                 # Launch interactive TUI explicitly  
  hab -t 3m              # Launch TUI with 3 month timeline
  hab --timeline 6m      # Launch TUI with 6 month timeline
  hab --no-legend        # Launch TUI without legend
  hab new exercise       # Create a new habit called 'exercise'
  hab exercise           # Add an entry for 'exercise' today
  hab list               # List all habits with statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no arguments provided, or -i flag used, launch TUI
		if len(args) == 0 || interactiveMode {
			// Parse timeline flag
			var timeline ui.TimelineDays
			switch timelineFlag {
			case "3m", "3":
				timeline = ui.Timeline3Months
			case "6m", "6":
				timeline = ui.Timeline6Months
			case "12m", "1y", "y", "12", "":
				timeline = ui.Timeline12Months
			default:
				fmt.Fprintf(os.Stderr, "Invalid timeline '%s'. Use 3m, 6m, or 12m\n", timelineFlag)
				os.Exit(1)
			}
			ui.RunTUIWithOptions(timeline, !hideLegend)
			return
		}

		// Handle single argument - treat as "add entry for habit"
		if len(args) == 1 {
			addEntry(args[0], "")
			return
		}

		// Show help for invalid usage
		cmd.Help()
	},
	// Custom command validation to handle habit names
	Args: cobra.ArbitraryArgs,
	// Disable unknown command suggestions to allow habit names as arguments
	DisableSuggestions: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Add the interactive flag
	rootCmd.Flags().BoolVarP(&interactiveMode, "interactive", "i", false, "Launch interactive TUI mode")
	
	// Add timeline flag
	rootCmd.Flags().StringVarP(&timelineFlag, "timeline", "t", "12m", "Timeline to display (3m, 6m, 12m)")
	
	// Add legend visibility flag
	rootCmd.Flags().BoolVar(&hideLegend, "no-legend", false, "Hide the completion legend")
}