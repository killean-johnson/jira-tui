package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type EditAssignee struct {
    parent *TUI
    client *api.JiraClient
}

func (pl *EditAssignee) Layout(g *gocui.Gui) error {
    return nil
}
