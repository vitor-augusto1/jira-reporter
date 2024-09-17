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
	IssueType string `yaml:"issueType"`
}

type JiraProjectYamlConfig struct {
	Jira     JiraProjectConfig `yaml:"jira"`
	Dirs     string            `yaml:"dirs"`
	RepoURL  string            `yaml:"repoUrl"`
	Keywords IssueTypeKeywords `yaml:"keywords"`
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
