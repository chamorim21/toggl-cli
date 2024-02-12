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
	ProjectId   int
}

func (t *Toggl) setup() error {
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
	projectId := os.Getenv("TOGGL_PROJECT_ID")
	if projectId == "" {
		return fmt.Errorf("error: TOGGL_PROJECT_ID is not set")
	}
	t.ApiToken = togglApiToken
	wId, err := strconv.Atoi(workspaceId)
	if err != nil {
		return fmt.Errorf("error: TOGGL_WORKSPACE_ID is not a number")
	}
	t.WorkspaceId = wId
	pId, err := strconv.Atoi(projectId)
	if err != nil {
		return fmt.Errorf("error: TOGGL_PROJECT_ID is not a number")
	}
	t.ProjectId = pId
	return nil
}

func New() (*Toggl, error) {
	t := &Toggl{}
	err := t.setup()
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Toggl) newRequest(url string, method string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	username := t.ApiToken
	password := "api_token"
	req.SetBasicAuth(username, password)
	return req
}

func (t *Toggl) Start(description string) error {
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
		ProjectId   int      `json:"project_id"`
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
		ProjectId:   t.ProjectId,
	}
	m, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/time_entries", t.WorkspaceId)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(m))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	username := t.ApiToken
	password := "api_token"
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return nil
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

func (t *Toggl) currentTimeEntryId() (int, error) {
	url := "https://api.track.toggl.com/api/v9/me/time_entries/current"
	client := &http.Client{}
	req := t.newRequest(url, http.MethodGet)
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	var body TimeEntry
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return body.Id, nil
}

func (t *Toggl) Stop() error {
	eId, err := t.currentTimeEntryId()
	if err != nil {
		return err
	}
	if eId == 0 {
		return fmt.Errorf("no time entry running, please start tracking")
	}
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/time_entries/%d/stop", t.WorkspaceId, eId)
	req := t.newRequest(url, http.MethodPut)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	e := &TimeEntry{}
	err = json.NewDecoder(resp.Body).Decode(e)
	if err != nil {
		return err
	}
	return nil
}

func (t *Toggl) LastTimeEntryDescription() (string, error) {
	req := t.newRequest(http.MethodGet,
		"https://api.track.toggl.com/api/v9/me/time_entries")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var e []TimeEntry
	err = json.NewDecoder(resp.Body).Decode(&e)
	if err != nil {
		return "", err
	}
	return e[0].Description, nil
}
