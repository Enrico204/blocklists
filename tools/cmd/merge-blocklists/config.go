package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type mergedListConfig struct {
	Description string   `yaml:"description"`
	Include     []string `yaml:"include"`
}

type config struct {
	ListsDirectory string                      `yaml:"lists_directory"`
	Lists          map[string]mergedListConfig `yaml:"lists"`
}

func loadConfig(fname string) (config, error) {
	var cfg config

	fp, err := os.Open(fname)
	if err != nil {
		return cfg, err
	}
	defer func() { _ = fp.Close() }()

	err = yaml.NewDecoder(fp).Decode(&cfg)
	return cfg, err
}
