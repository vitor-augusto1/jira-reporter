package main

import (
	"net/http"

	"github.com/vitor-augusto1/jira-weasel/client"
	"github.com/vitor-augusto1/jira-weasel/pkg/assert"
)

type CreatedIssueResponse struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

type IssueStatusResponse struct {
	Fields struct {
		Status struct {
			StatusCategory struct {
				Self      string `json:"self"`
				Id        uint64 `json:"id"`
				Key       string `json:"key"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
	} `json:"fields"`
}

type RequestStatusCode uint16

type JiraClient struct {
	HttpClient *client.HttpClient
}

// Returns new instance of JiraClient
func NewJiraClient(bURL, encodedCredentials string) *JiraClient {
	return &JiraClient{
		HttpClient: client.NewHttpClient(bURL, "Basic "+encodedCredentials),
	}
}

func (jc JiraClient) CreateNewIssueFromTODO(td Todo, issueTp IssueType) *Issue {
	assert.Nil(td.ReportedID, "Already reported todo passed", "ReportedId", td.ReportedID)
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

func (jc *JiraClient) CheckJiraIssueStatusFromAnExistingTodo(td Todo) string {
	assert.NotNil(td.ReportedID, "checking issue status of an unreported todo", "td.ReportedId", *td.ReportedID)
	var issueStatusResponse IssueStatusResponse
	requestOpts := &client.RequestOptions{
		Method:   http.MethodGet,
		Path:     "/issue/" + *td.ReportedID + "?fields=status",
		Payload:  nil,
		Response: &issueStatusResponse,
		Headers:  nil,
	}
	jc.HttpClient.DoRequest(requestOpts, func(err error) {
		assert.NoError(err, "Error trying to check an issue status")
	})
	return issueStatusResponse.Fields.Status.StatusCategory.Key
}

// Creates a new jira issue
func (jc *JiraClient) ReportIssueAsJiraTicket(issue *Issue) (*CreatedIssueResponse, error) {
  assert.NotNil(issue, "issue parameter cannot be nil", "issue", issue)
	var createdIssueResp CreatedIssueResponse
	requestOpts := &client.RequestOptions{
		Method:   http.MethodPost,
		Path:     "/issue",
		Payload:  IssuePayload{Fields: issue},
		Response: &createdIssueResp,
		Headers:  nil,
	}
	jc.HttpClient.DoRequest(requestOpts, func(err error) {
		assert.NoError(err, "Error trying to check an issue status")
	})
	return &createdIssueResp, nil
}

// TODO: Update this ugly ass function
func (jc JiraClient) HandleResponseStatusCode(resp *http.Response) *RequestError {
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
