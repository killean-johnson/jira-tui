package main

import (
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
	"github.com/killean-johnson/jira-tui/tui"
)

func main() {
	// Load in the config
	conf := new(config.Config)
	err := conf.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Set up the client connection
	client := &api.JiraClient{}
	client.Connect(conf.Email, conf.APIToken, conf.JiraURL)
	/* issueKey := "VV-392"
	transitions, _, _ := client.GetConnection().Issue.GetTransitions(issueKey)
	for _, tran := range transitions {
		fmt.Printf("%#v\n", tran)
	}

	tid := ""
	for _, t := range transitions {
		if t.Name == "Pending / Testing" {
			tid = t.ID
			break
		}
	}

	_, err = client.GetConnection().Issue.DoTransition(issueKey, tid)
	panic(err)
	return */
	//client.Issue.GetTransitions(issueKey)
	// Run the TUI
	t := new(tui.TUI)
	t.SetupTUI(client, conf)
	err = t.Run()
	if err != nil {
		panic(err)
	}

	//tui.CreateTUI(client, conf)
}
