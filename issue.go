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
	IssueTypeName Name             `json:"issueType"`
	Project       Key              `json:"project"`
	Summary       string           `json:"summary"`
}

func NewIssue() *Issue {
	return &Issue{}
}
