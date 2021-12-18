package tui

import (
	"github.com/killean-johnson/jira-tui/api"
)

type Widget struct {
    parent *TUI
    client *api.JiraClient
}
