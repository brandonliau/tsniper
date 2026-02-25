package config

import (
	"fmt"
)

type DiscordConfig struct {
	Token         string `yaml:"token"`
	ApplicationID string `yaml:"application_id"`
	GuildID       string `yaml:"guild_id"`
	AdminID       string `yaml:"admin_id"`
	Roles         struct {
		Unlimited string `yaml:"unlimited"`
		Standard  string `yaml:"standard"`
	} `yaml:"roles"`
	Channels struct {
		Admin       string `yaml:"admin"`
		Onboarding  string `yaml:"onboarding"`
		Autosniping string `yaml:"autosniping"`
		Pricing     string `yaml:"pricing"`
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
	if c.AdminID == "" {
		return fmt.Errorf("admin ID is required")
	}
	if c.Roles.Unlimited == "" {
		return fmt.Errorf("unlimited role ID is required")
	}
	if c.Roles.Standard == "" {
		return fmt.Errorf("standard role ID is required")
	}
	if c.Channels.Onboarding == "" {
		return fmt.Errorf("onboarding channel ID is required")
	}
	if c.Channels.Autosniping == "" {
		return fmt.Errorf("autosniping channel ID is required")
	}
	if c.Channels.Pricing == "" {
		return fmt.Errorf("pricing channel ID is required")
	}
	return nil
}
