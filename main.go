package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/tui"
)

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func main() {
	godotenv.Load()
	jiraToken := os.Getenv("JIRA_API_TOKEN")
    email := os.Getenv("JIRA_EMAIL")

	client := &api.JiraClient{}
	client.Connect(email, jiraToken)

    tui.CreateGUI(client)
}
