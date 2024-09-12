package main

import "strings"

type Todo struct {
	Prefix   string
	Keyword  string
	Priority PrioritiesID
	Title    string
	Body     []string
	FilePath string
	Line     uint32
  // TODO: Add a context field to this struct to hold the context around the todo
  // Maybe taking every contiguous text below the todo and set it as a context.
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
