package api

import (
	"fmt"

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

	//TODO: get url from user config
	client, err := jira.NewClient(authTransport.Client(), "https://stairsupplies-voe.atlassian.net")

	if err != nil {
		fmt.Println("ERR IN createClient FUNC")
	}

	jc.client = client
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
func (jc *JiraClient) GetBoardList() ([]jira.BoardsList, error) {
	return nil, nil
}

func (jc *JiraClient) GetSprintList() ([]jira.SprintsList, error) {
	return nil, nil
}

// get all statuses that a jira card could be in
func (jc *JiraClient) GetStatusList() ([]jira.StatusCategory, error) {
	return nil, nil
}
