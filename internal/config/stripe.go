package config

import (
	"fmt"
)

type StripeConfig struct {
	Token         string `yaml:"token"`
	WebhookSecret string `yaml:"webhook_secret"`
	SuccessURL    string `yaml:"success_url"`
	CancelURL     string `yaml:"cancel_url"`
	Port          string `yaml:"port"`
}

func (c *StripeConfig) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token is required")
	}
	if c.WebhookSecret == "" {
		return fmt.Errorf("webhook secret is required")
	}
	if c.SuccessURL == "" {
		return fmt.Errorf("success URL is required")
	}
	if c.CancelURL == "" {
		return fmt.Errorf("cancel URL is required")
	}
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}
	return nil
}
