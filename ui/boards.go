package ui

import (
	"fmt"
	"io"

	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)
var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item string

type itemDelegate struct{}

func (i item) FilterValue() string                               { return "" }
func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type BoardModel struct {
	common        *commonModel
	boards        []jira.Board
	selectedBoard *jira.Board
	boardsList    list.Model
}

func newBoardModel(common *commonModel) BoardModel {
	boards, err := common.jc.GetBoardList()
	if err != nil {
		fmt.Println(err)
	}
	var boardsItems []list.Item

	for _, board := range boards {
		boardsItems = append(boardsItems, item(board.Name))
	}

	const defaultWidth = 20
	const listHeight = 14

	boardsList := list.NewModel(boardsItems, itemDelegate{}, common.width, listHeight)

	boardsList.Title = "Select your board"
	boardsList.SetShowStatusBar(false)
	boardsList.SetShowPagination(false)
	boardsList.SetFilteringEnabled(false)
	boardsList.Styles.Title = titleStyle
	boardsList.Styles.HelpStyle = helpStyle

	return BoardModel{
		common:     common,
		boards:     boards,
		boardsList: boardsList,
	}
}

func (m BoardModel) View() string {
	return "\n" + m.boardsList.View()
}

func updateState(m BoardModel) tea.Cmd {
	return func() tea.Msg {
		if m.selectedBoard != nil {
			return updateModelState(stateShowIssues)
		}
		return nil
	}

}
func (m BoardModel) Update(msg tea.Msg) (BoardModel, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.boardsList.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.selectedBoard = &(m.boards[m.boardsList.Index()])
			cmds = append(cmds, updateState(m))
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.boardsList, cmd = m.boardsList.Update(msg)
	return m, cmd
}
