package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	toggl "github.com/chamorim21/toggl-cli/internal/toggl"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a new time entry",
	Long: `Start a new time entry with a description. 
				 If -l (last) flag is provided, the description of the last time entry will be used. 
				 Otherwise, the description must be provided as an argument.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		t, err := toggl.New()
		if err != nil {
			color.Red("%s", err)
			return
		}
		l, _ := cmd.Flags().GetBool("last")
		var description string
		if l {
			d, err := t.LastTimeEntryDescription()
			if err != nil {
				color.Red("Error getting last time entry description\n")
				return
			}
			if d == "" {
				color.Red("No time entry found\n")
				return
			}
			description = d
		} else {
			if len(args) != 1 {
				color.Red("Please provide only one argument as description for the time entry\n")
				return
			}
			description = args[0]
		}
		err = t.Start(description)
		if err != nil {
			color.Red("%s", err)
			return
		}
		color.Green("Time entry started: %s", description)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("last", "l", false, "Use the last time entry description as the description for the new time entry")
}
