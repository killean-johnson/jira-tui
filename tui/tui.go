package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

func CreateGUI(client *api.JiraClient) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

    var sprintId = new(int)
	gui.SetManagerFunc(boardSetupLayout(client))

	if err := keybindings(gui, client, sprintId); err != nil {
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

func swapLayout(client * api.JiraClient, sprintId *int) func (*gocui.Gui, *gocui.View) error {
    return func (g *gocui.Gui, v *gocui.View) error {
        var l string
        var err error
        _, cy := v.Cursor()
        if l, err = v.Line(cy); err != nil {
            l = ""
        }

        *sprintId, _ = strconv.Atoi(strings.Split(l, "|")[0])

        g.SetManagerFunc(issueLayout(client, sprintId))

        if err := keybindings(g, client, sprintId); err != nil {
            log.Panicln(err)
        }
        return nil
    }
}

func keybindings(g *gocui.Gui, client * api.JiraClient, sprintId *int) error {
    if err := g.SetKeybinding("boardlist", gocui.KeyEnter, gocui.ModNone, swapLayout(client, sprintId)); err != nil {
        return err
    }
	if err := g.SetKeybinding("boardlist", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("boardlist", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
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

func issueLayout(client *api.JiraClient, sprintId *int) func(*gocui.Gui) error {
    return func(g *gocui.Gui) error {
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
            // View settings
            v.Highlight = true
            v.SelBgColor = gocui.ColorGreen
            v.SelFgColor = gocui.ColorBlack
            v.Title = "Issue List"

            // Issues
            issues, err := client.GetIssuesForSprint(*sprintId)
            if err != nil {
                return err
            }

            for _, is := range(issues) {
                fmt.Fprintf(v, "%s | %s\n", is.Key, is.Fields.Summary)
            }

            if _, err := g.SetCurrentView("issuelist"); err != nil {
                return err
            }
        }

        return nil
    }
}

func boardSetupLayout(client *api.JiraClient) func(*gocui.Gui) error {
    return func (g *gocui.Gui) error {
        maxX, maxY := g.Size()
        if v, err := g.SetView("boardlist", 0, 0, maxX - 1, maxY - 1); err != nil {
            if err != gocui.ErrUnknownView {
                return err
            }

            // View settings
            v.Highlight = true
            v.SelBgColor = gocui.ColorGreen
            v.SelFgColor = gocui.ColorBlack
            v.Title = "Board List"

            // Get the boards in this list
            boards, err := client.GetBoardList()
            if err != nil {
                return err
            }

            for _, board := range(boards) {
                fmt.Fprintf(v, "%d|%#v|%s\n", board.ID, board, board.Name)
            }

            if _, err := g.SetCurrentView("boardlist"); err != nil {
                return err
            }
        }
        return nil
    }
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
