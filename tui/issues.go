package tui

import (
	"fmt"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
)

type IssueLayout struct {
	gui         *gocui.Gui
	client      *api.JiraClient
	activeIssue *jira.Issue
    issueList   []jira.Issue
    boardId int
	sprintId    int
    config *config.Config
    keymap map[string]func(*gocui.Gui, *gocui.View) error
    helpbar *chan string
}

func (il *IssueLayout) SetView(viewname string, x0, y0, x1, y1 int) (*gocui.View, error) {
    v, err := il.gui.SetView(viewname, x0, y0, x1, y1)
    updateHelpbar(*il.helpbar, il.gui, il.config)
    return v, err
}

func (il *IssueLayout) SetCurrentView(viewname string) (*gocui.View, error) {
    v, err := il.gui.SetCurrentView(viewname)
    updateHelpbar(*il.helpbar, il.gui, il.config)
    return v, err
}

func (il *IssueLayout) getLocalIssueUtil(key string) *jira.Issue {
    for _, issue := range(il.issueList) {
        if issue.Key == key {
            return &issue
        }
    }
    return nil
}

func (il *IssueLayout) redrawIssueList(g *gocui.Gui) error {
    curView := g.CurrentView()
    if curView.Name() != "issuelist" {
        return nil
    }

    issueView, err := il.SetCurrentView("issuelist")
    if err != nil {
        return err
    }

    issueView.Clear()
	maxX, _ := g.Size()

    // Set up the issue list
    issueListWidth := maxX / 3 * 2 - 1
    issueTextWidth := issueListWidth / 3 * 2
    titleIssueTextWidth := issueListWidth / 3 * 2 - 9
    issueInfoWidth := issueListWidth / 3
    issueView.Title = fmt.Sprintf("%-5s | %-" + fmt.Sprint(titleIssueTextWidth) + "s | %-20s | %-12s", "Key", "Issue", "Assignee", "Status")

    for _, is := range il.issueList {
        var issueText, issueInfo string

        if len(is.Fields.Summary) + 9 > issueTextWidth {
            issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary[:issueTextWidth - 9])
        } else {
            issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary)
        }

        if is.Fields.Assignee != nil {
            issueInfo = fmt.Sprintf("%-20s | %-12s", is.Fields.Assignee.DisplayName, is.Fields.Status.StatusCategory.Name)
        } else {
            issueInfo = fmt.Sprintf("%-20s | %-12s", "Unassigned", is.Fields.Status.StatusCategory.Name)
        }

        fmt.Fprintf(issueView, "%-" + fmt.Sprint(issueTextWidth) + "s | %-" + fmt.Sprint(issueInfoWidth) + "s\n", issueText, issueInfo)
    }

    if _, err := il.SetCurrentView("issuelist"); err != nil {
        return err
    }
    return nil
}

