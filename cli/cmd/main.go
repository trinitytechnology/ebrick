package main

import (
	"github.com/spf13/cobra"
	"github.com/trinitytechnology/ebrick/cli/app"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "ebrick",
	}

	var newAppCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new ebrick application",
		Run: func(cmd *cobra.Command, args []string) {
			app.NewApp()
		},
	}

	rootCmd.AddCommand(newAppCmd)
	rootCmd.Execute()
}
