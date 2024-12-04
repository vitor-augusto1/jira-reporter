package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/vitor-augusto1/jira-weasel/pkg/assert"
	"github.com/vitor-augusto1/jira-weasel/pkg/colors"
)

type Weasel struct {
	Files         []string
	Keywords      []string
	Todos         []Todo
	baseRemoteUrl string
}

// Walk a file searching for todos and execute TodoTransformer
func (wl *Weasel) searchTodos(filePath string, ttr TodoTransformer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var currentLineNumber uint64 = 0
	var currentOffset uint64 = uint64(0)
	var lineOffsets []uint64 = []uint64{}
	var todo *Todo
	for scanner.Scan() {
		currentLineNumber++
		line := scanner.Text()
		lineOffsets = append(lineOffsets, currentOffset)
		currentOffset += uint64(len(line)) + 1
		if todo == nil {
			todo = wl.GrabTodoFromLine(line, currentLineNumber, filePath)
			continue
		}
		lineIsPartOfTheTodoBody := todo.LineHasTodoPrefix(line)
		if lineIsPartOfTheTodoBody != nil {
      todo.EndLine = currentLineNumber
			todo.Body = append(todo.Body, *lineIsPartOfTheTodoBody)
		} else {
			wl.StoreTodoFullRemoteAddrs(todo)
			wl.Todos = append(wl.Todos, *todo)
			ttr(*todo)
			todo = nil
		}
	}
	if todo != nil {
		wl.StoreTodoFullRemoteAddrs(todo)
		wl.Todos = append(wl.Todos, *todo)
		ttr(*todo)
	}
	return nil
}

func (wl *Weasel) returnNewTodoFromLine(
	lineContent string,
	lineNumber uint64,
	filePath string,
) *Todo {
	for _, keyword := range wl.Keywords {
		todo := regexp.MustCompile(wl.TodoRegex(keyword))
		groups := todo.FindStringSubmatch(lineContent)
		if groups != nil {
			prefix := groups[1]
			keyword := groups[2]
			priority := groups[3]
			title := groups[4]
			return &Todo{
				Prefix:   prefix,
				Keyword:  keyword,
				Priority: PrioritiesID(priority),
				Title:    title,
				FilePath: filePath,
				Line:     lineNumber,
				Regex:    wl.TodoRegex(keyword),
			}
		}
	}
	return nil
}

func (wl *Weasel) returnReportedTodoFromLine(
	lineContent string,
	lineNumber uint64,
	filePath string,
) *Todo {
	for _, keyword := range wl.Keywords {
		todo := regexp.MustCompile(wl.ReportedTodoRegex(keyword))
		groups := todo.FindStringSubmatch(lineContent)
		if groups != nil {
			prefix := groups[1]
			keyword := groups[2]
			priority := groups[3]
			reportedId := groups[4]
			title := groups[5]
			return &Todo{
				Prefix:     prefix,
				Keyword:    keyword,
				Priority:   PrioritiesID(priority),
				Title:      title,
				FilePath:   filePath,
				Line:       lineNumber,
				EndLine:    lineNumber,
				ReportedID: &reportedId,
				Regex:      wl.ReportedTodoRegex(keyword),
			}
		}
	}
	return nil
}

func (wl Weasel) GrabTodoFromLine(
	lineContent string,
	lineNumber uint64,
	filePath string,
) *Todo {
	if todo := wl.returnNewTodoFromLine(lineContent, lineNumber, filePath); todo != nil {
		assert.Nil(
			todo.ReportedID,
			"An already reported todo was mistakenly returned as an unreported todo",
			"todo.ReportedID",
			todo.ReportedID,
		)
		return todo
	}
	if todo := wl.returnReportedTodoFromLine(lineContent, lineNumber, filePath); todo != nil {
		assert.NotNil(
			todo.ReportedID,
			"An unreported todo was mistakenly returned as a reported todo",
			"todo.ReportedID",
			todo.ReportedID,
		)
		return todo
	}
	return nil
}

func (wl Weasel) TodoRegex(keyword string) string {
	return "^(.*)" + "(" + regexp.QuoteMeta(keyword) + ")" + " P([1-5])" + ": (.*)$"
}

func (wl Weasel) ReportedTodoRegex(keyword string) string {
	return "^(.*)" + "(" + regexp.QuoteMeta(keyword) + ")" + " P([1-5]) " + "\\((.*)\\)" + ": (.*)$"
}

func (wl Weasel) RemoteIsAGithubRepo(str string) bool {
	matched, _ := regexp.MatchString("https://github.", str)
	return matched
}

func (wl Weasel) RemoteIsAGitlabRepo(str string) bool {
	matched, _ := regexp.MatchString("https://gitlab.", str)
	return matched
}

func (wl Weasel) GetProjectCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	stdout, err := cmd.Output()
	assert.NoError(err, "Error getting the current git branch. Please check if the current directory is a git work tree.")
	cb := strings.TrimSpace(string(stdout))
	return cb
}

func (wl *Weasel) LoadProjectFiles() {
	cmd := exec.Command("git", "ls-files", "--full-name")
	stdout, err := cmd.Output()
	assert.NoError(err, "Error getting the project's files. Make sure the current directory is a git work tree.\n")
	output := strings.TrimSpace(string(stdout))
	wl.Files = strings.Split(output, "\n")
}

func (wl *Weasel) GetRemoteBlobPath(filePath string, line uint64) string {
	assert.Assert(len(filePath) > 0, "invalid file path was provided")
	assert.Assert(line > 0, "invalid line number was provided")
	if wl.RemoteIsAGithubRepo(wl.baseRemoteUrl) {
		return fmt.Sprintf(
			"%s/blob/%s/%s/#L%d",
			wl.baseRemoteUrl, wl.GetProjectCurrentBranch(), filePath, line,
		)
	}
	if wl.RemoteIsAGitlabRepo(wl.baseRemoteUrl) {
		return fmt.Sprintf(
			"%s/-/blob/%s/%s#L%d",
			wl.baseRemoteUrl, wl.GetProjectCurrentBranch(), filePath, line,
		)
	}
	failedMessage := fmt.Sprintf(
		"Failed to store the remote todo blob. Remote is not a github nor a gitlab repo. Remote: %s",
		wl.baseRemoteUrl,
	)
	fmt.Fprintf(os.Stderr, colors.Error(failedMessage))
	return ""
}

func (wl Weasel) StoreTodoFullRemoteAddrs(todo *Todo) {
	assert.NotNil(todo, "expected todo to be a *Todo but got nil", "todo", todo)
	remoteAddr := wl.GetRemoteBlobPath(todo.FilePath, todo.Line)
	todo.RemoteAddr = remoteAddr
	todo.Body = append(todo.Body, todo.RemoteAddr)
}

func (wl *Weasel) VisitTodosInWeaselFiles(ttr TodoTransformer) error {
	for _, file := range wl.Files {
		wl.searchTodos(file, ttr)
	}
	return nil
}
