package tui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
)

type TUI struct {
    // Application structures
    gui *gocui.Gui
    client *api.JiraClient
    conf *config.Config

    // Ties together the keymapping for every view
    keymap map[string]func(*gocui.Gui,*gocui.View) error

    // Current state information
    activeProjectKey string
    activeBoardId int
    activeSprintId int

    // UI
    pl *ProjectList
    il *IssueList
    iv *IssueView
    ic *IssueCreate
    ed *EditDesc
    es *EditStatus
    ea *EditAssignee
    hb *Helpbar
}

func (t *TUI) SetupTUI(client *api.JiraClient, conf *config.Config) error {
    // These defaults keep things from getting messy later on when we need to check if they're initialized
    t.activeProjectKey = ""
    t.activeBoardId = -1
    t.activeSprintId = -1

    gui, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        return err
    }

    t.gui = gui
    t.gui.InputEsc = true

    t.pl = new(ProjectList)
    t.il = new(IssueList)
    t.iv = new(IssueView)
    t.ic = new(IssueCreate)
    t.ed = new(EditDesc)
    t.es = new(EditStatus)
    t.ea = new(EditAssignee)
    t.hb = new(Helpbar)

    // Set up the starting manager and the keymap for it
    t.gui.SetManager(t.pl, t.hb)
    if err = t.ProjectLayoutKeymap(); err != nil {
        return err
    }
    
    return nil
}

func (t *TUI) Run() error {
	if err := t.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
        return err
	}
    return nil
}

func (tui *TUI) ProjectLayoutKeymap() error {
    if err := tui.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, tui.Quit); err != nil {
        return err
    }

    return nil
}

func (tui *TUI) IssueViewLayoutKeymap() error {
    return nil
}

func (tui *TUI) Quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}











func CreateTUI(client *api.JiraClient, conf *config.Config) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

    gui.InputEsc = true

    var bl *BoardLayout = new(BoardLayout)
    bl.client = client
    bl.gui = gui
    bl.config = conf
    bl.keymap = make(map[string]func(*gocui.Gui, *gocui.View) error)
    bl.keymap["blcursordown"] = cursorDown
    bl.keymap["blcursorup"] = cursorUp
    bl.keymap["blselect"] = bl.switchToIssueLayout
    bl.keymap["blquit"] = boardQuit

    var helpbar = make(chan string)
    bl.helpbar = &helpbar

    go func(hb chan string, gui *gocui.Gui) error {
        for {
            var helptext = <-hb
            gui.Update(func (gui *gocui.Gui) error {
                curView := gui.CurrentView()
                maxX, maxY := gui.Size()
                v, err := gui.SetView("helpbar", 0, maxY - 3, maxX - 1, maxY - 1)
                if err != nil && err != gocui.ErrUnknownView {
                    panic(err)
                } else if err == gocui.ErrUnknownView {
                    v.Title = "Keybindings"
                }
                gui.SetViewOnTop(v.Name())

                v.Clear()
                fmt.Fprintf(v, "%s", helptext)

                if curView != nil {
                    gui.SetCurrentView(curView.Name())
                }

                return nil
            })
        }
    }(helpbar, gui)

    gui.SetManagerFunc(bl.Layout)

	if err := bl.boardLayoutKeybindings(); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func updateHelpbar(helpbar chan string, gui *gocui.Gui, conf *config.Config) {
    curView := gui.CurrentView()

    if curView != nil {
        helpString := ""
        for _, l := range(conf.Board) {
            if l.View == curView.Name() {
                for _, key := range(l.Keys) {
                    helpString += fmt.Sprintf("| %s - %s |", key.Key, key.Description)
                }
                break
            }
        }

        for _, l := range(conf.Issue) {
            if l.View == curView.Name() {
                for _, key := range(l.Keys) {
                    helpString += fmt.Sprintf("| %s - %s |", key.Key, key.Description)
                }
                break
            }
        }
        helpbar <- helpString
    }
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
