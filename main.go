package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/killean-johnson/jira-tui/api"

	// "github.com/killean-johnson/jira-tui/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/killean-johnson/jira-tui/ui"
)

func main() {
	godotenv.Load()
	jiraToken := os.Getenv("JIRA_API_TOKEN")

	client := &api.JiraClient{}
	client.Connect("nate.cunningham@stairsupplies.com", jiraToken)

	p := tea.NewProgram(ui.InitialTui(client), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Error has occured: %v", err)
		os.Exit(1)
	}

	// tui.CreateGUI(client)
}
