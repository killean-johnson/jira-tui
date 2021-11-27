package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type EditDesc struct {
    parent *TUI
    client *api.JiraClient
}

func (pl *EditDesc) Layout(g *gocui.Gui) error {
    return nil
}