func (il *IssueLayout) createIssue(g *gocui.Gui, v *gocui.View) error {
    i := jira.Issue {
        Fields: &jira.IssueFields {
            Description: "Test Issue",
            Type: jira.IssueType {
                Name: "Bug",
            },
            Project: il.activeIssue.Fields.Project,
            Summary: "Just A Test Issue",
        },
    }
    
    _, err := il.client.CreateIssue(&i)
    if err != nil {
        panic(err)
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

	if v, err = il.SetCurrentView("issueview"); err != nil {
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
    assignee := ""
	if foundIssue.Fields.Assignee != nil {
        assignee = foundIssue.Fields.Assignee.DisplayName
	} else {
        assignee = "Unassigned"
	}
    fmt.Fprintf(v, "%s\nAssigned To %s\nStatus: %s\nSummary: %s\nDescription: %s\n", 
        foundIssue.Key, assignee, foundIssue.Fields.Status.StatusCategory.Name, foundIssue.Fields.Summary, foundIssue.Fields.Description)

    // Return to the issuelist view
	if _, err := il.SetCurrentView("issuelist"); err != nil {
		return err
	}

	return nil
}

func (il *IssueLayout) editDescDialogue(g *gocui.Gui, v *gocui.View) error {
	if il.activeIssue != nil {
		maxX, maxY := g.Size()
		if v, err := il.SetView("editdesc", maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Editable = true
			v.Title = "Edit Description"
			v.Wrap = true
			g.Cursor = true
			fmt.Fprint(v, il.activeIssue.Fields.Description)
			if _, err := il.SetCurrentView("editdesc"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (il *IssueLayout) confirmEditDesc(g *gocui.Gui, v *gocui.View) error {
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
	if _, err := il.SetCurrentView("issuelist"); err != nil {
		return err
	}
	g.Cursor = false
	return nil
}

func (il *IssueLayout) cancelEditDesc(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("editdesc"); err != nil {
		return err
	}
	if _, err := il.SetCurrentView("issuelist"); err != nil {
		return err
	}
	g.Cursor = false
	return nil
}

func (il *IssueLayout) editStatusDialogue(g *gocui.Gui, v *gocui.View) error {
    if il.activeIssue != nil {
        statuses, err := il.client.GetStatusList()
        if err != nil {
            return err
        }
		maxX, maxY := g.Size()
		if v, err := il.SetView("editstatus", maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
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

			if _, err := il.SetCurrentView("editstatus"); err != nil {
				return err
			}
		}
    }

    return nil
}

func (il *IssueLayout) confirmEditStatus(g *gocui.Gui, v *gocui.View) error {
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

        il.redrawIssueList(g)
    }
	if err := g.DeleteView("editstatus"); err != nil {
		return err
	}
	if _, err := il.SetCurrentView("issuelist"); err != nil {
		return err
	}
    return nil
}

func (il *IssueLayout) cancelEditStatus(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("editstatus"); err != nil {
		return err
	}
	if _, err := il.SetCurrentView("issuelist"); err != nil {
		return err
	}
    return nil
}

func (il *IssueLayout) editAssigneeDialogue(g *gocui.Gui, v *gocui.View) error {
    if il.activeIssue != nil {
        board, err := il.client.GetBoard(il.boardId)
        if err != nil {
            return err
        }

        users, err := il.client.GetUsers(strings.Split(board.Name, " ")[0])
        if err != nil {
            return err
        }

		maxX, maxY := g.Size()
		if v, err := il.SetView("editassignee", maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
            v.Highlight = true
            v.SelBgColor = gocui.ColorGreen
            v.SelFgColor = gocui.ColorBlack
			v.Title = "Set Status"

            for _, user := range(*users) {
                fmt.Fprintf(v, "%s\n", user.DisplayName)
            }

			if _, err := il.SetCurrentView("editassignee"); err != nil {
				return err
			}
		}
    }

    return nil
}

func (il *IssueLayout) confirmEditAssignee(g *gocui.Gui, v *gocui.View) error {
	if il.activeIssue != nil {
        var l string
        var err error
        _, cy := v.Cursor()
        if l, err = v.Line(cy); err != nil {
            l = ""
        }

        board, err := il.client.GetBoard(il.boardId)
        if err != nil {
            return err
        }

        users, err := il.client.GetUsers(strings.Split(board.Name, " ")[0])
        if err != nil {
            return err
        }

        assignee := &jira.User{}
        for _, user := range(*users) {
            if user.DisplayName == l {
                assignee = &user
                break
            }
        }

		err = il.client.UpdateIssue(&jira.Issue{
			Key: il.activeIssue.Key,
			Fields: &jira.IssueFields{
				Assignee: assignee,
			},
		})
		if err != nil {
			return err
		}

        localIssue := il.getLocalIssueUtil(il.activeIssue.Key)
        if localIssue != nil {
            localIssue.Fields.Assignee = assignee
        }
	}
	if err := g.DeleteView("editassignee"); err != nil {
		return err
	}
	if _, err := il.SetCurrentView("issuelist"); err != nil {
		return err
	}
	return nil
}

func (il *IssueLayout) cancelEditAssignee(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("editassignee"); err != nil {
		return err
	}
	if _, err := il.SetCurrentView("issuelist"); err != nil {
		return err
	}
    return nil
}

func (il *IssueLayout) issueLayoutKeybindings() error {
    for _, view := range(il.config.Issue) {
        for _, key := range(view.Keys) {
            if len(key.Key) > 1 {
                var keySet gocui.Key

                if strings.Contains(key.Key, "<C-") {
                    char := key.Key[3]
                    var val int = int(char) - 96
                    keySet = gocui.Key(val)
                } else {
                    switch key.Key { 
                    case "<ENTER>":
                        keySet = gocui.KeyEnter
                    case "<ESCAPE>":
                        keySet = gocui.KeyEsc
                    }
                }

                if err := il.gui.SetKeybinding(view.View, keySet, gocui.ModNone, il.keymap[key.Name]); err != nil {
                    return err
                }
            } else {
                if err := il.gui.SetKeybinding(view.View, rune(key.Key[0]), gocui.ModNone, il.keymap[key.Name]); err != nil {
                    return err
                }
            }
        }
    }

    if err := il.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, issueQuit); err != nil {
        return err
    }

	return nil
}

func (il *IssueLayout) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := il.SetView("issueview", maxX/3*2, 0, maxX-1, maxY-4); err != nil {
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
    issueTextWidth := issueListWidth / 3 * 2
    titleIssueTextWidth := issueListWidth / 3 * 2 - 9
    issueInfoWidth := issueListWidth / 3
	if v, err := il.SetView("issuelist", 0, 0, issueListWidth, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// View settings
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Title = fmt.Sprintf("%-5s | %-" + fmt.Sprint(titleIssueTextWidth) + "s | %-20s | %-12s", "Key", "Issue", "Assignee", "Status")

		for _, is := range il.issueList {
            var issueText, issueInfo string

            if len(is.Fields.Summary) + 9 > issueTextWidth {
                issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary[:issueTextWidth - 9])
            } else {
                issueText = fmt.Sprintf("%s | %s", is.Key, is.Fields.Summary)
            }

            if is.Fields.Assignee != nil {
                issueInfo = fmt.Sprintf("%-20s | %-12s", is.Fields.Assignee.DisplayName, is.Fields.Status.StatusCategory.Name)
            } else {
                issueInfo = fmt.Sprintf("%-20s | %-12s", "Unassigned", is.Fields.Status.StatusCategory.Name)
            }

            fmt.Fprintf(v, "%-" + fmt.Sprint(issueTextWidth) + "s | %-" + fmt.Sprint(issueInfoWidth) + "s\n", issueText, issueInfo)
		}

		if _, err := il.SetCurrentView("issuelist"); err != nil {
			return err
		}
	}

	return nil
}

func issueQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
