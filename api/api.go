package api

import (
	"fmt"
	"net/url"

	"github.com/andygrunwald/go-jira"
)

type JiraClient struct {
	client *jira.Client
}

func (jc *JiraClient) GetConnection() *jira.Client {
	return jc.client
}

func (jc *JiraClient) Connect(username string, token string, url string) {
	// Set up client
	authTransport := jira.BasicAuthTransport{
		Username: username,
		Password: token,
	}

	client, err := jira.NewClient(authTransport.Client(), url)

	if err != nil {
		fmt.Println("ERR IN createClient FUNC")
	} else {
		jc.client = client
	}
}

// TODO: Fix this, it doesn't actually work
func (jc *JiraClient) GetProject(projectKey string) (*jira.Project, error) {
    project, _, err := jc.client.Project.Get(projectKey)
    if err != nil {
        return nil, err
    }

    return project, nil
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

func (jc *JiraClient) GetBoard(boardId int) (*jira.Board, error) {
	board, _, err := jc.client.Board.GetBoard(boardId)
	if err != nil {
		return nil, err
	}
	return board, nil
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
func (jc *JiraClient) GetStatusList() ([]jira.Status, error) {
	statuses, _, err := jc.client.Status.GetAllStatuses()
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

func (jc *JiraClient) GetIssueTransitions(id string) ([]jira.Transition, error) {
    transitions, _, err := jc.client.Issue.GetTransitions(id)
    if err != nil {
        return nil, err
    }
    return transitions, nil
}

// Get all the issues on a sprint
func (jc *JiraClient) GetIssuesForSprint(sprintId int) ([]jira.Issue, error) {
	issues, _, err := jc.client.Issue.Search("Sprint="+fmt.Sprint(sprintId), &jira.SearchOptions{
		Fields: []string{"summary", "status", "assignee", "description", "sprint", "project"},
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
	for _, t := range transitions {
		if t.Name == transitionName {
			transitionId = t.ID
			break
		}
	}
	_, err := jc.client.Issue.DoTransition(issueKey, transitionId)
	return err
}

func (jc *JiraClient) CreateIssue(sprintId int, issue *jira.Issue) (*jira.Issue, error) {
	is, _, err := jc.client.Issue.Create(issue)
	if err != nil {
		return nil, err
	}

	_, err = jc.client.Sprint.MoveIssuesToSprint(sprintId, []string{is.Key})
	if err != nil {
		return nil, err
	}

	return is, nil
}

func (jc *JiraClient) GetActiveSprint(boardId int) (*jira.Sprint, error) {
	sprints, _, err := jc.client.Board.GetAllSprintsWithOptions(boardId, &jira.GetAllSprintsOptions{
		State: "active",
	})
	if err != nil {
		return nil, err
	}
	return &sprints.Values[0], nil
}

func (jc *JiraClient) GetUsers(boardKey string) (*[]jira.User, error) {
	u := url.URL{
		Path: "/rest/api/2/user/assignable/multiProjectSearch",
	}
	uv := url.Values{}
	uv.Add("projectKeys", boardKey)
	uv.Add("startAt", "0")
	uv.Add("maxResults", "50")

	u.RawQuery = uv.Encode()

	req, _ := jc.client.NewRequest("GET", u.String(), nil)

	users := new([]jira.User)
	_, err := jc.client.Do(req, users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
