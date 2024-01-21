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

type TimeEntry struct {
	Id              int        `json:"id"`
	WorkspaceId     int        `json:"workspace_id"`
	ProjectId       *int       `json:"project_id"`
	TaskId          *int       `json:"task_id"`
	Billable        bool       `json:"billable"`
	Start           time.Time  `json:"start"`
	Stop            *time.Time `json:"stop"`
	Duration        int        `json:"duration"`
	Description     string     `json:"description"`
	Tags            []string   `json:"tags"`
	TagIds          []int      `json:"tag_ids"`
	Duronly         bool       `json:"duronly"`
	At              time.Time  `json:"at"`
	ServerDeletedAt *time.Time `json:"server_deleted_at"`
	UserId          int        `json:"user_id"`
	Uid             int        `json:"uid"`
	Wid             int        `json:"wid"`
}

func (t *Toggl) getCurrentEntry() (*TimeEntry, error) {
	url := "https://api.track.toggl.com/api/v9/me/time_entries/current"
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	username := t.ApiToken
	password := "api_token"
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	var body TimeEntry
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("error: %s", err)
	}
	defer resp.Body.Close()
	return &body, nil
}

func (t *Toggl) Stop() (string, error) {
	t.Setup()
	c, err := t.getCurrentEntry()
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	if c.Id == 0 {
		return "", fmt.Errorf("no time entry running, please start tracking")
	}
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/time_entries/%d/stop", t.WorkspaceId, c.Id)
	client := &http.Client{}
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	username := t.ApiToken
	password := "api_token"
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	var s TimeEntry
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return "", fmt.Errorf("error: %s", err)
	}
	parsedTime := s.Stop.Local().Format(time.RFC3339)
	return fmt.Sprintf("Time entry stopped at %s: %s", parsedTime, c.Description), nil
}
