package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type IssueCreate struct {
    parent *TUI
    client *api.JiraClient
}

func (pl *IssueCreate) Layout(g *gocui.Gui) error {
    return nil
}
