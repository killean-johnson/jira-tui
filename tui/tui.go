package tui

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
)

type TUI struct {
	// Application structures
	gui    *gocui.Gui
	client *api.JiraClient
	config *config.Config

	// Ties together the keymapping for every view
	keymap map[string]func(*gocui.Gui, *gocui.View) error

	// Current state information
	activeProjectKey string
	activeBoardId    int
	activeSprintId   int

	// UI
	pl *ProjectList
	il *IssueList
	iv *IssueView
	ic *IssueCreate
	ed *EditDesc
	es *EditStatus
	ea *EditAssignee
	hb *Helpbar
	mb *MessageBox
}

func (t *TUI) SetupTUI(client *api.JiraClient, conf *config.Config) error {
	// Jira api and user configuration
	t.client = client
	t.config = conf
	t.keymap = make(map[string]func(*gocui.Gui, *gocui.View) error)

	// These defaults keep things from getting messy later on when we need to check if they're initialized
	t.activeProjectKey = ""
	t.activeBoardId = -1
	t.activeSprintId = -1

	// Create the gui and set some settings for it
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}

	t.gui = gui
	t.gui.InputEsc = true

	// Set up the layouts
    parentWidget := Widget{parent: t, client: t.client}
    t.pl = &ProjectList{Widget: parentWidget}
    t.il = &IssueList{Widget: parentWidget}
    t.iv = &IssueView{Widget: parentWidget}
    t.ic = &IssueCreate{Widget: parentWidget}
    t.ed = &EditDesc{Widget: parentWidget}
    t.es = &EditStatus{Widget: parentWidget}
    t.ea = &EditAssignee{Widget: parentWidget}
    t.hb = &Helpbar{Widget: parentWidget}
    t.mb = &MessageBox{}

	// Setting up keymaps has to happen AFTER the layouts are created
	t.SetupProjectLayoutKeymap()
	t.SetupIssueViewLayoutKeymap()

	// TODO: Replace this with actual caching for the project id
	alreadyHasProjectInCache := false

	// Pull up the project list if they don't have a cached project id
	if !alreadyHasProjectInCache {
		err = t.SetupProjectLayout()
		if err != nil {
			return err
		}
	} else {
		err = t.SetupIssueViewLayout()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *TUI) Run() error {
	defer t.gui.Close()
	if err := t.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return nil
}

// Set the keymap for the project layout
func (t *TUI) SetupProjectLayoutKeymap() error {
	t.keymap[PLCURSORDOWN] = cursorDown
	t.keymap[PLCURSORUP] = cursorUp
	t.keymap[PLSELECT] = t.pl.SelectProject
	t.keymap[PLQUIT] = t.Quit

	t.keymap[MBCLEAR] = t.mb.ExitMessageBox
	return nil
}

// Set the actual keybindings for the project layout
func (t *TUI) ProjectLayoutKeybind() error {
	err := t.SetKeybinding(t.config.Project)
	if err != nil {
		return err
	}
	return nil
}

func (t *TUI) SetupProjectLayout() error {
	// Set up the starting manager and the keymap for it
	t.gui.SetManager(t, t.pl, t.hb, t.mb)
	if err := t.ProjectLayoutKeybind(); err != nil {
		return err
	}
	return nil
}

// Set the keymap for the issue view layout
func (t *TUI) SetupIssueViewLayoutKeymap() error {
	t.keymap[ILCURSORDOWN] = cursorDown
	t.keymap[ILCURSORUP] = cursorUp
	t.keymap[ILSELECTISSUE] = t.il.SelectIssue
	t.keymap[ILEDITDESCRIPTION] = t.ed.Dialogue
	t.keymap[ILEDITSTATUS] = t.es.Dialogue
	t.keymap[ILEDITASSIGNEE] = t.ea.Dialogue
	t.keymap[ILADDISSUE] = t.ic.Dialogue
	t.keymap[ILQUIT] = t.Quit

	t.keymap[IVCURSORDOWN] = cursorDown
	t.keymap[IVCURSORUP] = cursorUp

	t.keymap[EDCONFIRM] = t.ed.Confirm
	t.keymap[EDCANCEL] = t.ed.Cancel

	t.keymap[ESCURSORDOWN] = cursorDown
	t.keymap[ESCURSORUP] = cursorUp
	t.keymap[ESCONFIRM] = t.es.Confirm
	t.keymap[ESCANCEL] = t.es.Cancel

	t.keymap[EACURSORDOWN] = cursorDown
	t.keymap[EACURSORUP] = cursorUp
	t.keymap[EACONFIRM] = t.ea.Confirm
	t.keymap[EACANCEL] = t.ea.Cancel

	t.keymap[CISCYCLE] = t.ic.Cycle
	t.keymap[CISCONFIRM] = t.ic.Confirm
	t.keymap[CISCANCEL] = t.ic.Cancel

	t.keymap[CIACURSORDOWN] = cursorDown
	t.keymap[CIACURSORUP] = cursorUp
	t.keymap[CIASETASSIGNEE] = t.ic.SetAssignee
	t.keymap[CIACYCLE] = t.ic.Cycle
    t.keymap[CIACONFIRM] = t.ic.Confirm
	t.keymap[CIACANCEL] = t.ic.Cancel

	t.keymap[CIDCYCLE] = t.ic.Cycle
    t.keymap[CIDCONFIRM] = t.ic.Confirm
	t.keymap[CIDCANCEL] = t.ic.Cancel

	t.keymap[MBCLEAR] = t.mb.ExitMessageBox
	return nil
}

