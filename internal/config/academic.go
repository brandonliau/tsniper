package config

import (
	"fmt"

	"tsniper/internal/domain/scope"
)

type AcademicConfig struct {
	Campuses []scope.Campus `yaml:"campuses"`
	Seasons  []scope.Season `yaml:"seasons"`
}

func (c *AcademicConfig) Validate() error {
	if len(c.Campuses) == 0 {
		return fmt.Errorf("campuses are required")
	}
	if len(c.Seasons) == 0 {
		return fmt.Errorf("seasons are required")
	}
	return nil
}
