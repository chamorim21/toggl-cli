package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "toggl-cli",
	Short: "A smooth CLI for Toggl time tracking",
	Long: `Toggl-cli is a CLI tool to interact with the Toggl time tracking service.
You can start and stop time entries.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
