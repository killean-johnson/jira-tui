package tui

import "github.com/jroimartin/gocui"

type Helpbar struct {
    parent *TUI
}

func (pl *Helpbar) Layout(g *gocui.Gui) error {
    return nil
}
