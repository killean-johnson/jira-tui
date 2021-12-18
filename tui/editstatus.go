package tui

import (
	"github.com/jroimartin/gocui"
)

type EditStatus struct {
    Widget
    isActive bool
}

func (pl *EditStatus) Layout(g *gocui.Gui) error {
    return nil
}
