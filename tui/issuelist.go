package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type IssueList struct {
    parent *TUI
    client *api.JiraClient
}

func (pl *IssueList) Layout(g *gocui.Gui) error {
    return nil
}
