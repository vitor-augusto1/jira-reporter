package logger

import (
	"fmt"
	"os"

	"github.com/vitor-augusto1/jira-weasel/pkg/colors"
)

// Exit logging an error message
func LogErrorExitingOne(message string) {
	fmt.Fprintf(
    os.Stderr,
    colors.Error(message),
  )
	os.Exit(1)
}
