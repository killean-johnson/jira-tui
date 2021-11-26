package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
	"github.com/killean-johnson/jira-tui/api"

	// "github.com/killean-johnson/jira-tui/tui"
	tea "github.com/charmbracelet/bubbletea"
)

type boardList struct {
	boards []jira.Board
	cursor int
}

type Tui struct {
	client      *api.JiraClient
	boards      []jira.Board
	chosenBoard int
	issues      []jira.Issue
	chosenIssue jira.Issue
}

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func initialBoards(client *api.JiraClient) boardList {
	boards, err := client.GetBoardList()
	if err != nil {
		fmt.Println(err)
	}
	return boardList{
		boards: boards,
		cursor: 0,
	}
}

func (t Tui) Init() tea.Cmd {
	/* boards, err := client.GetBoardList()
	   bl.boards = boards
	   if err != nil {
	       fmt.Println(err)
	   } */
	return nil
}

func chooseBoardView(t Tui) string {
	s := "Which Board"
	for i, board := range t.boards {
		cursor := " "
		if t.boards[i].ID == t.chosenBoard {
			cursor = ">"
		}
		s += fmt.Sprintf("\n%s %s", cursor, board.Name)
	}
	return s
}

func (t Tui) View() string {
	var s string
	if t.chosenBoard == 0 {
		s = chooseBoardView(t)
	}

	return s
}

/* func (bl boardList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return bl, tea.Quit
		case "k":
			if bl.cursor > 0 {
				bl.cursor--
			}
		case "j":
			if bl.cursor < len(bl.boards)-1 {
				bl.cursor++
			}
		case "enter":

		}

	}
	return bl, nil
} */

func main() {
	godotenv.Load()
	jiraToken := os.Getenv("JIRA_API_TOKEN")

	client := &api.JiraClient{}
	client.Connect("killean.johnson@stairsupplies.com", jiraToken)

	p := tea.NewProgram(initialBoards(client), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error has occured: %v", err)
		os.Exit(1)
	}

	// tui.CreateGUI(client)
}
