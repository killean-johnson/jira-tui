package tui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

func CreateGUI() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.SetManagerFunc(layout)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
    return nil
}

func cursorDown(g *gocui.Gui, v * gocui.View) error {
    return nil
}

func keybindings(g *gocui.Gui) error {
    if err := g.SetKeybinding("issuelist", "j", gocui.ModNone, cursorDown); err != nil {
        return err
    }
    if err := g.SetKeybinding("issuelist", "k", gocui.ModNone, cursorUp); err != nil {
        return err
    }
    return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("issuelist", 0, 0, maxX / 3 * 2, maxY / 3 * 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello World")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
