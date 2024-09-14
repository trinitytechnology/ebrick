package main

import (
	"github.com/spf13/cobra"
	"github.com/trinitytechnology/ebrick/cli/app"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "ebrick",
	}

	rootCmd.AddCommand(createAppCommands())
	rootCmd.AddCommand(createRunCommand())
	rootCmd.Execute()
}

func createAppCommands() *cobra.Command {
	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new ebrick application, module or service..",
	}

	var newAppCmd = &cobra.Command{
		Use:   "app",
		Short: "Create a new ebrick application",
		Run: func(cmd *cobra.Command, args []string) {
			app.NewApp()
		},
	}
	newCmd.AddCommand(newAppCmd)

	return newCmd
}

func createRunCommand() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the ebrick application",
		Run: func(cmd *cobra.Command, args []string) {
			app.RunApp()
		},
	}
	return runCmd
}
