package config

import (
	"fmt"
)

type CustomizationConfig struct {
	Thumbnail string            `yaml:"thumbnail"`
	Emojis    map[string]string `yaml:"emojis"`
}

func (c *CustomizationConfig) Validate() error {
	if c.Thumbnail == "" {
		return fmt.Errorf("thumbnail is required")
	}
	if len(c.Emojis) == 0 {
		return fmt.Errorf("emojis are required")
	}
	return nil
}
