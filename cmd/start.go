/*
Copyright © 2024 Cristiano Amorim cristianoretiro2003@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	internal "github.com/chamorim21/toggl-cli/internal"
	trackers "github.com/chamorim21/toggl-cli/internal/trackers"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("Please provide a description for the time entry\n")
			return
		}
		if len(args) > 1 {
			fmt.Printf("Please provide only one argument\n")
			return
		}
		timeEntryDescription := args[0]
		tracker := &trackers.Toggl{}
		start(tracker, timeEntryDescription)
	},
}

func start(t internal.Tracker, description string) {
	loading := true
	go func() {
		if loading {
			color.Cyan("Starting tracking...")
		}
	}()
	err := t.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := t.Start(description)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	color.Green("%s", s)
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
