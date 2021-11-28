package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Board struct {
}

func chooseBoardView(t Tui) string {
	s := "Which Board"
	for i, board := range t.boards {
		cursor := " "
		if t.boardCursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("\n%s %s", cursor, board.Name)
	}
	return s
}

func boardKeys(msg tea.Msg, t Tui) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			t.boardCursor++
			if t.boardCursor >= len(t.boards) {
				t.boardCursor = len(t.boards)
			}
		case "k":
			t.boardCursor--
			if t.boardCursor < 0 {
				t.boardCursor = 0
			}
		}
	}

	return t, nil
}
