package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/vitor-augusto1/jira-weasel/pkg/assert"
)

const (
	ALL_PROJECTS_PATH string = "/project/search"
	NEW_ISSUE_PATH    string = "/issue"
)

type CreatedIssueResponse struct {
  Id   string `json:"id"`
  Key  string `json:"key"`
  Self string `json:"self"`
}

type RequestStatusCode uint16

type JiraClient struct {
	creds   *JiraBasicAuthCreds
	baseURL string
}

func (jc JiraClient) CreateNewIssueFromTODO(td Todo, issueTp IssueType) *Issue {
  assert.Nil(td.ReportedID, "Already reported todo passed", "ReportedId", *td.ReportedID)
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
	newIssue.IssueTypeName = Name{Name: issueTp}
	newIssue.Project = Key{Key: "SCRUM"}
	newIssue.Summary = td.Title
  newIssue.Todo = td
	return newIssue
}

// Creates a new jira issue
func (jc *JiraClient) ReportIssueAsJiraTicket(issue *Issue) (*CreatedIssueResponse, error) {
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
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+jc.creds.ReturnEncodedCredentials())
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error. Cannot send the request:\n%s", err)
		return nil, err
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
		return nil, err
	}
  var createdIssueResp *CreatedIssueResponse
  err = json.Unmarshal(body, &createdIssueResp)
  if err != nil {
		fmt.Fprintf(os.Stderr, "💢 %s\n", err)
		os.Exit(1)
  }
	fmt.Fprintf(os.Stdout, "🔊 Issue reported: '%s'\n", issue.Summary)
	return createdIssueResp, nil
}

func (jc *JiraClient) HandleResponseStatusCode(resp *http.Response) *RequestError {
	requestError := &RequestError{}
	responseStatusCode := resp.StatusCode
	switch responseStatusCode {
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
