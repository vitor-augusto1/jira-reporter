package main

type Todo struct {
	Keyword  string
	Priority PrioritiesID
	Title    string
	Body     []string
}

type TodoTransformer func(Todo) error
