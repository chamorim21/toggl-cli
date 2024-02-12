package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	toggl "github.com/chamorim21/toggl-cli/internal/toggl"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the current time entry",
	Long:  `Stop the current time entry. If there is no time entry running, the command fails.`,
	Run: func(cmd *cobra.Command, args []string) {
		t, err := toggl.New()
		if err != nil {
			color.Red("%s", err)
			return
		}
		err = t.Stop()
		if err != nil {
			color.Red("%s", err)
			return
		}
		color.Green("Time entry stopped")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
