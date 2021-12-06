package tui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type IssueView struct {
	parent *TUI
	client *api.JiraClient

	activeIssue *jira.Issue
}

func (iv *IssueView) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	v, err := g.SetView(ISSUEVIEW, maxX/3*2, 0, maxX-1, maxY-4)
	if err == gocui.ErrUnknownView {
		v.Title = "Issue View"
	} else if err != nil {
		return err
	} else {
		if iv.activeIssue != nil {
			v.Clear()

			// Display the issue information in the issueinfo view
			assignee := ""
			if iv.activeIssue.Fields.Assignee != nil {
				assignee = iv.activeIssue.Fields.Assignee.DisplayName
			} else {
				assignee = "Unassigned"
			}
			fmt.Fprintf(v, "%s\nAssigned To %s\nStatus: %s\nSummary: %s\nDescription: %s\n",
				iv.activeIssue.Key, assignee, iv.activeIssue.Fields.Status.Name, iv.activeIssue.Fields.Summary, iv.activeIssue.Fields.Description)
		}
	}

	return nil
}
