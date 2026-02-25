package main

import (
	"fmt"
	"math"
	"strings"

	"tsniper/internal/config"
	"tsniper/internal/domain/product"
	"tsniper/internal/interfaces/discord/component"
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

	catalog := product.NewCatalog(product.TrialProduct, product.SemesterProduct, product.YearlyProduct, product.GraduateProduct)

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Onboarding, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{privacySettingsEmbed(cfg)},
	})
	if err != nil {
		logger.Fatal("Failed to send privacy settings announcement: %v", err)
	}

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Onboarding, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{gettingStartedEmbed(cfg)},
	})
	if err != nil {
		logger.Fatal("Failed to send getting started announcement: %v", err)
	}

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Pricing, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{pricingEmbed(cfg, catalog.GetAll()...)},
	})
	if err != nil {
		logger.Fatal("Failed to send pricing announcement: %v", err)
	}

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Autosniping, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{autoSnipingEmbed(cfg)},
	})
	if err != nil {
		logger.Fatal("Failed to send autosniping announcement: %v", err)
	}

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Autosniping, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{autoSnipingTipsEmbed(cfg)},
	})
	if err != nil {
		logger.Fatal("Failed to send autosniping tips announcement: %v", err)
	}

	_, err = s.ChannelMessageSendComplex(cfg.Discord.Channels.Admin, adminConsoleMessage(cfg))
	if err != nil {
		logger.Fatal("Failed to send admin console announcement: %v", err)
	}

	logger.Info("Announcements sent")
}

func privacySettingsEmbed(cfg *config.Config) *discordgo.MessageEmbed {
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
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cfg.Customization.Thumbnail,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/B4kCDbw.png",
		},
	}
}

func gettingStartedEmbed(cfg *config.Config) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Getting Started",
		Description: "Commands can be executed anywhere in this server or in your direct messages. " +
			"All responses from the bot are visible only to you and will eventually disappear.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "`/add {index}`",
				Value:  ">>> Add the index to your snipe requests.",
				Inline: false,
			},
			{
				Name:   "`/remove {index}`",
				Value:  ">>> Remove the index from your snipe requests.",
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
				Value:  ">>> View course information for the index.",
				Inline: false,
			},
		},
		Color: presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cfg.Customization.Thumbnail,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/ybyUvXT.png",
		},
	}
}

func pricingEmbed(cfg *config.Config, products ...product.Product) *discordgo.MessageEmbed {
	fields := make([]*discordgo.MessageEmbedField, 0, len(products)+1)
	for _, p := range products {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   pricingTitle(p),
			Value:  pricingDescription(p),
			Inline: false,
		})
	}
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "**Referrals**",
		Value:  ">>> Half price on your next purchase after your referent purchases a plan.",
		Inline: false,
	})

	return &discordgo.MessageEmbed{
		Title: "SnipeR Pricing",
		Description: "This bot only supports **Fall** and **Spring** terms. " +
			"Each semester purchased provides access to the autosniper until the end of the next add/drop period. " +
			"Use `/purchase` to select and purchase a plan.",
		Fields: fields,
		Color:  presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cfg.Customization.Thumbnail,
		},
	}
}

func pricingTitle(p product.Product) string {
	if p.Semesters == 1 {
		return fmt.Sprintf("**[%s Plan] %d Semester - $%.0f**", p.Name, p.Semesters, p.Price)
	}
	return fmt.Sprintf("**[%s Plan] %d Semesters - $%.0f**", p.Name, p.Semesters, p.Price)
}

func pricingDescription(p product.Product) string {
	if !strings.Contains(p.Description, "%v") {
		return p.Description
	}
	cost := p.Price / float64(p.Semesters)
	if cost == math.Trunc(cost) {
		return fmt.Sprintf(p.Description, int(cost))
	}
	return fmt.Sprintf(p.Description, fmt.Sprintf("%.2f", cost))
}

func autoSnipingEmbed(cfg *config.Config) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Getting Started with Autosniping",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Getting your Credentials",
				Value: ">>> Before enabling the autosniper, try logging into [WebReg](https://sims.rutgers.edu/webreg/pacLogin.htm) to ensure that you have the correct ruid and pac. " +
					"See the image below or this [page](https://scarlethub.rutgers.edu/registrar/personal-information-updates/pac-number-change/) for more information.",
				Inline: false,
			},
			{
				Name:   "Autosniper Management",
				Value:  ">>> Use `/console` to open the autosniper management console.",
				Inline: false,
			},
			{
				Name: "Managing your Credentials",
				Value: fmt.Sprintf(">>> Press <:green_id_card:%s> to add your credentials.\n", cfg.Customization.Emojis["green_id_card"]) +
					fmt.Sprintf("Press <:red_id_card:%s> to remove your credentials.", cfg.Customization.Emojis["red_id_card"]),
				Inline: false,
			},
			{
				Name: "Toggling the autosniper",
				Value: fmt.Sprintf(">>> Press <:green_gun:%s> to enable autosniping.\n", cfg.Customization.Emojis["green_gun"]) +
					fmt.Sprintf("Press <:red_gun:%s> to disable autosniping.", cfg.Customization.Emojis["red_gun"]),
				Inline: false,
			},
		},
		Color: presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cfg.Customization.Thumbnail,
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "https://i.imgur.com/6Z41zUS.png",
		},
	}
}

func autoSnipingTipsEmbed(cfg *config.Config) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Tips for Autosniping",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Plan Your Schedule",
				Value: ">>> The autosniper will attempt to register you for any course in your snipe list. " +
					"It is up to you to ensure that sections you provide are compatible with your current schedule.",
				Inline: false,
			},
			{
				Name: "Replacing Courses",
				Value: ">>> If you add a different section for a course you are registered in, " +
					"the autosniper will replace your current section with the new one when it opens.",
				Inline: false,
			},
		},
		Color: presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cfg.Customization.Thumbnail,
		},
	}
}

func adminConsoleMessage(cfg *config.Config) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{presentation.AdminConsoleEmbed(nil, cfg.Customization.Thumbnail)},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				component.AdminRefreshComponentDefinition(),
				component.AdminGrantComponentDefinition(),
				component.AdminRevokeComponentDefinition(),
				component.AdminGrantUnlimitedComponentDefinition(),
			}},
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				component.AdminUserSelectComponentDefinition(),
			}},
		},
	}
}
