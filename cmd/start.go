/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/chamorim21/toggl-cli/internal"
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
		startTimeEntry(timeEntryDescription)
	},
}

func startTimeEntry(description string) {
	var t *internal.Toggl = &internal.Toggl{}
	err := t.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	type body struct {
		CreatedWith string   `json:"created_with"`
		Description string   `json:"description"`
		Tags        []string `json:"tags"`
		Billable    bool     `json:"billable"`
		Duration    int      `json:"duration"`
		Start       string   `json:"start"`
		Stop        *string  `json:"stop"`
		WorkspaceId int      `json:"workspace_id"`
	}
	payload := body{
		CreatedWith: "Toggl CLI",
		Description: description,
		Tags:        []string{},
		Billable:    false,
		Duration:    -1,
		Start:       time.Now().Format(time.RFC3339),
		Stop:        nil,
		WorkspaceId: t.WorkspaceId,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/time_entries", t.WorkspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	username := t.ApiToken
	password := "api_token"
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	var target interface{}
	json.NewDecoder(resp.Body).Decode(&target)
	fmt.Printf("Response: %s\n", target)
	defer resp.Body.Close()
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
