package tui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
)

type EditStatus struct {
    Widget

    transitions []jira.Transition
    isActive bool
}

func (es *EditStatus) Dialogue(g *gocui.Gui, v *gocui.View) error {
    if es.parent.ActiveIssue() != nil {
        es.isActive = true
    } else {
        es.parent.ShowMessageBox("No issue selected!")
    }
    return nil
}

func (es *EditStatus) Confirm(g *gocui.Gui, v *gocui.View) error {
    // Get the highlighted line
    var l string
    var err error
    _, cy := v.Cursor()
    if l, err = v.Line(cy); err != nil {
        l = ""
    }

    // Get the actual transition object we want
    var transition jira.Transition
    for _, t := range es.transitions {
        if t.Name == l {
            transition = t
        }
    }

    // Move the issue to the new status
    err = es.client.DoTransition(es.parent.ActiveIssue().Key, transition.To.Name)
    if err != nil {
        return err
    }
    
    // Update the local issue
	for _, issue := range(*es.parent.Issues()) {
		if issue.Key == es.parent.ActiveIssue().Key {
			issue.Fields.Status = &transition.To
            break
		}
	}

    return es.Cancel(g, v)
}

func (es *EditStatus) Cancel(g *gocui.Gui, v *gocui.View) error {
    if err := g.DeleteView(EDITSTATUS); err != nil {
        return nil
    }
    if _, err := g.SetCurrentView(ISSUELIST); err != nil {
        return err
    }
    es.isActive = false
    return nil
}

func (es *EditStatus) Layout(g *gocui.Gui) error {
    if es.isActive {
        maxX, maxY := g.Size()
        if v, err := g.SetView(EDITSTATUS, maxX/4, maxY/6, maxX/4*3, maxY/6*5); err != nil {
            if err != gocui.ErrUnknownView {
                return err
            }

			v.Highlight = true
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack
			v.Title = "Set Status"

            transitions, err := es.client.GetIssueTransitions(es.parent.ActiveIssue().ID)
            if err != nil {
                return err
            }
            es.transitions = transitions

            for _, transition := range transitions {
                fmt.Fprintf(v, "%s\n", transition.Name)
            }

            if _, err := g.SetCurrentView(EDITSTATUS); err != nil {
                return err
            }
        }
    }

    return nil
}
