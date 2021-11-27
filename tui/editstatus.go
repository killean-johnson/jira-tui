package tui

import "github.com/jroimartin/gocui"

type EditStatus struct {
    parent *TUI
}

func (pl *EditStatus) Layout(g *gocui.Gui) error {
    return nil
}
