package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"hab/ui"
)

var interactiveMode bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hab",
	Short: "A terminal-based habit tracker with GitHub-style contribution grids",
	Long: `hab is a fast, terminal-based habit tracker that visualizes your daily habits 
using GitHub-style contribution grids. Track multiple habits with different frequencies
and view your progress over time.

Examples:
  hab                    # Launch interactive TUI (default)
  hab -i                 # Launch interactive TUI explicitly  
  hab new exercise       # Create a new habit called 'exercise'
  hab exercise           # Add an entry for 'exercise' today
  hab list               # List all habits with statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		// If no arguments provided, or -i flag used, launch TUI
		if len(args) == 0 || interactiveMode {
			ui.RunTUI()
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
}