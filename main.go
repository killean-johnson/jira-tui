package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
	"github.com/killean-johnson/jira-tui/jirautils"
)

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func main() {
	// Set up env
	godotenv.Load()
	jiraToken := os.Getenv("JIRA_API_TOKEN")
	VVID := os.Getenv("JIRA_VV_TABLE_ID")

	client := jirautils.CreateClient("killean.johnson@stairsupplies.com", jiraToken)

	// Set up client
	authTransport := jira.BasicAuthTransport{
		Username: "killean.johnson@stairsupplies.com",
		Password: jiraToken,
	}

	jiraClient, err := jira.NewClient(authTransport.Client(), "https://stairsupplies-voe.atlassian.net")

	if err != nil {
		fmt.Println("Bork")
		os.Exit(0)
	}

	// List out boards
	//threeBoards, _, _ := //jiraClient.Board.GetAllSprints(
	// projects, _, _ := jiraClient.Project.GetList()
	// s, _ :=json.MarshalIndent(projects, "", "\t")
	// fmt.Printf("projects: %v\n", string(s))

	boardOpt := &jira.BoardListOptions{
		ProjectKeyOrID: VVID,
	}
	board, _, _ := .Board.GetAllBoards(boardOpt)
	var boardId string = fmt.Sprint(board.Values[0].ID)

	sprints, _, _ := jiraClient.Board.GetAllSprints(boardId)

	for i := 0; i < len(sprints); i++ {
		spr := sprints[i]
		issues, _, _ := jiraClient.Sprint.GetIssuesForSprint(spr.ID)
		MarshalPrint(issues)
	}

	// Search for issues
	// opt := &jira.SearchOptions{
	// MaxResults: 10,
	// Expand:     "fields",
	// }
	// issue, _, er := jiraClient.Issue.Search("", opt)

	// if er != nil {
	//     fmt.Printf("Failed! %+v", er)
	// } else {
	//     fmt.Printf("\nIssues: %+v\nSuccess!\n", issue)
	// }
}

// import (
// 	"fmt"
// 	"log"

// 	"github.com/jroimartin/gocui"
// )

// func main() {
// 	g, err := gocui.NewGui(gocui.OutputNormal)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	defer g.Close()

// 	g.SetManagerFunc(layout)

// 	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
// 		log.Panicln(err)
// 	}

// 	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
// 		log.Panicln(err)
// 	}
// }

// func layout(g *gocui.Gui) error {
// 	maxX, maxY := g.Size()
// 	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
// 		if err != gocui.ErrUnknownView {
// 			return err
// 		}
// 		fmt.Fprintln(v, "Hello World")
// 	}
// 	return nil
// }

// func quit(g *gocui.Gui, v *gocui.View) error {
// 	return gocui.ErrQuit
// }
