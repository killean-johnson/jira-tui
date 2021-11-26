package tui

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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

	var issueInfo *IssueLayout = new(IssueLayout)
    issueInfo.gui = bl.gui
	issueInfo.client = bl.client
	issueInfo.sprintId = sprints[0].ID

	// Get the issues
	issues, err := issueInfo.client.GetIssuesForSprint(issueInfo.sprintId)
	if err != nil {
		return err
	}

    // TODO: Maybe find some way to have this not be a large copy operation
    issueInfo.issueList = issues

	bl.gui.SetManagerFunc(issueInfo.issueLayout)

	if err := issueInfo.issueLayoutKeybindings(); err != nil {
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

func (bl *BoardLayout) boardLayoutfunc(g *gocui.Gui) error {
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
