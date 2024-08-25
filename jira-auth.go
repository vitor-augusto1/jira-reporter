package main

import "encoding/base64"

type JiraBasicAuthCreds struct {
	username string
	password string
}

// Returns jira encoded credentials
func (b *JiraBasicAuthCreds) ReturnEncodedCredentials() string {
  authStr := b.username + ":" + b.password
  return base64.StdEncoding.EncodeToString([]byte(authStr))
}

// Returns new instance of NewJiraBasicAuthCreds
func NewJiraBasicAuthCreds() *JiraBasicAuthCreds {
  return &JiraBasicAuthCreds{}
}
