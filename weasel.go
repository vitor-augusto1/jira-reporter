package main

type Weasel struct {
  keywords []string
}

// Walk a file searching for todos or comments with specific keywords
// and execute TodoTransformer
func (wl *Weasel) searchTodos(filePath string, ttr TodoTransformer) error {
  // TODO: Walk todos of a file and execute ttr (TodoTransformer)
  return nil
}
