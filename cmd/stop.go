/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	internal "github.com/chamorim21/toggl-cli/internal"
	trackers "github.com/chamorim21/toggl-cli/internal/trackers"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		tracker := &trackers.Toggl{}
		stop(tracker)
	},
}

func stop(t internal.Tracker) {
	loading := true
	go func() {
		if loading {
			color.Cyan("Stopping tracking...")
		}
	}()
	err := t.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	s, err := t.Stop()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	color.Green("%s", s)
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
