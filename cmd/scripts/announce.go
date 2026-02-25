package main

import (
	"tsniper/internal/config"
	"tsniper/internal/interfaces/discord/presentation"

	"tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func announce() {
	logger := logger.NewStdLogger(logger.LevelInfo)

	cfg, err := config.NewYamlConfig("./config/config.yml")
	if err != nil {
		logger.Fatal("Failed to create config: %v", err)
	}

	s, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection: %v", err)
	}
	defer s.Close()

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Onboarding, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{privacySettingsEmbed(cfg.Customization)},
	})
	if err != nil {
		logger.Fatal("Failed to send privacy settings announcement: %v", err)
	}

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Onboarding, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{gettingStartedEmbed(cfg.Customization)},
	})
	if err != nil {
		logger.Fatal("Failed to send getting started announcement: %v", err)
	}

	logger.Info("Announcements sent")
}

func privacySettingsEmbed(cfg *config.CustomizationConfig) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Privacy Settings",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Step 1",
				Value:  ">>> Enable `allow direct messages from server members` in your `Privacy & Safety` settings",
				Inline: false,
			},
			{
				Name:   "Step 2",
				Value:  ">>> Enable `allow direct messages from server members` in this server's `Privacy Settings`",
				Inline: false,
			},
		},
		Color: presentation.Blue,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/B4kCDbw.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cfg.Thumbnail,
		},
	}
}

func gettingStartedEmbed(cfg *config.CustomizationConfig) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Getting Started",
		Description: "Commands can be excecuted anywhere in this server or in your direct messages. All responses from the bot are visible only to you and will eventually disappear. When a section in your snipe requests opens, tsniper will DM you and provide a link that autofills the course index in WebReg.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "`/add {index}`",
				Value:  ">>> Adds the given index to your snipe requests.",
				Inline: false,
			},
			{
				Name:   "`/remove {index}`",
				Value:  ">>> Remove the given index from your snipe requests.",
				Inline: false,
			},
			{
				Name:   "`/clear`",
				Value:  ">>> Remove all indices from your snipe requests.",
				Inline: false,
			},
			{
				Name:   "`/check`",
				Value:  ">>> View all active snipe requests.",
				Inline: false,
			},
			{
				Name:   "`/search {index}`",
				Value:  ">>> View course information for given index.",
				Inline: false,
			},
		},
		Color: presentation.Blue,
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/tCd1v6q.png",
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cfg.Thumbnail,
		},
	}
}
