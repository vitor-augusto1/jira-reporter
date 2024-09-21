package main

type DefaultIssueTypesID uint16

const (
	TASK DefaultIssueTypesID = 10001
	BUG  DefaultIssueTypesID = 10002
)

type PrioritiesID string

const (
	HIGHEST PrioritiesID = "1"
	HIGH    PrioritiesID = "2"
	MEDIUM  PrioritiesID = "3"
	LOW     PrioritiesID = "4"
	LOWEST  PrioritiesID = "5"
)

type ID struct {
	ID string `json:"id"`
}

type Name struct {
	Name string `json:"name"`
}

type Key struct {
	Key string `json:"key"`
}

type ContentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type DescriptionParagraph struct {
	Type    string        `json:"type"`
	Content []ContentItem `json:"content"`
}

type IssueDescription struct {
	Type    string                 `json:"type"`
	Version uint8                  `json:"version"`
	Content []DescriptionParagraph `json:"content"`
}

type Issue struct {
	Description   IssueDescription `json:"description"`
	Priority      ID               `json:"priority"`
	IssueTypeName Name             `json:"issuetype"`
	Project       Key              `json:"project"`
	Summary       string           `json:"summary"`
  Todo          Todo
}

type IssuePayload struct {
	Fields *Issue `json:"fields"`
}

func NewIssue() *Issue {
	return &Issue{}
}
