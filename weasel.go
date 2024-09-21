package main

import (
	"bufio"
	"fmt"
	"io"
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
			todo = wl.returnTodoFromLine(line, currentLineNumber, filePath)
			continue
		}
		lineIsPartOfTheTodoBody := todo.LineHasTodoPrefix(line)
		if lineIsPartOfTheTodoBody != nil {
			todo.Body = append(todo.Body, *lineIsPartOfTheTodoBody)
		} else {
			wl.StoreTodoFullRemoteAddrs(todo)
			wl.Todos = append(wl.Todos, *todo)
			err = ttr(*todo)
			// If the TodoTransformer executes successfully we change the todo
			// line to indicate that it has been reported
      if err == nil {}
			todo = nil
		}
	}
	return nil
}

// Back tracks to a specific line based on the offset
func (wl *Weasel) BackTrackLine(
  file *os.File,
  offset int64,
  handler func(string) string,
) error {
	_, err := file.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf(
      "Could not back track to the line. Byte offset %d: %v", offset, err,
    )
	}
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
    fmt.Println("BACKTRACTED LINE CONTENT ::: ", scanner.Text())
    handler(scanner.Text())
		return nil
	}
	return nil
}

func (wl *Weasel) returnTodoFromLine(
	lineContent string,
	lineNumber uint64,
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
