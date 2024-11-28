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
	// TODO: Add ui with bubbletea
	CheckIfCurrentDirectoryIsAGitRepository()
	CheckIfWeaselConfigFileExists()
	var conf *Config
	conf, err := LoadConfigs("./weasel.toml")
	assert.NotNil(conf, "conf cannot be nil", "conf", conf)
	if err != nil {
		logger.LogErrorExitingOne("There is an error with you config file. Please check it again.")
	}
	creds := NewJiraBasicAuthCreds()
	assert.NotNil(creds, "creds cannot be nil", "creds", creds)
	creds.username = conf.Credentials.Username
	creds.password = conf.Credentials.Password
	jiraClient := NewJiraClient(conf.Project.Url, creds.ReturnBasicAuthEncodedCredentials())
	assert.NotNil(jiraClient, "jiraClient cannot be nil", "jiraClient", jiraClient)
	keywords := conf.returnIssuesKeywordsSlice()
	weasel := &Weasel{
		Keywords:      keywords,
		baseRemoteUrl: conf.Remote,
	}
	weasel.LoadProjectFiles()
	var (
		flagVar    bool
		silentFlag bool
	)
	flag.BoolVarP(&flagVar, REPORT_FLAG_NAME.full, REPORT_FLAG_NAME.short, REPORT_FLAG_NAME.defaultV, static.REPORT_MESSAGE)
	flag.BoolVarP(&flagVar, PURGE_FLAG_NAME.full, PURGE_FLAG_NAME.short, PURGE_FLAG_NAME.defaultV, static.PURGE_MESSAGE)
	flag.BoolVarP(&flagVar, LIST_FLAG_NAME.full, LIST_FLAG_NAME.short, LIST_FLAG_NAME.defaultV, static.LIST_MESSAGE)
	flag.BoolVarP(&silentFlag, QUIET_FLAG_NAME.full, QUIET_FLAG_NAME.short, QUIET_FLAG_NAME.defaultV, static.SILENT_MESSAGE)
	flag.Parse()
	if flag.Lookup(REPORT_FLAG_NAME.full).Changed {
		reportCommand(weasel, jiraClient, &conf.Keywords, silentFlag, static.Banner)
	} else if flag.Lookup(PURGE_FLAG_NAME.full).Changed {
		purgeCommand(weasel, jiraClient, silentFlag, static.Banner)
	} else if flag.Lookup(LIST_FLAG_NAME.full).Changed {
		listCommand(weasel, jiraClient, silentFlag, static.Banner)
	} else {
		helperCommand()
		os.Exit(1)
	}
}

// TODO: Add interative mode to the commands

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
