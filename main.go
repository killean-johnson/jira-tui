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
	client        *api.JiraClient
	boards        []jira.Board
	selectedBoard jira.Board
	chosenBoard   bool
	boardCursor   int
	issues        []jira.Issue
	chosenIssue   jira.Issue
}

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func initialTui(client *api.JiraClient) Tui {
	boards, err := client.GetBoardList()
	if err != nil {
		fmt.Println(err)
	}
	return Tui{
		client:      client,
		boards:      boards,
		chosenBoard: false,
		boardCursor: 0,
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
		if t.boardCursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("\n%s %s", cursor, board.Name)
	}
	return s
}

func boardKeys(msg tea.Msg, t Tui) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			t.boardCursor++
			if t.boardCursor >= len(t.boards) {
				t.boardCursor = len(t.boards)
			}
		case "k":
			t.boardCursor--
			if t.boardCursor < 0 {
				t.boardCursor = 0
			}
		}
	}

	return t, nil
}

func (t Tui) View() string {
	var s string
	if !t.chosenBoard {
		s = chooseBoardView(t)
	}

	return s
}

func (t Tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			return t, tea.Quit
		}
	}
	if !t.chosenBoard {
		return boardKeys(msg, t)
	}
	return t, nil
}

func main() {
	godotenv.Load()
	jiraToken := os.Getenv("JIRA_API_TOKEN")

	client := &api.JiraClient{}
	client.Connect("killean.johnson@stairsupplies.com", jiraToken)

	p := tea.NewProgram(initialTui(client), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error has occured: %v", err)
		os.Exit(1)
	}

	// tui.CreateGUI(client)
}
