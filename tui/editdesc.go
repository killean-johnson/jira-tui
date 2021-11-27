package tui

import "github.com/jroimartin/gocui"

type EditDesc struct {
    parent *TUI
}

func (pl *EditDesc) Layout(g *gocui.Gui) error {
    return nil
}
