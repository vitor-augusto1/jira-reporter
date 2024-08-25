package main

import (
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
  newIssue := NewIssue()
  newIssue.Description = IssueDescription{
    Type: "doc",
    Version: 1,
    Content: []DescriptionParagraph{
      {
        Type: "paragraph",
        Content: []ContentItem{
          {
            Type: "text",
            Text: "JIRA-WEASEL body test",
          },
        },
      },
    },
  }
  newIssue.Priority = ID{ID: string(HIGHEST)}
  newIssue.IssueTypeName = Name{Name: "Bug"}
  newIssue.Project = Key{Key: "SCRUM"}
  newIssue.Summary = "JIRA-WEASEL TESTE 2"
  _ = jc.CreateNewIssue(newIssue)
}
