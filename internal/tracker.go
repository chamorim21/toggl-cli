package internal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Tracker interface {
	Setup() error
	Start() (string, error)
	Current() (string, error)
}

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
