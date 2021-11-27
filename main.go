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

    // Run the TUI
    t := new(tui.TUI)
    t.SetupTUI(client, conf)
    err = t.Run()
    if err != nil {
        panic(err)
    }

    //tui.CreateTUI(client, conf)
}
