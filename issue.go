package main

type DefaultIssueTypesID uint16

const (
	TASK DefaultIssueTypesID = 10001
	BUG  DefaultIssueTypesID = 10002
)

type PrioritiesID uint16

const (
	HIGHEST PrioritiesID = iota + 1
	HIGH
	MEDIUM
	LOW
	LOWEST
)

type ID struct {
	ID string `json:"id"`
}

type ContentItem struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type DescriptionParagraph struct {
	Content []ContentItem `json:"content"`
	Type    string        `json:"type"`
}

type IssueDescription struct {
	Content []DescriptionParagraph `json:"content"`
	Type    string                 `json:"type"`
	Version uint8                  `json:"version"`
}

type Issue struct {
	IssueTypeID ID               `json:"issueType"`
	Priority    ID               `json:"priority"`
	Project     ID               `json:"project"`
	Summary     string           `json:"summary"`
	Description IssueDescription `json:"description"`
}

func NewIssue() *Issue {
	return &Issue{}
}
