package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Discord       *DiscordConfig       `yaml:"discord"`
	Stripe        *StripeConfig        `yaml:"stripe"`
	Customization *CustomizationConfig `yaml:"customization"`
	Academic      *AcademicConfig      `yaml:"academic"`
}

func NewYamlConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Validate() error {
	err := c.Discord.Validate()
	if err != nil {
		return fmt.Errorf("Discord: %v", err)
	}
	err = c.Stripe.Validate()
	if err != nil {
		return fmt.Errorf("Stripe: %v", err)
	}
	err = c.Customization.Validate()
	if err != nil {
		return fmt.Errorf("Customization: %v", err)
	}
	err = c.Academic.Validate()
	if err != nil {
		return fmt.Errorf("Academic: %v", err)
	}
	return nil
}
