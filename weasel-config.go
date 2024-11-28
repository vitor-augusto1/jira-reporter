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
