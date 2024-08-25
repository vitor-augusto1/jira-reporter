package main

import (
	"fmt"
	"os"
)

var (
	PROJECT_BASE_URL = os.Getenv("PROJECT_BASE_URL")
)

func main() {
  creds := NewJiraBasicAuthCreds()
  creds.username = os.Getenv("PROJECT_USERNAME")
  creds.password = os.Getenv("PROJECT_PASSWORD")
  jc := NewJiraClient(creds, PROJECT_BASE_URL)
  projects, _ := jc.GetAllProjects()
  fmt.Println(string(projects))
}
