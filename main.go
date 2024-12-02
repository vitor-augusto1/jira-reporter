package main

import (
	"fmt"
	"os"
	"sort"

	flag "github.com/spf13/pflag"
	"github.com/vitor-augusto1/jira-weasel/logger"
	"github.com/vitor-augusto1/jira-weasel/pkg/assert"
	"github.com/vitor-augusto1/jira-weasel/pkg/colors"
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
		listCommand(weasel, silentFlag, static.Banner)
	} else {
		helperCommand()
		os.Exit(1)
	}
}

func reportCommand(weasel *Weasel, jiraClient *JiraClient, keywordMap *map[string]string, quiet bool, bannerFunc func()) {
	assert.NotNil(weasel, "weasel can't be nil", "weasel", weasel)
	assert.NotNil(jiraClient, "jiraClient can't be nil", "jiraClient", jiraClient)
	assert.NotNil(keywordMap, "keywordMap can't be nil", "keywordMap", keywordMap)
	if !quiet {
		bannerFunc()
	}
	issuesToReport := []*Issue{}
	weasel.VisitTodosInWeaselFiles(func(t Todo) error {
		if t.ReportedID != nil {
			return nil
		}
		assert.Nil(t.ReportedID, "Trying to report an already reported todo.", "t.ReportedID", t.ReportedID)
		mappedKeyword, ok := (*keywordMap)[t.Keyword]
		assert.Assert(ok, "The provided todo keyword was not mapped", "t.Keyword", t.Keyword)
		issue := jiraClient.CreateNewIssueFromTODO(t, mappedKeyword)
		if issue != nil {
			issuesToReport = append(issuesToReport, issue)
		}
		return nil
	})
	for _, issue := range issuesToReport {
		if !quiet {
			fmt.Fprintf(
				os.Stdout,
				" 📢 Reporting a new '%s' [%s %s]\n\n",
				colors.Info(issue.Todo.Keyword),
				colors.Success(issue.Todo.FilePath),
				colors.Success(fmt.Sprint(issue.Todo.Line)),
			)
		}
		createdIssueResp, err := jiraClient.ReportIssueAsJiraTicket(issue)
		if err != nil {
			logger.LogErrorExitingOne(fmt.Sprintf("Can't report the following issue: '%s'. Skipping for now.\n", issue.Summary))
		}
		issue.Todo.ReportedID = &createdIssueResp.Key
		err = issue.Todo.ChangeTodoStatusToReported()
		if err != nil {
			continue
		}
		commitMessage := fmt.Sprintf("weasel: Report TODO (%s)", *issue.Todo.ReportedID)
		err = issue.Todo.CommitTodoUpdate(commitMessage)
		if err != nil {
			continue
		}
		if !quiet {
			issue.Todo.PrintCurrentStatus()
		}
	}
}

func purgeCommand(weasel *Weasel, jiraClient *JiraClient, quiet bool, bannerFunc func()) {
	assert.NotNil(weasel, "weasel can't be nil", "weasel", weasel)
	assert.NotNil(jiraClient, "jiraClient can't be nil", "jiraClient", jiraClient)
	if !quiet {
		bannerFunc()
    fmt.Fprintf(
      os.Stdout,
      " 🧹 Purging %s TODOS...\n\n",
      colors.Success("'DONE'"),
    )
	}
	todosToCheckStatus := []*Todo{}
	todosToPurge := []*Todo{}
	weasel.VisitTodosInWeaselFiles(func(td Todo) error {
		if td.ReportedID != nil {
			assert.NotNil(td.ReportedID, "Trying to purge an unreported todo", "td.ReportedID", td.ReportedID)
			todosToCheckStatus = append(todosToCheckStatus, &td)
		}
		return nil
	})
	for _, td := range todosToCheckStatus {
		status := jiraClient.CheckJiraIssueStatusFromAnExistingTodo(*td)
		if len(status) > 0 {
			if status == "done" || status == "DONE" {
				todosToPurge = append(todosToPurge, td)
				if !quiet {
					fmt.Println()
					td.PrintCurrentStatus()
				}
			}
		}
	}
	sort.Slice(todosToPurge, func(i, j int) bool {
		if todosToPurge[i].FilePath == todosToPurge[j].FilePath {
			return todosToPurge[i].Line > todosToPurge[j].Line
		}
		return todosToPurge[i].FilePath < todosToPurge[j].FilePath
	})
  if len(todosToPurge) == 0 {
    fmt.Fprintf(os.Stdout, colors.Info(" There are no TODOS to be purge\n"))
    os.Exit(0)
  }
	for _, td := range todosToPurge {
		err := td.SelfPurge()
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				colors.Error(fmt.Sprintf("Error trying to purge the todo of id: '%s'. Skipping for now.\n%s", *td.ReportedID, err)),
			)
		}
		commitMessage := fmt.Sprintf("weasel: Purge TODO (%s)", *td.ReportedID)
		err = td.CommitTodoUpdate(commitMessage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't commit changes: '%s'.\n", err)
		}
	}
}

func listCommand(weasel *Weasel, quiet bool, bannerFunc func()) {
	if !quiet {
		bannerFunc()
	}
	weasel.VisitTodosInWeaselFiles(func(t Todo) error {
		if t.ReportedID != nil {
			fmt.Println()
			t.PrintCurrentStatus()
			return nil
		}
		fmt.Println()
		t.PrintCurrentStatus()
		return nil
	})
}

func helperCommand() {
	flag.Usage()
}
