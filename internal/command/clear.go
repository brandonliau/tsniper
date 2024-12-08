package command

import (
	"time"

	"Tsniper/internal/repository"
	"Tsniper/internal/shared"
	"Tsniper/pkg/config"
	"Tsniper/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type clearCommand struct {
	config *config.DiscordConfig
	repo   repository.Repository
	db     database.Database
	auth   bool
}

func NewClearCommand(cfg config.Config, repo repository.Repository, db database.Database) *clearCommand {
	return &clearCommand{
		config: cfg.(*config.DiscordConfig),
		repo:   repo,
		db:     db,
		auth:   false,
	}
}

func (c *clearCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "clear",
		Description: "Remove all active snipe requests.",
	}
}

func (c *clearCommand) Auth() bool {
	return c.auth
}

func (c *clearCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	snipes := c.repo.Snipes(args.UserID)
	if len(snipes) == 0 {
		rsp := shared.EphemeralEmbedResponse(c.InvalidClear())
		return rsp, nil
	}

	c.repo.ClearSnipe(args.UserID)
	for _, snipe := range snipes {
		index, campus, season := snipe[0], snipe[1], snipe[2]
		c.repo.Remove(index, campus, season)
	}

	rsp := shared.EphemeralEmbedResponse(c.SuccessfulClear())
	return rsp, nil
}

func (c *clearCommand) InvalidClear() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: "You have no active snipe requests.",
		Color:       shared.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *clearCommand) SuccessfulClear() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Success!",
		Description: "All active snipes requests have been removed.",
		Color:       shared.Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
