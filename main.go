package main

import (
	"fmt"
	"os"
)

var DEFAULT_KEYWORDS = []string{"TODO", "FIXME", "REFACTOR"}

func main() {
	parsedJiraConfig, err := parseYamlConfigFile("./test.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing the yaml config file.\n%s\n", err)
	}
	creds := NewJiraBasicAuthCreds()
	creds.username = os.Getenv("PROJECT_USERNAME")
	creds.password = os.Getenv("PROJECT_PASSWORD")
  keywordSlice := parsedJiraConfig.ReturnKeywordSlice()
  if len(keywordSlice) == 0 {
    keywordSlice = DEFAULT_KEYWORDS
  }
	jc := NewJiraClient(creds, parsedJiraConfig.Jira.Project.BaseUrl)
	wsl := Weasel{
		Keywords:      keywordSlice,
		baseRemoteUrl: parsedJiraConfig.RepoURL,
	}
	fmt.Fprintf(os.Stdout, "TODO regex: %s\n", wsl.todoRegex("TODO"))
	// TODO: Implement depth searchs reports. Visit every file in the "dirs" param
	issuesToReport := []*Issue{}
	wsl.searchTodos("test.txt", func(todo Todo) error {
		issue := jc.CreateNewIssueFromTODO(todo)
		if issue != nil {
			// TODO: Store the created issue to issuesToReport slice and report after
			issuesToReport = append(issuesToReport, issue)
		}
		return nil
	})
	for _, issue := range issuesToReport {
		createdIssueResp, err := jc.ReportIssueAsJiraTicket(issue)
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Cant report the following issue: '%s'. Skipping for now.\n",
				issue.Summary,
			)
		}
    issue.Todo.ReportedID = &createdIssueResp.Key
    err = issue.Todo.ChangeTodoStatus()
    if err != nil {
      continue
    }
    err = issue.Todo.CommitReportedTodo()
    if err != nil {
      continue
    }
	}
}

// Returns new instance of JiraClient
func NewJiraClient(creds *JiraBasicAuthCreds, bURL string) *JiraClient {
	return &JiraClient{
		creds:   creds,
		baseURL: bURL,
	}
}
