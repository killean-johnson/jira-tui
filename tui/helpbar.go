package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type Helpbar struct {
    parent *TUI
    client *api.JiraClient
}

func (pl *Helpbar) Layout(g *gocui.Gui) error {
    return nil
}
