package command

import (
	"time"

	"tsniper/internal/application/usecase"
	"tsniper/internal/config"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type clearCommand struct {
	snipeService  *usecase.SnipeService
	customization *config.CustomizationConfig
}

func ClearCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "clear",
		Description: "Remove all active snipe requests.",
	}
}

func ClearCommandHandler(snipeService *usecase.SnipeService, customization *config.CustomizationConfig) interaction.HandleFunc {
	c := &clearCommand{
		snipeService:  snipeService,
		customization: customization,
	}
	return c.execute
}

func (c *clearCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	req := usecase.ClearSnipeRequest{
		UserID: interaction.GetUserID(i),
	}

	res, err := c.snipeService.Clear(req)
	if err != nil {
		return nil, err
	}

	if res.Count == 0 {
		rsp := interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.invalidClear()),
			interaction.WithEphemeral(),
		)
		return rsp, nil
	}

	rsp := interaction.InteractionInitialResponse(
		interaction.WithEmbeds(c.successfulClear()),
		interaction.WithEphemeral(),
	)
	return rsp, nil
}

func (c *clearCommand) successfulClear() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Success!",
		Description: "All active snipes requests have been removed.",
		Color:       presentation.Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *clearCommand) invalidClear() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: "You have no active snipe requests.",
		Color:       presentation.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
