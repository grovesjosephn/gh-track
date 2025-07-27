package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"hab/internal"
)

var (
	color        string
	targetPerDay int
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [habit-name]",
	Short: "Create a new habit",
	Long: `Create a new habit to track. You can specify the habit name as an argument
or you'll be prompted to enter it interactively.

Examples:
  hab new exercise              # Create a habit called 'exercise'
  hab new --color red exercise  # Create with red color
  hab new --target 2 brushing   # Create with target of 2 times per day`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hm := internal.NewHabitManager()
		if err := hm.Load(); err != nil {
			fmt.Printf("Error loading habits: %v\n", err)
			os.Exit(1)
		}

		var habitName string
		var habitKey string

		if len(args) > 0 {
			habitKey = args[0]
			habitName = strings.Title(strings.ReplaceAll(habitKey, "_", " "))
		} else {
			// Interactive mode
			fmt.Print("Enter habit name: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error reading input: %v\n", err)
				os.Exit(1)
			}
			habitName = strings.TrimSpace(input)
			habitKey = strings.ToLower(strings.ReplaceAll(habitName, " ", "_"))
		}

		if habitName == "" {
			fmt.Println("Error: habit name cannot be empty")
			os.Exit(1)
		}

		// Interactive prompts for missing values
		if color == "" {
			color = promptForColor()
		}

		if targetPerDay == 0 {
			targetPerDay = promptForTarget()
		}

		// Validate color
		validColors := []string{"red", "blue", "green", "magenta", "cyan", "yellow"}
		colorValid := false
		for _, validColor := range validColors {
			if color == validColor {
				colorValid = true
				break
			}
		}
		if !colorValid {
			fmt.Printf("Error: invalid color '%s'. Valid colors: %s\n", color, strings.Join(validColors, ", "))
			os.Exit(1)
		}

		// Create the habit
		if err := hm.CreateActivity(habitKey, habitName, color, targetPerDay); err != nil {
			fmt.Printf("Error creating habit: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Created habit '%s' with color %s", habitName, color)
		if targetPerDay > 1 {
			fmt.Printf(" (target: %d times per day)", targetPerDay)
		}
		fmt.Println()
		fmt.Printf("Add an entry with: hab %s\n", habitKey)
	},
}

func promptForColor() string {
	fmt.Print("Choose color (red, blue, green, magenta, cyan, yellow) [green]: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "green"
	}
	color := strings.TrimSpace(input)
	if color == "" {
		return "green"
	}
	return color
}

func promptForTarget() int {
	fmt.Print("Target per day [1]: ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 1
	}
	targetStr := strings.TrimSpace(input)
	if targetStr == "" {
		return 1
	}
	target, err := strconv.Atoi(targetStr)
	if err != nil || target < 1 {
		return 1
	}
	return target
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&color, "color", "c", "", "Color for the habit (red, blue, green, magenta, cyan, yellow)")
	newCmd.Flags().IntVarP(&targetPerDay, "target", "t", 0, "Target number of times per day")
}