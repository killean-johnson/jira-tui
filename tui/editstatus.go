package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type EditStatus struct {
    parent *TUI
    client *api.JiraClient
}

func (pl *EditStatus) Layout(g *gocui.Gui) error {
    return nil
}
