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

	if issue.Fields.Assignee != nil {
		fmt.Fprintf(v, "ID: %s\nAssigned To %s\nDescription: %s\n", issue.ID, issue.Fields.Assignee.DisplayName, issue.Fields.Description)
	} else {
		fmt.Fprintf(v, "ID: %s\nUnassigned\nDescription: %s\n", issue.ID, issue.Fields.Description)
	}

	if _, err := g.SetCurrentView("issuelist"); err != nil {
		return err
	}

	return nil
}

func (il *IssueLayout) editDescription(g *gocui.Gui, v *gocui.View) error {
	if il.activeIssue != nil {
		maxX, maxY := g.Size()
		if v, err := g.SetView("editdesc", maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Editable = true
			v.Title = "Edit Description"
			v.Wrap = true
			g.Cursor = true
			fmt.Fprint(v, il.activeIssue.Fields.Description)
			if _, err := g.SetCurrentView("editdesc"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (il *IssueLayout) changeDescription(g *gocui.Gui, v *gocui.View) error {
	if il.activeIssue != nil {
		err := il.client.UpdateIssue(&jira.Issue{
			Key: il.activeIssue.Key,
			Fields: &jira.IssueFields{
				Description: v.Buffer(),
			},
		})
		if err != nil {
			return err
		}
	}
	if err := g.DeleteView("editdesc"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("issuelist"); err != nil {
		return err
	}
	g.Cursor = false
	return nil
}

func (il *IssueLayout) exitEditDescription(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("editdesc"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("issuelist"); err != nil {
		return err
	}
	g.Cursor = false
	return nil
}

func (il *IssueLayout) editStatus(g *gocui.Gui, v *gocui.View) error {
    if il.activeIssue != nil {
        statuses, err := il.client.GetStatusList()
        if err != nil {
            return err
        }
		maxX, maxY := g.Size()
		if v, err := g.SetView("editstatus", maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
            v.Highlight = true
            v.SelBgColor = gocui.ColorGreen
            v.SelFgColor = gocui.ColorBlack
			v.Title = "Set Status"

            for _, status := range(statuses) {
                fmt.Fprintf(v, "%s\n", status.Name)
            }

			if _, err := g.SetCurrentView("editstatus"); err != nil {
				return err
			}
		}
    }

    return nil
}

func (il *IssueLayout) changeStatus(g *gocui.Gui, v *gocui.View) error {
    if il.activeIssue != nil {
        var err error
        statuses, err := il.client.GetStatusList()
        if err != nil {
            return err
        }
        
        var l string
        _, cy := v.Cursor()
        if l, err = v.Line(cy); err != nil {
            l = ""
        }

        var newStatusId string
        for _, status := range(statuses) {
            if status.Name == l {
                newStatusId = fmt.Sprint(status.Name)
                break
            }
        }

		if err != nil {
			return err
		}

        err = il.client.DoTransition(il.activeIssue.Key, newStatusId)
        if err != nil {
            return err
        }
    }
	if err := g.DeleteView("editstatus"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("issuelist"); err != nil {
		return err
	}
    return nil
}

func (il *IssueLayout) exitEditStatus(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("editstatus"); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("issuelist"); err != nil {
		return err
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
	if err := il.gui.SetKeybinding("issuelist", gocui.KeyCtrlD, gocui.ModNone, il.editDescription); err != nil {
		return err
	}
    if err := il.gui.SetKeybinding("issuelist", gocui.KeyCtrlS, gocui.ModNone, il.editStatus); err != nil {
        return err
    }
	if err := il.gui.SetKeybinding("issueview", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issueview", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editdesc", gocui.KeyCtrlS, gocui.ModNone, il.changeDescription); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editdesc", gocui.KeyCtrlX, gocui.ModNone, il.exitEditDescription); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", gocui.KeyEnter, gocui.ModNone, il.changeStatus); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", gocui.KeyCtrlX, gocui.ModNone, il.exitEditStatus); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, issueQuit); err != nil {
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

	// Get the issues
	issues, err := il.client.GetIssuesForSprint(il.sprintId)
	if err != nil {
		return err
	}

    // Set up the issue info column
    if v, err := g.SetView("issueinfo", maxX/2 - 9, 0, maxX/3*2-1, maxY-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "Issue Info"

		for _, is := range issues {
			if is.Fields.Assignee != nil {
				fmt.Fprintf(v, "%s | %s\n", is.Fields.Assignee.DisplayName, is.Fields.Status.StatusCategory.Name)
			} else {
				fmt.Fprintf(v, "%s | %s\n", "Unassigned", is.Fields.Status.StatusCategory.Name)
			}
		}
    }

    // Set up the issue list
	if v, err := g.SetView("issuelist", 0, 0, maxX/2 - 10, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// View settings
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Issue List"

		for _, is := range issues {
			if is.Fields.Assignee != nil {
				fmt.Fprintf(v, "%s | %s\n", is.Key, is.Fields.Summary)
			} else {
				fmt.Fprintf(v, "%s | %s\n", is.Key, is.Fields.Summary)
			}
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
