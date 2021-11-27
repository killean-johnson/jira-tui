package tui

import "github.com/jroimartin/gocui"

type IssueView struct {
    parent *TUI
}

func (pl *IssueView) Layout(g *gocui.Gui) error {
    return nil
}
