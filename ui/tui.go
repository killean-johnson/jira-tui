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
	jc    *api.JiraClient
	width int
}

type Tui struct {
	common    *commonModel
	uiState   uiState
	boardList BoardModel
}

type updateModelState uiState

func InitialTui(client *api.JiraClient) Tui {
	common := commonModel{
		jc:    client,
		width: 80,
	}
	return Tui{
		common:    &common,
		uiState:   stateShowBoards,
		boardList: newBoardModel(&common),
	}
}

func (t Tui) Init() tea.Cmd {
	return nil
}

func (t Tui) View() string {
	switch t.uiState {
	case stateShowBoards:
		return t.boardList.View()
	default:
		return t.boardList.View()
	}
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.common.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return t, tea.Quit
		}
	case updateModelState:
		t.uiState = stateShowIssues
	}
	switch t.uiState {
	case stateShowBoards:
		newBoardsModel, cmd := t.boardList.Update(msg)
		t.boardList = newBoardsModel
		cmds = append(cmds, cmd)
	}

	return t, tea.Batch(cmds...)
}
