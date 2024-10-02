package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Todo struct {
	Prefix      string
	Keyword     string
	Priority    PrioritiesID
	Title       string
	Body        []string
	FilePath    string
	Line        uint64
	RemoteAddr  string
	ReportedID  *string
}

type TodoTransformer func(Todo) error

func (td *Todo) LineHasTodoPrefix(line string) *string {
	if strings.HasPrefix(line, td.Prefix) {
		lineContent := strings.TrimPrefix(line, td.Prefix)
		return &lineContent
	}
	return nil
}

func (td *Todo) UpdatedTodoString(defaultStr string) string {
  if td.ReportedID != nil {
    updatedTodo := fmt.Sprintf(
      "%s %s P%s (%s): %s",
      td.Prefix,
      td.Keyword,
      td.Priority,
      *td.ReportedID,
      td.Title,
    )
    fmt.Println(">>> Updated todo content: ", updatedTodo)
    return updatedTodo
  }
  return defaultStr
}

func (td *Todo) StringBody() string {
	return strings.Join(td.Body, "\n")
}

func (td *Todo) ChangeTodoStatus() {
  fmt.Println("Changing the todo status...")
  tmpFileName := "tmp-wasel.weasel"
  tmpFile, err := os.Create(tmpFileName)
  if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Can't create the tmp file: '%s'. Skipping for now.\n",
        err,
			)
      return
  }
  defer tmpFile.Close()
  todoFile, err := os.Open(td.FilePath)
  if err != nil {
    fmt.Fprintf(
      os.Stderr,
      "Can't open the Todo file: '%s'. Skipping for now.\n",
      err,
    )
    return
  }
  defer todoFile.Close()
  todoFileInfo, _ := os.Stat(td.FilePath)
	scanner := bufio.NewScanner(todoFile)
  lnn := uint64(0)
  for scanner.Scan() {
    lnContent := scanner.Text()
    if td.Line == (lnn + 1) {
      fmt.Fprintln(tmpFile, td.UpdatedTodoString(lnContent))
    } else {
      fmt.Fprintln(tmpFile, lnContent)
    }
    lnn++
  }
  err = os.Chmod(tmpFileName, todoFileInfo.Mode())
  if err != nil {
    fmt.Fprintf(
      os.Stderr,
      "Can't set permissions: '%s'. Skipping for now.\n",
      err,
    )
    return
  }
  err = os.Rename(tmpFileName, td.FilePath)
  if err != nil {
    fmt.Fprintf(
      os.Stderr,
      "Can't rename the file: '%s'. Skipping for now.\n",
      err,
    )
    return
  }
}
