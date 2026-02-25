package config

import (
	"fmt"
)

type CustomizationConfig struct {
	Thumbnail string `yaml:"thumbnail"`
}

func (c *CustomizationConfig) Validate() error {
	if c.Thumbnail == "" {
		return fmt.Errorf("thumbnail is required")
	}
	return nil
}
