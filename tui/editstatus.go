package tui

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/logger"
)

type EditStatus struct {
    Widget

    statuses []jira.Status
    isActive bool
}

func (es *EditStatus) Dialogue(g *gocui.Gui, v *gocui.View) error {
    statuses, err := es.client.GetIssueTransitions(es.parent.ActiveIssue().ID)
    if err != nil {
        return err
    }
    for _, stat := range statuses {
        logger.InfoLogger.Printf("%#v\n", stat)
    }
    return nil
}

func (es *EditStatus) Confirm(g *gocui.Gui, v *gocui.View) error {
    return nil
}

func (es *EditStatus) Cancel(g *gocui.Gui, v *gocui.View) error {
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

            statuses, err := es.client.GetStatusList()
            if err != nil {
                return err
            }

            for _, status := range statuses {
                fmt.Fprintf(v, "%s\n", status.Name)
            }

            if _, err := g.SetCurrentView(EDITSTATUS); err != nil {
                return err
            }
        }
    }

    return nil
}
