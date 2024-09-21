package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	ALL_PROJECTS_PATH string = "/project/search"
	NEW_ISSUE_PATH    string = "/issue"
)

type RequestStatusCode uint16

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
		fmt.Fprintf(
			os.Stderr,
			"Error. Cannot create a new request wrapper:\n%s",
			err,
		)
		return err
	}
	req.Header.Add("Authorization", "Basic "+jc.creds.ReturnEncodedCredentials())
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error. Cannot send the request:\n%s", err)
		return err
	}
  requestErr := jc.HandleResponseStatusCode(resp)
  if requestErr != nil {
    resp.Body.Close()
    fmt.Fprintf(os.Stderr, requestErr.Error())
    os.Exit(1)
  }
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: \n%s", err)
		return err
	}
	fmt.Fprintf(os.Stdout, "🔊 Issue reported: '%s'\n", issue.Summary)
	return nil
}

func (jc *JiraClient) HandleResponseStatusCode(resp *http.Response) *RequestError {
	requestError := &RequestError{}
	stCode := resp.StatusCode
	switch stCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		requestError.Message =
			"You do not have permission to create issues in this project. " +
				"Please check your credentials.\n"
		return requestError
	default:
		return nil
	}
}
