package tui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type MessageBox struct {
	message      string
	isActive     bool
	previousView string
}

func (mb *MessageBox) ShowMessageBox(g *gocui.Gui, msg string) error {
	mb.isActive = true
	mb.message = msg

    mb.previousView = g.CurrentView().Name()

	maxX, maxY := g.Size()
	if v, err := g.SetView(MESSAGEBOX, maxX/2-30, maxY/2, maxX/2+30, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, mb.message)
		if _, err := g.SetCurrentView(MESSAGEBOX); err != nil {
			return err
		}
	}
	return nil
}

func (mb *MessageBox) ExitMessageBox(g *gocui.Gui, v *gocui.View) error {
	mb.isActive = false

	err := g.DeleteView(MESSAGEBOX)
	if err != nil {
		return err
	}

	_, err = g.SetCurrentView(mb.previousView)
	if err != nil {
		return err
	}

	return nil
}

func (mb *MessageBox) Layout(g *gocui.Gui) error {
	if mb.isActive {
		v, err := g.SetCurrentView(MESSAGEBOX)
		if err != nil {
			return err
		}

		v.Clear()
		fmt.Fprint(v, mb.message)
	}

	return nil
}
