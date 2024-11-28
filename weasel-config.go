package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/vitor-augusto1/jira-weasel/pkg/logger"
)

type Project struct {
	Url string `toml:"url"`
	Key string `toml:"key"`
}

type Credentials struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type Config struct {
	Remote      string            `toml:"remote"`
	Project     Project           `toml:"project"`
	Credentials Credentials       `toml:"credentials"`
	Keywords    map[string]string `toml:"keywords"`
}

func (c *Config) returnIssuesTypesSlice() []string {
	values := make([]string, 0, len(c.Keywords))
	for _, v := range c.Keywords {
		values = append(values, v)
	}
	return values
}

func (c *Config) returnIssuesKeywordsSlice() []string {
	values := make([]string, 0, len(c.Keywords))
	for key := range c.Keywords {
    values = append(values, strings.ToUpper(key))
	}
	return values
}

func (c *Config) configGotValidIssuesTypeInKeywords() {
	validOnes := map[string]bool{"Bug": true, "Task": true}
	issues := c.returnIssuesTypesSlice()
	for _, item := range issues {
		if !validOnes[item] {
			logger.LogErrorExitingOne("There is an invalid issue type in your config file.")
		}
	}
}

func CheckIfCurrentDirectoryIsAGitRepository() {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	if err != nil {
		logger.LogErrorExitingOne("Error trying to run jira-weasel outside of a git working tree")
	}
}

func CheckIfWeaselConfigFileExists() {
	if _, err := os.Stat("./weasel.toml"); errors.Is(err, os.ErrNotExist) {
		logger.LogErrorExitingOne("Error trying to run jira-weasel outside of a git working tree")
	}
}

func LoadConfigs(filePath string) (*Config, error) {
	var config Config
	config.Keywords = make(map[string]string)
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		return nil, err
	}
	// config.configGotValidIssuesTypeInKeywords()
	return &config, nil
}
