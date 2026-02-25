package config

import (
	"fmt"
)

type DiscordConfig struct {
	Token         string            `yaml:"token"`
	ApplicationID string            `yaml:"application_id"`
	GuildID       string            `yaml:"guild_id"`
	Roles         map[string]string `yaml:"roles"`
	Channels      struct {
		Boarding   string `yaml:"boarding"`
		Onboarding string `yaml:"onboarding"`
	} `yaml:"channels"`
}

func (c *DiscordConfig) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	if c.ApplicationID == "" {
		return fmt.Errorf("application ID is required")
	}
	if c.GuildID == "" {
		return fmt.Errorf("guild ID is required")
	}
	if len(c.Roles) == 0 {
		return fmt.Errorf("roles are required")
	}
	if c.Channels.Boarding == "" {
		return fmt.Errorf("boarding channel ID is required")
	}
	if c.Channels.Onboarding == "" {
		return fmt.Errorf("onboarding channel ID is required")
	}
	return nil
}
