package colors

import "fmt"

var reset = "\033[0m" 
var Red = "\033[31m" 
var Green = "\033[32m" 
var Yellow = "\033[33m" 
var Blue = "\033[34m" 
var Magenta = "\033[35m" 
var Cyan = "\033[36m" 
var Gray = "\033[37m" 
var White = "\033[97m"

func Error(str string) string {
  return fmt.Sprintf("%s%s%s", Red, str, reset)
}

func Success(str string) string {
  return fmt.Sprintf("%s%s%s", Green, str, reset)
}

func Info(str string) string {
  return fmt.Sprintf("%s%s%s", Cyan, str, reset)
}

func Remote(str string) string {
  return fmt.Sprintf("%s%s%s", Magenta, str, reset)
}
