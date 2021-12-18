package tui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type Helpbar struct {
    Widget
    text string
}

func (hb *Helpbar) Update(newText string) error {
    hb.text = newText
    return nil
}

func (hb *Helpbar) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
    v, err := g.SetView("helpbar", 0, maxY - 3, maxX - 1, maxY - 1)

    if err == gocui.ErrUnknownView {
        // View settings
        v.Title = "Keybindings"
    } else if err != nil {
        return err
    } else {
        v.Clear()
        fmt.Fprint(v, hb.text)
    }

	return nil
}
