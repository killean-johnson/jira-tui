package tui

import "github.com/jroimartin/gocui"

type IssueList struct {
    parent *TUI
}

func (pl *IssueList) Layout(g *gocui.Gui) error {
    return nil
}
