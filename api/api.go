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
func (jc *JiraClient) GetBoardList() ([]jira.Board, error) {
	boards, _, err := jc.client.Board.GetAllBoards(nil)
	if err != nil {
		return nil, err
	}
	return boards.Values, nil
}

func (jc *JiraClient) GetSprintList(boardId int) ([]jira.Sprint, error) {
	sprints, _, err := jc.client.Board.GetAllSprintsWithOptions(boardId, &jira.GetAllSprintsOptions{
		State: "active",
	})
	if err != nil {
		return nil, err
	}
	return sprints.Values, nil
}

// get all statuses that a jira card could be in
func (jc *JiraClient) GetStatusList() ([]jira.StatusCategory, error) {
	statuses, _, err := jc.client.StatusCategory.GetList()
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// Get all the issues on a sprint
func (jc *JiraClient) GetIssuesForSprint(sprintId int) ([]jira.Issue, error) {
    issues, _, err := jc.client.Issue.Search("Sprint=" + fmt.Sprint(sprintId), &jira.SearchOptions {
        Fields: []string{"summary", "status", "assignee"},
    })
	if err != nil {
		return nil, err
	}
    return issues, nil
}

func (jc *JiraClient) GetIssue(issueId string) (*jira.Issue, error) {
    issue, _, err := jc.client.Issue.Get(issueId, nil)
    if err != nil {
        return nil, err
    }
    return issue, nil
}

func (jc *JiraClient) UpdateIssue(issue *jira.Issue) error {
    _, _, err := jc.client.Issue.Update(issue) 
    return err
}

func (jc *JiraClient) DoTransition(issueKey string, transitionName string) error {
    var transitionId string
    transitions, _, _ := jc.client.Issue.GetTransitions(issueKey)
    for _, t := range(transitions) {
        if t.Name == transitionName {
            transitionId = t.ID
            break
        }
    }
    _, err := jc.client.Issue.DoTransition(issueKey, transitionId)
    return err
}
