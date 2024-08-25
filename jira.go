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

// Returns the jira api response with the projects
func (jc *JiraClient) GetAllProjects() ([]byte, error) {
	var ALL_PROJECTS_URL string = jc.baseURL + ALL_PROJECTS_PATH
	client := &http.Client{}
	req, _ := http.NewRequest(http.MethodGet, ALL_PROJECTS_URL, nil)
	req.Header.Add("Authorization", "Basic "+jc.creds.ReturnEncodedCredentials())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return nil, err
	}
	return body, nil
}

// Creates a new jira issue
func (jc *JiraClient) CreateNewIssue(issue *Issue) error {
	var CREATE_ISSUE_URL string = jc.baseURL + NEW_ISSUE_PATH
	client := &http.Client{}
	body, _ := json.Marshal(issue)
	payload := bytes.NewBuffer(body)
	req, err := http.NewRequest(http.MethodPost, CREATE_ISSUE_URL, payload)
	if err != nil {
		fmt.Println("Error. Cannot create a new request wrapper: ", err)
		return err
	}
	req.Header.Add("Authorization", "Basic "+jc.creds.ReturnEncodedCredentials())
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
	fmt.Println(string(body))
	return nil
}

// Returns new instance of JiraClient
func NewJiraClient(creds *JiraBasicAuthCreds, bURL string) *JiraClient {
	return &JiraClient{
		creds:   creds,
		baseURL: bURL,
	}
}
