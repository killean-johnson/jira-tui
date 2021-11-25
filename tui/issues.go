package tui

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type IssueLayout struct {
	gui         *gocui.Gui
	client      *api.JiraClient
	activeIssue *jira.Issue
	sprintId    int
}

func (il *IssueLayout) selectIssue(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	if v, err = g.SetCurrentView("issueview"); err != nil {
		return err
	}

	v.Clear()

	id := strings.Trim(strings.Split(l, "|")[0], " ")
	issue, err := il.client.GetIssue(id)
	if err != nil {
		return err
	}

	il.activeIssue = issue

	fmt.Fprintf(v, "ID: %s\nAssigned To %s\nDescription: %s\n", issue.ID, issue.Fields.Assignee.DisplayName, issue.Fields.Description)

	if _, err := g.SetCurrentView("issuelist"); err != nil {
		return err
	}

	return nil
}

func (il *IssueLayout) updateIssue(g *gocui.Gui, v *gocui.View) error {
	if il.activeIssue != nil {
		err := il.client.UpdateIssue(&jira.Issue{
			Key: il.activeIssue.Key,
			Fields: &jira.IssueFields{
				Description: il.activeIssue.Fields.Description + "Please Work",
			},
		})
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (il *IssueLayout) issueLayoutKeybindings() error {
	if err := il.gui.SetKeybinding("issuelist", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issuelist", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issuelist", gocui.KeyEnter, gocui.ModNone, il.selectIssue); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issuelist", 'u', gocui.ModNone, il.updateIssue); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issueview", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issueview", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("", 'q', gocui.ModNone, issueQuit); err != nil {
		return err
	}
	return nil
}

func (il *IssueLayout) issueLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("issueview", maxX/3*2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
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
		issues, err := il.client.GetIssuesForSprint(il.sprintId)
		if err != nil {
			return err
		}

		for _, is := range issues {
			fmt.Fprintf(v, "%s | %s\n", is.Key, is.Fields.Summary)
		}

		if _, err := g.SetCurrentView("issuelist"); err != nil {
			return err
		}
	}

	return nil
}

func issueQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
