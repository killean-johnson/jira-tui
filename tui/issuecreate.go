package tui

import "github.com/jroimartin/gocui"

type IssueCreate struct {
    parent *TUI
}

func (pl *IssueCreate) Layout(g *gocui.Gui) error {
    return nil
}
