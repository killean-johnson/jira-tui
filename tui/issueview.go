package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type IssueView struct {
    parent *TUI
    client *api.JiraClient
}

func (pl *IssueView) Layout(g *gocui.Gui) error {
    return nil
}
