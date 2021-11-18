package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/andygrunwald/go-jira"
	"github.com/joho/godotenv"
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

    // Set up client
    authTransport := jira.BasicAuthTransport {
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

    boardOpt := &jira.BoardListOptions {
        ProjectKeyOrID: VVID,

    }
    board, _, _ := jiraClient.Board.GetAllBoards(boardOpt)
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
