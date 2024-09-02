package main

import (
	"bufio"
	"fmt"
	"os"
)

type Weasel struct {
  keywords []string
}

// Walk a file searching for todos or comments with specific keywords
// and execute TodoTransformer
func (wl *Weasel) searchTodos(filePath string, ttr TodoTransformer) error {
  // TODO: Walk todos of a file and execute ttr (TodoTransformer)
  file, err := os.Open(filePath)
  if err != nil {
    return err
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    fmt.Println(line)
  }
  return nil
}
