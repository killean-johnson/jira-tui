package ui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
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
	jc *api.JiraClient
}

type Tui struct {
	common  *commonModel
	uiState uiState
}

func initialTui(client *api.JiraClient) Tui {
	return Tui{}
}

func (t Tui) Init() tea.Cmd {
	/* boards, err := client.GetBoardList()
	   bl.boards = boards
	   if err != nil {
	       fmt.Println(err)
	   } */
	return nil
}

func (t Tui) View() string {
	var s string
	s = "Hello World"

	return s
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return t, tea.Quit
		}
	}
	return t, nil
}
