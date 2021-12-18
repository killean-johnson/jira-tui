package tui

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type IssueCreate struct {
	parent *TUI
	client *api.JiraClient

	isActive     bool
	activeWidget string
}

func (ic *IssueCreate) Dialogue(g *gocui.Gui, v *gocui.View) error {
	ic.isActive = true
	return nil
}

func (ic *IssueCreate) SetAssignee(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	name, err := v.Line(cy)
	if err != nil {
		return err
	}
	v.Title = name
	return nil
}

func (ic *IssueCreate) Cycle(g *gocui.Gui, v *gocui.View) error {
	switch ic.activeWidget {
	case CISUMMARY:
		g.SetCurrentView(CIASSIGNEE)
		ic.activeWidget = CIASSIGNEE
	case CIASSIGNEE:
		g.SetCurrentView(CIDESCRIPTION)
		ic.activeWidget = CIDESCRIPTION
	case CIDESCRIPTION:
		g.SetCurrentView(CISUMMARY)
		ic.activeWidget = CISUMMARY
	}
	return nil
}

func (ic *IssueCreate) Confirm(g *gocui.Gui, v *gocui.View) error {
	// Gather the info from the widgets
	v, err := g.SetCurrentView(CISUMMARY)
	if err != nil {
		return err
	}
	summary := strings.ReplaceAll(v.Buffer(), "\n", "")
	if summary == "" {
		ic.parent.mb.ShowMessageBox(g, "Summary can't be empty!")
		return nil
	}

	v, err = g.SetCurrentView(CIDESCRIPTION)
	if err != nil {
		return err
	}
	description := v.Buffer()
	if description == "" {
		ic.parent.mb.ShowMessageBox(g, "Description can't be empty!")
		return nil
	}

	v, err = g.SetCurrentView(CIASSIGNEE)
	if err != nil {
		return err
	}
	assigneeName := v.Title
	if assigneeName == "Assignee" {
		assigneeName = ""
	}

	// Get the jira user info
	board, err := ic.client.GetBoard(ic.parent.activeBoardId)
	if err != nil {
		return nil
	}

	users, err := ic.client.GetUsers(strings.Split(board.Name, " ")[0])
	if err != nil {
		return err
	}

	assignee := &jira.User{}
	for _, user := range *users {
		if user.DisplayName == assigneeName {
			assignee = &user
			break
		}
	}

	// Create the issue
	issue := jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     summary,
			Description: description,
			Type: jira.IssueType{
				Name: "Story",
			},
			Project: jira.Project{
                Key: ic.parent.activeProjectKey,
            },
			Assignee: assignee,
		},
	}

	_, err = ic.client.CreateIssue(ic.parent.activeSprintId, &issue)
	if err != nil {
		return err
	}

	// Does all of the regular view cleanup
	ic.Cancel(g, v)
	return nil
}

func (ic *IssueCreate) Cancel(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(CIBACKGROUND); err != nil {
		return err
	}
	if err := g.DeleteView(CISUMMARY); err != nil {
		return err
	}
	if err := g.DeleteView(CIASSIGNEE); err != nil {
		return err
	}
	if err := g.DeleteView(CIDESCRIPTION); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(ISSUELIST); err != nil {
		return err
	}
	g.Highlight = false
	g.Cursor = false
	ic.isActive = false
	ic.activeWidget = ""
	return nil
}

func (ic *IssueCreate) Layout(g *gocui.Gui) error {
	if ic.isActive {
		g.Highlight = true
		g.Cursor = true
		g.SelFgColor = gocui.ColorRed

		maxX, maxY := g.Size()
		backMaxX := maxX - 6 - 5

		// Background box
		if v, err := g.SetView(CIBACKGROUND, 5, 5, maxX-6, maxY-6); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Create Issue"
		}

		// Summary box
		if v, err := g.SetView(CISUMMARY, 7, 6, maxX-8, 8); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Summary"
			v.Editable = true

			// Set the current view to the summary entry box
			_, err := g.SetCurrentView(CISUMMARY)
			if err != nil {
				return err
			}
			ic.activeWidget = CISUMMARY
		}

		// Assignee box
		if v, err := g.SetView(CIASSIGNEE, 7, 9, backMaxX/3, maxY-7); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Assignee"
			v.Highlight = true
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack

			// Get the possible assignable people and print them out
			board, err := ic.client.GetBoard(ic.parent.activeBoardId)
			if err != nil {
				return err
			}

			users, err := ic.client.GetUsers(strings.Split(board.Name, " ")[0])
			if err != nil {
				return err
			}

			for _, user := range *users {
				fmt.Fprintln(v, user.DisplayName)
			}
		}

		// Description box
		if v, err := g.SetView(CIDESCRIPTION, 1+backMaxX/3, 9, maxX-8, maxY-7); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = "Description"
			v.Editable = true
		}
	}

	return nil
}
