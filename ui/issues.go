package ui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type IssuesModel struct {
	common     *commonModel
	issues     []jira.Issue
	issuesList list.Model
}

func newIssuesModel(common *commonModel, boardId int) IssuesModel {
	sprints, err := common.jc.GetSprintList(boardId)
	if err != nil {
		fmt.Println(err)
	}

	issues, err := common.jc.GetIssuesForSprint(sprints[0].ID)
	if err != nil {
		fmt.Println(err)
	}

	var issueItems []list.Item
	issueItems = append(issueItems, item("test"))

	/* for _, issue := range issues {
		issueItems = append(issueItems, item(issue.Fields.Summary))
	} */

	issuesList := list.NewModel(issueItems, itemDelegate{}, common.width, common.height)

	issuesList.Title = "Select your Issue"
	issuesList.SetShowStatusBar(false)
	issuesList.SetShowPagination(false)
	issuesList.SetFilteringEnabled(false)

	return IssuesModel{
		common:     common,
		issues:     issues,
		issuesList: issuesList,
	}

}

func (m IssuesModel) View() string {
	return m.issuesList.View()
}

func (m IssuesModel) Update(msg tea.Msg) (IssuesModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.issuesList.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	m.issuesList, cmd = m.issuesList.Update(msg)
	return m, cmd
}
