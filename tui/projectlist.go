package tui

import "github.com/jroimartin/gocui"

type ProjectList struct {
    parent *TUI
}

func (pl *ProjectList) Layout(g *gocui.Gui) error {
    return nil
}
