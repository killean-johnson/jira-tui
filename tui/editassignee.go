package tui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
)

type EditAssignee struct {
	Widget

	users    *[]jira.User
	isActive bool
}

func (ea *EditAssignee) Dialogue(g *gocui.Gui, v *gocui.View) error {
	if ea.parent.ActiveIssue() != nil {
		ea.isActive = true
	} else {
		ea.parent.ShowMessageBox("No issue selected!")
	}
	return nil
}

func (ea *EditAssignee) Confirm(g *gocui.Gui, v *gocui.View) error {
	// Grab the line we've highlighted
	var l string
	var err error
	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	// Get the specific user we've selected
	assignee := &jira.User{}
	for _, user := range *ea.users {
		if user.DisplayName == l {
			assignee = &user
			break
		}
	}

	// Update the issue to the new assignee
	err = ea.client.UpdateIssue(&jira.Issue{
		Key: ea.parent.ActiveIssue().Key,
		Fields: &jira.IssueFields{
			Assignee: assignee,
		},
	})
	if err != nil {
		return err
	}

	// Update the local issue so it displays immediately
	for _, issue := range *ea.parent.Issues() {
		if issue.Key == ea.parent.ActiveIssue().Key {
			issue.Fields.Assignee = assignee
			break
		}
	}

	// Clean up the widget
	return ea.Cancel(g, v)
}

func (ea *EditAssignee) Cancel(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(EDITASSIGNEE); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(ISSUELIST); err != nil {
		return err
	}
    ea.isActive = false
	return nil
}

func (ea *EditAssignee) Layout(g *gocui.Gui) error {
	if ea.isActive {
		maxX, maxY := g.Size()
		if v, err := g.SetView(EDITASSIGNEE, maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Highlight = true
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack
			v.Title = "Set Status"

			users, err := ea.client.GetUsers(ea.parent.activeProjectKey)
			if err != nil {
				return err
			}
			ea.users = users

			for _, user := range *users {
				fmt.Fprintf(v, "%s\n", user.DisplayName)
			}

			if _, err := g.SetCurrentView(EDITASSIGNEE); err != nil {
				return err
			}
		}
	}
	return nil
}
