package colors

import "fmt"

var reset = "\033[0m" 
var red = "\033[31m" 
var green = "\033[32m" 
var yellow = "\033[33m" 
var blue = "\033[34m" 
var magenta = "\033[35m" 
var cyan = "\033[36m" 
var gray = "\033[37m" 
var white = "\033[97m"

func Error(str string) string {
  return fmt.Sprintf("%s%s\n%s", red, str, reset)
}
