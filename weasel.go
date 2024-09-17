package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Weasel struct {
	Keywords      []string
	Todos         []Todo
	baseRemoteUrl string
}

// Walk a file searching for todos or comments with specific keywords
// and execute TodoTransformer
func (wl *Weasel) searchTodos(filePath string, ttr TodoTransformer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lineNumber uint32 = 0
	var todo *Todo
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if todo == nil {
			todo = wl.returnTodoFromLine(line, lineNumber, filePath)
			continue
		}
		lineIsPartOfTheTodoBody := todo.LineHasTodoPrefix(line)
		if lineIsPartOfTheTodoBody != nil {
			todo.Body = append(todo.Body, *lineIsPartOfTheTodoBody)
		} else {
      wl.StoreTodoFullRemoteAddrs(todo)
			wl.Todos = append(wl.Todos, *todo)
			ttr(*todo)
			todo = nil
		}
	}
	return nil
}

func (wl *Weasel) returnTodoFromLine(
	lineContent string,
	lineNumber uint32,
	filePath string,
) *Todo {
	for _, keyword := range wl.Keywords {
		todo := regexp.MustCompile(wl.todoRegex(keyword))
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
			}
		}
	}
	return nil
}

func (wl Weasel) todoRegex(keyword string) string {
	return "^(.*)" + "(" + regexp.QuoteMeta(keyword) + ")" + " P([1-5])" + ": (.*)$"
}

func (wl Weasel) GetProjectCurrentBranch() string {
	cmd := exec.Command("git", "branch", "--show-current")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting the current git branch.\n%s\n", err)
		return ""
	}
	cb := strings.TrimSpace(string(stdout))
	return cb
}

func (wl Weasel) StoreTodoFullRemoteAddrs(todo *Todo) {
	rtAddr := fmt.Sprintf(
		"%s/blob/%s/%s/#L%d",
		wl.baseRemoteUrl, wl.GetProjectCurrentBranch(), todo.FilePath, todo.Line,
	)
  todo.RemoteAddr = rtAddr
  todo.Body = append(todo.Body, todo.RemoteAddr)
}
