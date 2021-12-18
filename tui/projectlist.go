package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type ProjectList struct {
	parent *TUI
	client *api.JiraClient
}

func (pl *ProjectList) SelectProject(g *gocui.Gui, v *gocui.View) error {
	// Get the line we're highlighting
	var l string
	var err error
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	// From that line, get the selected project and the sprints on it
	boardId, _ := strconv.Atoi(strings.Split(l, "|")[0])
	sprints, err := pl.client.GetSprintList(boardId)
	if err != nil {
		return err
	}

	// Throw an error if there is no active sprint on this project
	if len(sprints) < 1 {
		pl.parent.mb.ShowMessageBox(g, "No active sprint on selected project!")
		return nil
	}

	// Transfer the info to the parent
    pl.parent.activeProjectKey = strings.Split(strings.Split(l, "|")[1], " ")[0]
	pl.parent.activeBoardId = boardId
	pl.parent.activeSprintId = sprints[0].ID

	// Set up for the issue view
	err = pl.parent.SetupIssueViewLayout()
	if err != nil {
		return err
	}

	// Begin the updater routine
	go func() {
		for {
			pl.parent.Updater(pl.parent.gui)
			time.Sleep(time.Second * 5)
		}
	}()

	return nil
}

func (pl *ProjectList) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(PROJECTLIST, 0, 0, maxX-1, maxY-4)

	if err == gocui.ErrUnknownView {
		// View settings
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Project List"

		// Get the boards and list them out
		boards, err := pl.client.GetBoardList()
		if err != nil {
			return err
		}

		for _, board := range boards {
			fmt.Fprintf(v, "%d|%s\n", board.ID, board.Name)
		}

		// Make sure that this view is set to be active
		if _, err := g.SetCurrentView(PROJECTLIST); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
