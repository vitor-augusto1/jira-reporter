package main

import (
	"fmt"
	"os"
)

func main() {
	parsedJiraConfig, err := parseYamlConfigFile("./test.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing the yaml config file.\n%s\n", err)
	}
	creds := NewJiraBasicAuthCreds()
	creds.username = os.Getenv("PROJECT_USERNAME")
	creds.password = os.Getenv("PROJECT_PASSWORD")
	jc := NewJiraClient(creds, parsedJiraConfig.Jira.Project.BaseUrl)
	wsl := Weasel{
		Keywords:      []string{"TODO", "FIXME", "REFACTOR"},
		baseRemoteUrl: parsedJiraConfig.RepoURL,
	}
	fmt.Fprintf(os.Stdout, "TODO regex: %s\n", wsl.todoRegex("TODO"))
	// TODO: Implement depth searchs reports. Visit every file in the "dirs" param
	wsl.searchTodos("test.txt", func(todo Todo) error {
		issue := jc.CreateNewIssueFromTODO(todo)
		if issue != nil {
			err := jc.ReportIssueAsJiraTicket(issue)
			if err != nil {
				fmt.Fprintf(
          os.Stderr,
          "Cant report this following issue: '%s'. Skipping for now.\n",
          issue.Summary,
        )
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
