package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	ALL_PROJECTS_PATH string = "/project/search"
	NEW_ISSUE_PATH    string = "/issue"
)

type JiraClient struct {
	creds   *JiraBasicAuthCreds
	baseURL string
}

func (jc *JiraClient) CreateNewIssueFromTODO(td Todo) *Issue {
	newIssue := NewIssue()
	newIssue.Description = IssueDescription{
		Type:    "doc",
		Version: 1,
		Content: []DescriptionParagraph{
			{
				Type: "paragraph",
				Content: []ContentItem{
					{
						Type: "text",
						Text: td.StringBody(),
					},
				},
			},
		},
	}
	newIssue.Priority = ID{ID: string(td.Priority)}
	newIssue.IssueTypeName = Name{Name: "Bug"}
	newIssue.Project = Key{Key: "SCRUM"}
	newIssue.Summary = td.Title 
  return newIssue
}

// Creates a new jira issue
func (jc *JiraClient) ReportIssueAsJiraTicket(issue *Issue) error {
	var CREATE_ISSUE_URL string = jc.baseURL + NEW_ISSUE_PATH
	client := &http.Client{}
	structPayload := IssuePayload{
		Fields: issue,
	}
	body, _ := json.Marshal(structPayload)
	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, CREATE_ISSUE_URL, payload)
	if err != nil {
		fmt.Println("Error. Cannot create a new request wrapper: ", err)
		return err
	}
	req.Header.Add("Authorization", "Basic "+jc.creds.ReturnEncodedCredentials())
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error. Cannot send request: ", err)
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return err
	}
	return nil
}
