package config

import (
	"fmt"
	"os"

	"Tsniper/pkg/logger"

	"gopkg.in/yaml.v3"
)

type Season struct {
	Term string `yaml:"term"`
	Year string `yaml:"year"`
}

type SnipeConfig struct {
	Campuses   []string          `yaml:"campuses"`
	Seasons    []string          `yaml:"seasons"`
	SeasonData map[string]Season `yaml:"season_data"`
	logger     logger.Logger
}

func NewSnipeConfig(file string, logger logger.Logger) *SnipeConfig {
	cfg := &SnipeConfig{
		logger: logger,
	}
	err := cfg.load(file)
	if err != nil {
		logger.Fatal("Failed to load config file: %v", err)
	}
	err = cfg.validate()
	if err != nil {
		logger.Fatal("Failed to validate config file: %v", err)
	}
	return cfg
}

func (c *SnipeConfig) load(file string) error {
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("readfile: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("unmarshal: %v", err)
	}
	return nil
}

func (c *SnipeConfig) validate() error {
	if c.Campuses == nil {
		return fmt.Errorf("empty campuses")
	}
	if c.Seasons == nil {
		return fmt.Errorf("empty seasons")
	}
	if c.SeasonData == nil {
		return fmt.Errorf("empty season data")
	}
	return nil
}
