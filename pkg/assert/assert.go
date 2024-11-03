package assert

import (
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime/debug"
)

func runAssert(msg string, args ...any) {
	slogValuesMessages := []any{
    "Assert message",
		msg,
	}
	slogValuesMessages = append(slogValuesMessages, args...)
	for i := 0; i < len(slogValuesMessages); i += 2 {
    fmt.Fprintf(os.Stderr, "  %s: %v\n", slogValuesMessages[i], slogValuesMessages[i+1])
	}
	fmt.Fprintln(os.Stderr, string(debug.Stack()))
	os.Exit(1)
}

func Assert(truthy bool, msg string, args ...any) {
  if !truthy {
    runAssert(msg, args...)
  }
}

func NotNil(item any, msg string, args ...any) {
	if item == nil || reflect.ValueOf(item).Kind() == reflect.Ptr && reflect.ValueOf(item).IsNil() {
		slog.Error("NotNil#nil value encountered")
    runAssert(msg, args...)
	}
}

func NoError(err error, msg string, data ...any) {
  if err != nil {
    data = append(data, "error", err)
    runAssert(msg, data...)
  }
}
