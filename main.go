package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
	"github.com/killean-johnson/jira-tui/tui"
)

func main() {
    kb := new(config.Keybindings)
    err := kb.LoadKeybindings()
    if err != nil {
        panic(err)
    }
    return 
	godotenv.Load()
	jiraToken := os.Getenv("JIRA_API_TOKEN")
    email := os.Getenv("JIRA_EMAIL")

	client := &api.JiraClient{}
	client.Connect(email, jiraToken)

    tui.CreateGUI(client)
}