// Set the actual keybindings for the issue view layout
func (t *TUI) IssueViewLayoutKeybind() error {
	err := t.SetKeybinding(t.config.Issue)
	if err != nil {
		return err
	}
	return nil
}

func (t *TUI) SetupIssueViewLayout() error {
	// Set up the starting manager and the keymap for it
	t.gui.SetManager(t, t.il, t.iv, t.ic, t.ed, t.ea, t.es, t.hb)
	if err := t.IssueViewLayoutKeybind(); err != nil {
		return err
	}
	return nil
}

func (t *TUI) SetKeybinding(views []config.LayoutStruct) error {
	for _, view := range views {
		for _, key := range view.Keys {
			if len(key.Key) > 1 {
				var keySet gocui.Key

				if strings.Contains(key.Key, "<C-") {
					char := key.Key[3]
					var val int = int(char) - 96
					keySet = gocui.Key(val)
				} else {
					switch key.Key {
					case "<ENTER>":
						keySet = gocui.KeyEnter
					case "<ESCAPE>":
						keySet = gocui.KeyEsc
					case "<TAB>":
						keySet = gocui.KeyTab
					}
				}

				if err := t.gui.SetKeybinding(view.View, keySet, gocui.ModNone, t.keymap[key.Name]); err != nil {
					return err
				}
			} else {
				if err := t.gui.SetKeybinding(view.View, rune(key.Key[0]), gocui.ModNone, t.keymap[key.Name]); err != nil {
					return err
				}
			}
		}
	}

	if err := t.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, t.Quit); err != nil {
		return err
	}

	return nil
}

// Runs in a routine. Checks jira to see if there were any updates to the active list
func (t *TUI) Updater(g *gocui.Gui) error {
	// Make sure we've actually gone into a project and have a sprint
	if t.activeBoardId != -1 && t.activeSprintId != -1 {
		// Need to have the current view be the issue list, otherwise prooooobably don't update
		if t.gui.CurrentView() == nil || t.gui.CurrentView().Name() != ISSUELIST {
			return nil
		}

		// Get the issues
		issues, err := t.client.GetIssuesForSprint(t.activeSprintId)
		if err != nil {
			return err
		}

		// TODO: Update this at some point to do something smarter than just replacing the list. Really need to be wiser about this...
		t.il.issues = issues

		// Update the actual GUI interface
		t.gui.Update(t.il.RedrawList)
	}
	return nil
}

// Handle things here that will always be on the UI (ie. Helpbar)
func (t *TUI) Layout(g *gocui.Gui) error {
	// Update the helpbar based on the present context
	curView := g.CurrentView()
	if curView != nil {
		// Gather all of the bindings in the current view
		helpString := ""
		allBindings := append(t.config.Project, t.config.Issue...)
		for _, l := range allBindings {
			if l.View == curView.Name() {
				for _, key := range l.Keys {
					helpString += fmt.Sprintf("| %s - %s |", key.Key, key.Description)
				}
				break
			}
		}

		// Update the helpbar to display those bindings
		t.hb.Update(helpString)
	}
	return nil
}

func (t *TUI) ActiveIssue() *jira.Issue {
    return t.iv.activeIssue
}

func (t *TUI) Issues() *[]jira.Issue {
    return &t.il.issues
}

func (t *TUI) ShowMessageBox(msg string) error {
    return t.mb.ShowMessageBox(t.gui, msg)
}

func (t *TUI) Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
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
		}
	}
	return nil
}
