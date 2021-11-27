package tui

import "github.com/jroimartin/gocui"

type EditAssignee struct {
    parent *TUI
}

func (pl *EditAssignee) Layout(g *gocui.Gui) error {
    return nil
}
