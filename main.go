package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/killean-johnson/jira-tui/api"
)

func MarshalPrint(obj interface{}) {
	s, _ := json.MarshalIndent(obj, "", "\t")
	fmt.Printf("%v\n", string(s))
}

func main() {
	// Set up env
	godotenv.Load()
	jiraToken := os.Getenv("JIRA_API_TOKEN")
	// VVID := os.Getenv("JIRA_VV_TABLE_ID")

	client := &api.JiraClient{}
	client.Connect("killean.johnson@stairsupplies.com", jiraToken)

	/* projectList, err := client.GetProjectList()
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(projectList) */

	/* boardsList, err := client.GetBoardList()
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(boardsList) */

	sprintList, err := client.GetSprintList(6)
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(sprintList)

	/* statusList, err := client.GetStatusList()
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(statusList) */

	/* projects, _, err := client.Project.GetList()
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(projects) */

	/* boardList, _, err := client.Board.GetAllBoards(&jira.BoardListOptions{})
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(boardList.Values) */

	/* searchFields := []string{"summary", "status"}
	issues, _, err := client.Issue.Search("sprint = ", &jira.SearchOptions{
		Fields: searchFields,
	})
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(issues) */

	/* categoryList, _, err := client.StatusCategory.GetList()
	if err != nil {
		fmt.Println(err)
	}
	MarshalPrint(categoryList) */

	// List out boards
	//threeBoards, _, _ := //jiraClient.Board.GetAllSprints(
	// projects, _, _ := jiraClient.Project.GetList()
	// s, _ :=json.MarshalIndent(projects, "", "\t")
	// fmt.Printf("projects: %v\n", string(s))

	/* boardOpt := &jira.BoardListOptions{
		ProjectKeyOrID: VVID,
	}
	board, _, _ := client.Board.GetAllBoards(boardOpt)
	var boardId string = fmt.Sprint(board.Values[0].ID)

	sprints, _, _ := jiraClient.Board.GetAllSprints(boardId)

	for i := 0; i < len(sprints); i++ {
		spr := sprints[i]
		issues, _, _ := jiraClient.Sprint.GetIssuesForSprint(spr.ID)
		MarshalPrint(issues)
	} */

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
