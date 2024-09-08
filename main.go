package main

import (
	"fmt"
	"os"
)

var (
	PROJECT_BASE_URL = os.Getenv("PROJECT_BASE_URL")
)

func main() {
	creds := NewJiraBasicAuthCreds()
	creds.username = os.Getenv("PROJECT_USERNAME")
	creds.password = os.Getenv("PROJECT_PASSWORD")
	jc := NewJiraClient(creds, PROJECT_BASE_URL)
  wsl := Weasel{
    Keywords: []string{"TODO","FIXME","REFACTOR"},
  }
  fmt.Println(wsl.todoRegex("TODO"))
  wsl.searchTodos("./test.txt", func(todo Todo) error {
    issue := jc.CreateNewIssueFromTODO(todo)
    if issue != nil {
      fmt.Printf("Issue to be reported: '%s' \n", issue.Summary)
      err := jc.ReportIssueAsJiraTicket(issue)
      if err != nil {
        fmt.Printf("Cant report this following issue: '%s'. Skipping for now.\n", issue.Summary)
      }
    }
    return nil
  })
}

// Returns new instance of JiraClient
func NewJiraClient(creds *JiraBasicAuthCreds, bURL string) *JiraClient {
	return &JiraClient{
		creds:   creds,
		baseURL: bURL,
	}
}
