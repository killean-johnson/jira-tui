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
)

type BoardLayout struct {
	gui    *gocui.Gui
	client *api.JiraClient
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

	bl.gui.SetManagerFunc(il.issueLayout)

    // Start issue update goroutine
    go func(il *IssueLayout, g *gocui.Gui) {
        for {
            // Sleep for some time
            time.Sleep(time.Second * 15)
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
	if err := bl.gui.SetKeybinding("boardlist", gocui.KeyEnter, gocui.ModNone, bl.switchToIssueLayout); err != nil {
		return err
	}
	if err := bl.gui.SetKeybinding("boardlist", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := bl.gui.SetKeybinding("boardlist", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := bl.gui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, boardQuit); err != nil {
		return err
	}
	return nil
}

func (bl *BoardLayout) boardLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("boardlist", 0, 0, maxX-1, maxY-1); err != nil {
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
	}
	return nil
}

func boardQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
