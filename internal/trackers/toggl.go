package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Toggl struct {
	ApiToken    string
	WorkspaceId int
}

func (t *Toggl) Setup() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}
	togglApiToken := os.Getenv("TOGGL_API_TOKEN")
	if togglApiToken == "" {
		return fmt.Errorf("error: TOGGL_API_TOKEN is not set")
	}
	workspaceId := os.Getenv("TOGGL_WORKSPACE_ID")
	if workspaceId == "" {
		return fmt.Errorf("error: TOGGL_WORKSPACE_ID is not set")
	}
	t.ApiToken = togglApiToken
	wId, err := strconv.Atoi(workspaceId)
	if err != nil {
		return fmt.Errorf("error: TOGGL_WORKSPACE_ID is not a number")
	}
	t.WorkspaceId = wId
	return nil
}

func (t *Toggl) Start(description string) (string, error) {
	start := time.Now().Format(time.RFC3339)
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
		Start:       start,
		Stop:        nil,
		WorkspaceId: t.WorkspaceId,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
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
	defer resp.Body.Close()
	return fmt.Sprintf("Time entry started at %s: %s", start, description), nil
}
