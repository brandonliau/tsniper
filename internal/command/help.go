package command

import (
	"fmt"
	"time"

	"Tsniper/internal/repository"
	"Tsniper/internal/shared"
	"Tsniper/pkg/config"

	"github.com/bwmarrin/discordgo"
)

type helpCommand struct {
	config *config.DiscordConfig
	repo   repository.Repository
	auth   bool
}

func NewHelpCommand(cfg config.Config, repo repository.Repository) *helpCommand {
	return &helpCommand{
		config: cfg.(*config.DiscordConfig),
		repo:   repo,
		auth:   false,
	}
}

func (c *helpCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "List all commands.",
	}
}

func (c *helpCommand) Auth() bool {
	return c.auth
}

func (c *helpCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	rsp := shared.EphemeralEmbedResponse(c.HelpEmbed())
	return rsp, nil
}

func (c *helpCommand) HelpEmbed() *discordgo.MessageEmbed {
	registered := c.repo.Registered()
	return &discordgo.MessageEmbed{
		Title: "Commands",
		Color: shared.Blue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Sniping",
				Value: fmt.Sprintf(
					"</add:%s> - Add a snipe request.\n"+
						"</remove:%s> - Remove a snipe request.\n"+
						"</clear:%s> - Remove all active snipe requests.\n"+
						"</check:%s> - View all active snipe requests.\n"+
						"</search:%s> - View course information for given index.",
					registered["add"],
					registered["remove"],
					registered["clear"],
					registered["check"],
					registered["search"],
				),
				Inline: false,
			},
			{
				Name: "Miscellaneous",
				Value: fmt.Sprintf(
					"</help:%s> - List all commands.\n"+
						"</uptime:%s> - Check bot uptime.\n"+
						"</ping:%s> - Check bot latency.",
					registered["help"],
					registered["uptime"],
					registered["ping"],
				),
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
