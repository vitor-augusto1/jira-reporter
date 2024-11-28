package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/vitor-augusto1/jira-weasel/pkg/assert"
	"github.com/vitor-augusto1/jira-weasel/pkg/colors"
	"github.com/vitor-augusto1/jira-weasel/pkg/logger"
	"github.com/vitor-augusto1/jira-weasel/static"
)

type Flag[T any] struct {
	full     string
	short    string
	defaultV T
}

var (
	REPORT_FLAG_NAME Flag[bool] = Flag[bool]{full: "report", short: "r", defaultV: false}
	PURGE_FLAG_NAME  Flag[bool] = Flag[bool]{full: "purge", short: "p", defaultV: false}
	LIST_FLAG_NAME   Flag[bool] = Flag[bool]{full: "list", short: "l", defaultV: false}
	QUIET_FLAG_NAME  Flag[bool] = Flag[bool]{full: "quiet", short: "q", defaultV: false}

	DEFAULT_KEYWORDS []string = []string{"TODO", "FIXME", "REFACTOR"}
)

func init() {
	flag.Usage = static.WeaselDetails
}

func main() {
	parsedJiraConfig, err := parseYamlConfigFile("./test.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing the yaml config file.\n%s\n", err)
	}
  sliceContainsValidIssuesType(
    parsedJiraConfig.ReturnIssuesTypesSlice(),
    map[string]bool{"Bug": true, "Task": true},
    func(invalidType string) {
      fmt.Fprintf(
        os.Stderr,
        "An invalid issue type was provided in your yaml config file: '%s'\n",
        invalidType,
      )
      os.Exit(1)
    },
  )
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
	fmt.Fprintf(os.Stdout, "TODO regex: %s\n", wsl.TodoRegex("TODO"))
	// TODO: Implement depth searchs reports. Visit every file in the "dirs" param
	issuesToReport := []*Issue{}
  wsl.LoadProjectFiles()
  wsl.VisitAndReportWeaselFiles(func (todo Todo) error {
    issue := jc.CreateNewIssueFromTODO(todo, parsedJiraConfig.Keywords[todo.Keyword].IssueType)
    if issue != nil {
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

var DEFAULT_KEYWORDS = []string{"TODO", "FIXME", "REFACTOR"}

// Returns new instance of JiraClient
func NewJiraClient(creds *JiraBasicAuthCreds, bURL string) *JiraClient {
	return &JiraClient{
		creds:   creds,
		baseURL: bURL,
	}
}
