package config

import (
	_ "embed"
	"encoding/json"
)

//go:embed groups.json
var groupsJSON []byte

type Group struct {
	Prefix string `json:"prefix"`
	Output string `json:"output"`
}

type Config struct {
	Groups  []Group `json:"groups"`
	WorkDir string  `json:"work_dir"`
}

func LoadConfig() (*Config, error) {
	var config Config
	err := json.Unmarshal(groupsJSON, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) GetGroups() []Group {
	return c.Groups
}

func (c *Config) GetWorkDir() string {
	if c.WorkDir == "" {
		return "~/Downloads"
	}
	return c.WorkDir
}