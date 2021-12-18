package tui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
)

type EditDesc struct {
	parent *TUI
	client *api.JiraClient

	isActive bool
}

func (ed *EditDesc) Dialogue(g *gocui.Gui, v *gocui.View) error {
	if ed.parent.ActiveIssue() != nil {
		ed.isActive = true
	} else {
        ed.parent.ShowMessageBox("No issue selected!")
	}
	return nil
}

func (ed *EditDesc) Confirm(g *gocui.Gui, v *gocui.View) error {
    // Update the jira issue
    activeIssue := ed.parent.ActiveIssue()
    buf := v.Buffer()
    err := ed.client.UpdateIssue(&jira.Issue {
        Key: activeIssue.Key,
        Fields: &jira.IssueFields {
            Description: buf,
        },
    })
    if err != nil {
        return err
    }

    // Update the local issue
	for _, issue := range(*ed.parent.Issues()) {
		if issue.Key == activeIssue.Key {
			issue.Fields.Description = buf
            break
		}
	}

    // Clean up views
    return ed.Cancel(g, v)
}

func (ed *EditDesc) Cancel(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView(EDITDESC); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(ISSUELIST); err != nil {
		return err
	}
	g.Cursor = false
	ed.isActive = false
	return nil
}

func (ed *EditDesc) Layout(g *gocui.Gui) error {
	if ed.isActive {
		activeIssue := ed.parent.iv.activeIssue

		maxX, maxY := g.Size()
		if v, err := g.SetView(EDITDESC, maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Editable = true
			v.Title = "Edit Description"
			v.Wrap = true
			g.Cursor = true
			fmt.Fprint(v, activeIssue.Fields.Description)
			if _, err := g.SetCurrentView(EDITDESC); err != nil {
				return err
			}
		}
	}
	return nil
}
