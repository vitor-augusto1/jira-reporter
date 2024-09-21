package main

import "strings"

type Todo struct {
	Prefix      string
	Keyword     string
	Priority    PrioritiesID
	Title       string
	Body        []string
	FilePath    string
	Line        uint64
	RemoteAddr  string
	ReportedID  string
}

type TodoTransformer func(Todo) error

func (td *Todo) LineHasTodoPrefix(line string) *string {
	if strings.HasPrefix(line, td.Prefix) {
		lineContent := strings.TrimPrefix(line, td.Prefix)
		return &lineContent
	}
	return nil
}

func (td *Todo) StringBody() string {
	return strings.Join(td.Body, "\n")
}

func (td *Todo) ChangeTodoStatus() {
  // TODO: Change the todo inside the file with the id of the reported todo
  // to identify that the todo was reported
}
