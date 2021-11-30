package tui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type IssueList struct {
	parent *TUI
	client *api.JiraClient

	issues []jira.Issue
}

func (il *IssueList) SelectIssue(g *gocui.Gui, v *gocui.View) error {
	/* // Get the string from the currently highlighted line
	_, cy := v.Cursor()
	l, err := v.Line(cy)
	if err != nil {
		l = ""
	}

	// Switch to the issue view
	if v, err = il.SetCurrentView("issueview"); err != nil {
		return err
	}

	// Clear out any previously written info
	v.Clear()

	// Match the ID from the string up with the one in our issue list
	id := strings.Trim(strings.Split(l, "|")[0], " ")
	var foundIssue *jira.Issue = il.getLocalIssueUtil(id)

	// Do nothing if the issue wasn't found
	if foundIssue == nil {
		return nil
	}

	// Update our actively selected issue
	il.activeIssue = foundIssue

	// Display the issue information in the issueinfo view
	assignee := ""
	if foundIssue.Fields.Assignee != nil {
		assignee = foundIssue.Fields.Assignee.DisplayName
	} else {
		assignee = "Unassigned"
	}
	fmt.Fprintf(v, "%s\nAssigned To %s\nStatus: %s\nSummary: %s\nDescription: %s\n",
		foundIssue.Key, assignee, foundIssue.Fields.Status.StatusCategory.Name, foundIssue.Fields.Summary, foundIssue.Fields.Description)

	// Return to the issuelist view
	if _, err := g.SetCurrentView("issuelist"); err != nil {
		return err
	} */

	return nil
}

func (il *IssueList) RedrawList(g *gocui.Gui) error {
	issueView, err := g.View(ISSUELIST)
	if err != nil {
		return err
	}

	issueView.Clear()
	maxX, _ := g.Size()

	// Set up the issue list
	issueListWidth := maxX/3*2 - 1
	issueTextWidth := issueListWidth / 3 * 2
	titleIssueTextWidth := issueListWidth/3*2 - 9
	issueInfoWidth := issueListWidth / 3
	issueView.Title = fmt.Sprintf("%-5s | %-"+fmt.Sprint(titleIssueTextWidth)+"s | %-20s | %-12s", "Key", "Issue", "Assignee", "Status")

	for _, is := range il.issues {
		var issueText, issueInfo string

		if len(is.Fields.Summary)+9 > issueTextWidth {
			issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary[:issueTextWidth-9])
		} else {
			issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary)
		}

		if is.Fields.Assignee != nil {
			issueInfo = fmt.Sprintf("%-20s | %-12s", is.Fields.Assignee.DisplayName, is.Fields.Status.StatusCategory.Name)
		} else {
			issueInfo = fmt.Sprintf("%-20s | %-12s", "Unassigned", is.Fields.Status.StatusCategory.Name)
		}

		fmt.Fprintf(issueView, "%-"+fmt.Sprint(issueTextWidth)+"s | %-"+fmt.Sprint(issueInfoWidth)+"s\n", issueText, issueInfo)
	}

	return nil
}

func (il *IssueList) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	issueListWidth := maxX/3*2 - 1
	v, err := g.SetView(ISSUELIST, 0, 0, issueListWidth, maxY-4)
	if err == gocui.ErrUnknownView {
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack

		_, err = g.SetCurrentView(ISSUELIST)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	il.RedrawList(g)

	return nil
}

//func (il *IssueList) getLocalIssueUtil()
