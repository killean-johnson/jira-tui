package ui

import (
	"github.com/killean-johnson/jira-tui/api"

	// "github.com/killean-johnson/jira-tui/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type uiState int

const (
	stateShowBoards uiState = iota
	stateShowIssues
)

type commonModel struct {
	jc        *api.JiraClient
	width     int
	height    int
	testState string
}

type Tui struct {
	common      *commonModel
	uiState     uiState
	boardModel  BoardModel
	issuesModel IssuesModel
}

type updateModelState struct {
	state uiState
}

type newIssuesList struct {
	boardId int
}

func InitialTui(client *api.JiraClient) Tui {
	common := commonModel{
		jc:        client,
		width:     80,
		testState: "boards",
	}
	return Tui{
		common:     &common,
		uiState:    stateShowBoards,
		boardModel: newBoardModel(&common),
	}
}

func (t Tui) Init() tea.Cmd {
	return nil
}

func (t Tui) View() string {
	switch t.uiState {
	case stateShowBoards:
		return t.boardModel.View()
	default:
		return "issues"
	}
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.common.width = msg.Width
		t.common.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return t, tea.Quit
		}
	case updateModelState:
		t.common.testState = "issues"
		t.uiState = msg.state

	case newIssuesList:
		t.issuesModel = newIssuesModel(t.common, msg.boardId)
		t.uiState = stateShowIssues
	}
	switch t.uiState {
	cas estateShowBoards:
		newBoardsModel, cmd := t.boardModel.Update(msg)
		t.boardModel = newBoardsModel
		cmds = append(cmds, cmd)
	case stateShowIssues:
		newIssuesModel, cmd := t.issuesModel.Update(msg)
		t.issuesModel = newIssuesModel
		cmds = append(cmds, cmd)
	}

	return t, tea.Batch(cmds...)
}
