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
    issueList   []jira.Issue
	sprintId    int
}

func (il *IssueLayout) getLocalIssueUtil(key string) *jira.Issue {
    for _, issue := range(il.issueList) {
        if issue.Key == key {
            return &issue
        }
    }
    return nil
}

func (il *IssueLayout) redrawIssueList(g *gocui.Gui, v *gocui.View) error {
    issueView, err := g.SetCurrentView("issuelist")
    if err != nil {
        return err
    }
    issueView.Clear()
	maxX, _ := g.Size()

    // Set up the issue list
    issueListWidth := maxX / 3 * 2 - 1

    issueTextWidth := fmt.Sprint(issueListWidth / 3 * 2)
    issueInfoWidth := fmt.Sprint(issueListWidth / 3)
    for _, is := range il.issueList {
        var issueText, issueInfo string
        if is.Fields.Assignee != nil {
            issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary)
            issueInfo = fmt.Sprintf("%-20s | %-12s", is.Fields.Assignee.DisplayName, is.Fields.Status.StatusCategory.Name)
        } else {
            issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary)
            issueInfo = fmt.Sprintf("%-20s | %-12s", "Unassigned", is.Fields.Status.StatusCategory.Name)
        }

        fmt.Fprintf(issueView, "%-" + issueTextWidth + "s | %-" + issueInfoWidth + "s\n", issueText, issueInfo)
    }

    if _, err := g.SetCurrentView("issuelist"); err != nil {
        return err
    }
    return nil
}

func (il *IssueLayout) selectIssue(g *gocui.Gui, v *gocui.View) error {
    // Get the string from the currently highlighted line
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
	if foundIssue.Fields.Assignee != nil {
		fmt.Fprintf(v, "ID: %s\nAssigned To %s\nDescription: %s\n", foundIssue.ID, foundIssue.Fields.Assignee.DisplayName, foundIssue.Fields.Description)
	} else {
		fmt.Fprintf(v, "ID: %s\nUnassigned\nDescription: %s\n", foundIssue.ID, foundIssue.Fields.Description)
	}

    // Return to the issuelist view
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
        buf := v.Buffer()
		err := il.client.UpdateIssue(&jira.Issue{
			Key: il.activeIssue.Key,
			Fields: &jira.IssueFields{
				Description: buf,
			},
		})
		if err != nil {
			return err
		}

        localIssue := il.getLocalIssueUtil(il.activeIssue.Key)
        if localIssue != nil {
            localIssue.Fields.Description = buf
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

        var newStatus jira.StatusCategory
        var newStatusId string
        for _, status := range(statuses) {
            if status.Name == l {
                newStatusId = fmt.Sprint(status.Name)
                newStatus = status
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

        il.activeIssue.Fields.Status.StatusCategory = newStatus

        il.redrawIssueList(g, v)
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
	if err := il.gui.SetKeybinding("issuelist", 'd', gocui.ModNone, il.editDescription); err != nil {
		return err
	}
    if err := il.gui.SetKeybinding("issuelist", 's', gocui.ModNone, il.editStatus); err != nil {
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
	if err := il.gui.SetKeybinding("editdesc", gocui.KeyEsc, gocui.ModNone, il.exitEditDescription); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", gocui.KeyEnter, gocui.ModNone, il.changeStatus); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", gocui.KeyEsc, gocui.ModNone, il.exitEditStatus); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("editstatus", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issuelist", 'q', gocui.ModNone, issueQuit); err != nil {
		return err
	}
	if err := il.gui.SetKeybinding("issueview", 'q', gocui.ModNone, issueQuit); err != nil {
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

    // Set up the issue list
    issueListWidth := maxX / 3 * 2 - 1
	if v, err := g.SetView("issuelist", 0, 0, issueListWidth, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// View settings
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = "Issue List"

        issueTextWidth := fmt.Sprint(issueListWidth / 3 * 2)
        issueInfoWidth := fmt.Sprint(issueListWidth / 3)
		for _, is := range il.issueList {
            var issueText, issueInfo string
            if is.Fields.Assignee != nil {
                issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary)
                issueInfo = fmt.Sprintf("%-20s | %-12s", is.Fields.Assignee.DisplayName, is.Fields.Status.StatusCategory.Name)
            } else {
                issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary)
                issueInfo = fmt.Sprintf("%-20s | %-12s", "Unassigned", is.Fields.Status.StatusCategory.Name)
            }

            fmt.Fprintf(v, "%-" + issueTextWidth + "s | %-" + issueInfoWidth + "s\n", issueText, issueInfo)
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
