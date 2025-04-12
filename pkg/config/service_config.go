package config

import (
	"fmt"
	"os"

	"tsniper/pkg/logger"

	"gopkg.in/yaml.v3"
)

type Season struct {
	Term string `yaml:"term"`
	Year string `yaml:"year"`
}

type ServiceConfig struct {
	Campuses      []string `yaml:"campuses"`
	Seasons       []string `yaml:"seasons"`
	ValidSeasons  map[string]struct{}
	SeasonData    map[string]Season `yaml:"season_data"`
	DefaultCampus string
	DefaultSeason string
	logger        logger.Logger
}

func NewServiceConfig(file string, logger logger.Logger) *ServiceConfig {
	cfg := &ServiceConfig{
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
	cfg.DefaultCampus = cfg.Campuses[0]
	cfg.DefaultSeason = cfg.Seasons[0]
	cfg.ValidSeasons = make(map[string]struct{})
	for _, season := range cfg.Seasons {
		cfg.ValidSeasons[season] = struct{}{}
	}
	return cfg
}

func (c *ServiceConfig) load(file string) error {
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

func (c *ServiceConfig) validate() error {
	if len(c.Campuses) == 0 {
		return fmt.Errorf("empty campuses")
	}
	if len(c.Seasons) == 0 {
		return fmt.Errorf("empty seasons")
	}
	if len(c.SeasonData) == 0 {
		return fmt.Errorf("empty season data")
	}
	return nil
}
