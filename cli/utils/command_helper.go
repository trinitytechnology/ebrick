package utils

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Helper function to get yes/no input
func GetYesOrNoInput(prompt string, defaultYes bool) bool {
	defaultPrompt := "Y/n"
	if !defaultYes {
		defaultPrompt = "y/N"
	}

	for {
		input := GetUserInput(prompt+" ["+defaultPrompt+"]: ", false, "")
		input = strings.ToLower(strings.TrimSpace(input))

		if input == "" { // User pressed Enter, use default
			return defaultYes
		} else if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'Y' or 'N'.")
		}
	}
}

// Helper function to gather user input
// If required is true, it will keep asking for input until it gets a non-empty string. If requiredMsg is provided, it will print the message when input is required
func GetUserInput(prompt string, required bool, requiredMsg string) string {
	for {
		fmt.Print(prompt)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		if required && input == "" {
			if requiredMsg != "" {
				fmt.Println(requiredMsg)
			} else {
				fmt.Println("Input required.")
			}
		} else {
			return strings.TrimSpace(scanner.Text())
		}
	}
}

// Helper function to process modules input and remove duplicates
func ProcessSlicesInput(input string) []string {
	if input == "" {
		return []string{}
	}
	modules := strings.Split(strings.ReplaceAll(input, " ", ""), ",")
	slices.Sort(modules)
	return slices.Compact(modules)
}
