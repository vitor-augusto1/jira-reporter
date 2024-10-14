package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

type JiraProjectMainConfig struct {
	Name    string `yaml:"name"`
	BaseUrl string `yaml:"baseUrl"`
}

type JiraProjectConfig struct {
	Project JiraProjectMainConfig `yaml:"project"`
}

type IssueTypeKeywords map[string]struct {
	IssueType IssueType `yaml:"issueType"`
}

type JiraProjectYamlConfig struct {
	Jira     JiraProjectConfig `yaml:"jira"`
	Dirs     string            `yaml:"dirs"`
	RepoURL  string            `yaml:"repoUrl"`
	Keywords IssueTypeKeywords `yaml:"keywords"`
}

func (jpy *JiraProjectYamlConfig) ReturnKeywordSlice() []string {
  keywordSlice := make([]string, len(jpy.Keywords))
  i := 0
  for key := range jpy.Keywords {
    keywordSlice[i] = key
    i++
  }
  return keywordSlice
}

func (jpy *JiraProjectYamlConfig) ReturnIssuesTypesSlice() []string {
  issuesSlice := make([]string, len(jpy.Keywords))
  i := 0
  for _, key := range jpy.Keywords {
    issuesSlice[i] = key.IssueType
    i++
  }
  return issuesSlice
}

func parseYamlConfigFile(filePath string) (*JiraProjectYamlConfig, error) {
	var jpy JiraProjectYamlConfig
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &jpy)
	if err != nil {
		return nil, err
	}
	return &jpy, nil
}
