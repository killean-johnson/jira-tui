package api

import (
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
)

type JiraClient struct {
	client *jira.Client
}

func (jc *JiraClient) Connect(username string, token string) {
	// Set up client
	authTransport := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client, err := jira.NewClient(authTransport.Client(), os.Getenv("JIRA_BOARD_URL"))

	if err != nil {
		fmt.Println("ERR IN createClient FUNC")
	} else {
        jc.client = client
    }
}

// Get all projects
func (jc *JiraClient) GetProjectList() (*jira.ProjectList, error) {
	projects, _, err := jc.client.Project.GetList()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// get all boards that exist on project
// loads project from userConfig
func (jc *JiraClient) GetBoardList() (*jira.BoardsList, error) {
    boards, _, err := jc.client.Board.GetAllBoards(nil)
    if err != nil {
        return nil, err
    }
	return boards, nil
}

func (jc *JiraClient) GetSprintList(boardId string) ([]jira.Sprint, error) {
    sprints, _, err := jc.client.Board.GetAllSprints(boardId)
    if err != nil {
        return nil, err
    }
	return sprints, nil
}

// get all statuses that a jira card could be in
func (jc *JiraClient) GetStatusList() ([]jira.Status, error) {
    statuses, _, err := jc.client.Status.GetAllStatuses()
    if err != nil {
        return nil, err
    }
	return statuses, nil
}
