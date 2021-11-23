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

	if err := keybindings(gui); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func swapViews(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        if v.Name() == "issuelist" {
            newView, err := g.SetCurrentView("issueview")
            newView.Highlight = true
            v.Highlight = false
            return err
        } else {
            newView, err := g.SetCurrentView("issuelist")
            newView.Highlight = true
            v.Highlight = false
            return err
        }
    }
    return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

        // Check to make sure the next line actually exists
        if l, err := v.Line(cy + 1); err == nil && l != "" {
            // Set the cursor to the next line down
            if err := v.SetCursor(cx, cy+1); err != nil {
                // Change the origin if we've hit the bottom
                if err := v.SetOrigin(ox, oy+1); err != nil {
                    return err
                }
            }
        } else {
            return err
        }
	}
	return nil
}

func selectIssue(g *gocui.Gui, v *gocui.View) error {
    var l string
    var err error

    _, cy := v.Cursor()
    if l, err = v.Line(cy); err != nil {
        l = ""
    }

    if v, err = g.SetCurrentView("issueview"); err != nil {
        return err
    }

    fmt.Fprintln(v, l)

    if _, err := g.SetCurrentView("issuelist"); err != nil {
        return err
    }

	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("issuelist", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("issuelist", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("issueview", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("issueview", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
    if err := g.SetKeybinding("issuelist", gocui.KeyEnter, gocui.ModNone, selectIssue); err != nil {
        return err
    }
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlSpace, gocui.ModNone, swapViews); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("issueview", maxX/3*2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = false
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
        v.Title = "Issue View"
	}

	if v, err := g.SetView("issuelist", 0, 0, maxX/3*2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
        v.Title = "Issue List"
		fmt.Fprintln(v, "Issue #1")
		fmt.Fprintln(v, "Issue #2")
		fmt.Fprintln(v, "Issue #3")
		fmt.Fprintln(v, "Issue #4")
		fmt.Fprintln(v, "Issue #5")
		fmt.Fprintln(v, "Issue #6")
		fmt.Fprintln(v, "Issue #7")

		if _, err := g.SetCurrentView("issuelist"); err != nil {
			return err
		}
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
