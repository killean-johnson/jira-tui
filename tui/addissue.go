package tui

import (
	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
)


type AddIssueLayout struct {
    gui *gocui.Gui
    client *api.JiraClient
    sprintId int
    boardId int
    config *config.Config
    keymap map[string]func(*gocui.Gui, *gocui.View) error
    helpbar *chan string

    activeWidth int
    widgets map[int]string
}

func (ail *AddIssueLayout) SetView(viewname string, x0, y0, x1, y1 int) (*gocui.View, error) {
    v, err := ail.gui.SetView(viewname, x0, y0, x1, y1)
    updateHelpbar(*ail.helpbar, ail.gui, ail.config)
    return v, err
}

func (ail *AddIssueLayout) SetCurrentView(viewname string) (*gocui.View, error) {
    v, err := ail.gui.SetCurrentView(viewname)
    updateHelpbar(*ail.helpbar, ail.gui, ail.config)
    return v, err
}

func (ail *AddIssueLayout) Layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := ail.SetView("ailsummarybox", 0, 0, maxX - 1, maxY - 1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = ""
    }
    if v, err := ail.SetView("aildescriptionbox", 0, 0, maxX - 1, maxY - 1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = ""
    }
    return nil
}
