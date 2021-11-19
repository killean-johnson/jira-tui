package jirautils

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func createClient(username string, token string) *jira.Client {
    // Set up client
    authTransport := jira.BasicAuthTransport {
        Username: username,
        Password: token,
    }

    jiraClient, err := jira.NewClient(authTransport.Client(), "https://stairsupplies-voe.atlassian.net")

    if err != nil {
        fmt.Println("ERR IN createClient FUNC")
        return nil
    }

    return jiraClient
}

