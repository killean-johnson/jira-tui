package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type EditAssignee struct {
	Widget
	isActive bool
}

func (pl *EditAssignee) Layout(g *gocui.Gui) error {
    return nil
}
