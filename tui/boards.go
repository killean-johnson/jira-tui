package tui

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
)

type BoardLayout struct {
	gui    *gocui.Gui
	client *api.JiraClient
    config *config.Config
    keymap map[string]func(*gocui.Gui,*gocui.View) error
    helpbar *chan string
}

func (bl *BoardLayout) switchToIssueLayout(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	boardId, _ := strconv.Atoi(strings.Split(l, "|")[0])
	sprints, err := bl.client.GetSprintList(boardId)
	if err != nil {
		return err
	}

	var il *IssueLayout = new(IssueLayout)
    il.gui = bl.gui
	il.client = bl.client
	il.sprintId = sprints[0].ID
    il.boardId = boardId
    il.config = bl.config

    il.keymap = make(map[string]func(*gocui.Gui, *gocui.View) error)
    il.keymap["ilcursordown"] = cursorDown
    il.keymap["ilcursorup"] = cursorUp
    il.keymap["ilselectissue"] = il.selectIssue
    il.keymap["ileditdescription"] = il.editDescription
    il.keymap["ilchangestatus"] = il.editStatus
    il.keymap["ileditassignee"] = il.editAssignee
    il.keymap["ilquit"] = issueQuit
    il.keymap["ivcursordown"] = cursorDown
    il.keymap["ivcursorup"] = cursorUp
    il.keymap["edsavechanges"] = il.changeDescription
    il.keymap["edcancel"] = il.exitEditDescription
    il.keymap["escursordown"] = cursorDown
    il.keymap["escursorup"] = cursorUp
    il.keymap["essetstatus"] = il.changeStatus
    il.keymap["escancel"] = il.exitEditStatus
    il.keymap["eacursordown"] = cursorDown
    il.keymap["eacursorup"] = cursorUp
    il.keymap["easetassignee"] = il.changeAssignee
    il.keymap["eacancel"] = il.exitAssignee

    il.helpbar = bl.helpbar

	// Get the issues
	issues, err := il.client.GetIssuesForSprint(il.sprintId)
	if err != nil {
		return err
	}

    sorter := func (issues []jira.Issue) func(int, int) bool { 
        return func(i, j int) bool {
            if issues[i].Fields.Status.StatusCategory.Key > issues[j].Fields.Status.StatusCategory.Key {
                return true
            } else if issues[i].Fields.Status.StatusCategory.Key < issues[j].Fields.Status.StatusCategory.Key {
                return false
            }

            var iname, jname string
            if issues[i].Fields.Assignee == nil {
                iname = "Unassigned"
            } else {
                iname = issues[i].Fields.Assignee.DisplayName
            }
            if issues[j].Fields.Assignee == nil {
                jname = "Unassigned"
            } else {
                jname = issues[j].Fields.Assignee.DisplayName
            }

            return iname < jname
        }
    }

    sort.Slice(issues, sorter(issues))

    // TODO: Maybe find some way to have this not be a large copy operation
    il.issueList = issues

	bl.gui.SetManagerFunc(il.Layout)

    // Start issue update goroutine
    go func(il *IssueLayout, g *gocui.Gui) {
        for {
            // Sleep for some time
            time.Sleep(time.Second * 5)
            // Update the current issues
            issues, err := il.client.GetIssuesForSprint(il.sprintId)
            if err == nil {
                sort.Slice(issues, sorter(issues))
                il.issueList = issues
                il.redrawIssueList(g)
            }
        }
    }(il, g)

	if err := il.issueLayoutKeybindings(); err != nil {
		log.Panicln(err)
	}
	return nil
}


func (bl *BoardLayout) boardLayoutKeybindings() error {
    for _, view := range(bl.config.Board) {
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
                    }
                }

                if err := bl.gui.SetKeybinding(view.View, keySet, gocui.ModNone, bl.keymap[key.Name]); err != nil {
                    return err
                }
            } else {
                if err := bl.gui.SetKeybinding(view.View, rune(key.Key[0]), gocui.ModNone, bl.keymap[key.Name]); err != nil {
                    return err
                }
            }
        }
    }

    if err := bl.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, boardQuit); err != nil {
        return err
    }

	return nil
}

func (bl *BoardLayout) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("boardlist", 0, 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// View settings
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Board List"

		// Get the boards in this list
		boards, err := bl.client.GetBoardList()
		if err != nil {
			return err
		}

		for _, board := range boards {
			fmt.Fprintf(v, "%d|%s\n", board.ID, board.Name)
		}

		if _, err := g.SetCurrentView("boardlist"); err != nil {
			return err
		}

        updateHelpbar(*bl.helpbar, g, bl.config)
	}

	return nil
}

func boardQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
