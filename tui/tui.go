package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
)

type TUI struct {
    // Application structures
    gui *gocui.Gui
    client *api.JiraClient
    config *config.Config

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
    // Jira api and user configuration
    t.client = client
    t.config = conf
    t.keymap = make(map[string]func(*gocui.Gui,*gocui.View) error)

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
    t.pl = new(ProjectList)
    t.il = new(IssueList)
    t.iv = new(IssueView)
    t.ic = new(IssueCreate)
    t.ed = new(EditDesc)
    t.es = new(EditStatus)
    t.ea = new(EditAssignee)
    t.hb = new(Helpbar)

    t.pl.parent = t
    t.il.parent = t
    t.iv.parent = t
    t.ic.parent = t
    t.ed.parent = t
    t.es.parent = t
    t.ea.parent = t
    t.hb.parent = t
    t.pl.client = t.client
    t.il.client = t.client
    t.iv.client = t.client
    t.ic.client = t.client
    t.ed.client = t.client
    t.es.client = t.client
    t.ea.client = t.client
    t.hb.client = t.client

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
    t.keymap["plcursordown"] = cursorDown
    t.keymap["plcursorup"] = cursorUp
    t.keymap["plselect"] = t.pl.SelectProject
    t.keymap["plquit"] = t.Quit
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
    t.gui.SetManager(t.pl, t.hb)
    if err := t.ProjectLayoutKeybind(); err != nil {
        return err
    }
    return nil
}

// Set the keymap for the issue view layout
func (t *TUI) SetupIssueViewLayoutKeymap() error {
    t.keymap["ilcursordown"] = cursorDown
    t.keymap["ilcursorup"] = cursorUp
    // t.keymap["ilselectissue"] = il.selectIssue
    // t.keymap["ileditdescription"] = il.editDescDialogue
    // t.keymap["ileditstatus"] = il.editStatusDialogue
    // t.keymap["ileditassignee"] = il.editAssigneeDialogue
    // t.keymap["iladdissue"] = il.createIssueDialogue
    t.keymap["ilquit"] = t.Quit
    // t.keymap["ivcursordown"] = cursorDown
    // t.keymap["ivcursorup"] = cursorUp
    // t.keymap["edconfirm"] = il.confirmEditDesc
    // t.keymap["edcancel"] = il.cancelEditDesc
    // t.keymap["escursordown"] = cursorDown
    // t.keymap["escursorup"] = cursorUp
    // t.keymap["esconfirm"] = il.confirmEditStatus
    // t.keymap["escancel"] = il.cancelEditStatus
    // t.keymap["eacursordown"] = cursorDown
    // t.keymap["eacursorup"] = cursorUp
    // t.keymap["eaconfirm"] = il.confirmEditAssignee
    // t.keymap["eacancel"] = il.cancelEditAssignee

    // t.keymap["ciscycle"] = il.cycleCreateIssueWidgets
    // t.keymap["cisconfirm"] = il.confirmCreateIssue
    // t.keymap["ciscancel"] = il.cancelCreateIssue
    // t.keymap["ciacursordown"] = cursorDown
    // t.keymap["ciacursorup"] = cursorUp
    // t.keymap["ciasetassignee"] = il.setCreateIssueAssignee
    // t.keymap["ciacycle"] = il.cycleCreateIssueWidgets
    // t.keymap["ciaconfirm"] = il.confirmCreateIssue
    // t.keymap["ciacancel"] = il.cancelCreateIssue
    // t.keymap["cidcycle"] = il.cycleCreateIssueWidgets
    // t.keymap["cidconfirm"] = il.confirmCreateIssue
    // t.keymap["cidcancel"] = il.cancelCreateIssue
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
    t.gui.SetManager(t.il, t.iv, t.hb)
    if err := t.IssueViewLayoutKeybind(); err != nil {
        return err
    }
    return nil
}

func (t *TUI) SetKeybinding(views []config.LayoutStruct) error {
    for _, view := range(views) {
        for _, key := range(view.Keys) {
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
		} else {
			return err
		}
	}
	return nil
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
        for _, l := range(conf.Project) {
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
