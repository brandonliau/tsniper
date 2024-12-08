package command

import (
	"fmt"
	"time"

	"Tsniper/internal/repository"
	"Tsniper/internal/shared"
	"Tsniper/pkg/config"
	"Tsniper/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type removeCommand struct {
	dCfg *config.DiscordConfig
	sCfg *config.SnipeConfig
	repo repository.Repository
	db   database.Database
	auth bool
}

func NewRemoveCommand(dCfg config.Config, sCfg config.Config, repo repository.Repository, db database.Database) *removeCommand {
	return &removeCommand{
		dCfg: dCfg.(*config.DiscordConfig),
		sCfg: sCfg.(*config.SnipeConfig),
		repo: repo,
		db:   db,
		auth: false,
	}
}

func (c *removeCommand) Command() *discordgo.ApplicationCommand {
	seasonChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for _, season := range c.sCfg.Seasons {
		data := discordgo.ApplicationCommandOptionChoice{
			Name:  season,
			Value: season,
		}
		seasonChoices = append(seasonChoices, &data)
	}

	campusChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for _, campus := range c.sCfg.Campuses {
		data := discordgo.ApplicationCommandOptionChoice{
			Name:  CampusName[campus],
			Value: campus,
		}
		campusChoices = append(campusChoices, &data)
	}

	return &discordgo.ApplicationCommand{
		Name:        "remove",
		Description: "Remove a snipe request.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "index",
				Description: "index",
				Required:    true,
				MaxLength:   5,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "season",
				Description: "season",
				Required:    false,
				Choices:     seasonChoices,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "campus",
				Description: "campus",
				Required:    false,
				Choices:     campusChoices,
			},
		},
	}
}

func (c *removeCommand) Auth() bool {
	return c.auth
}

func (c *removeCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	opts := ParseInteractionOptions(args.Interaction.ApplicationCommandData())
	index := opts["index"]
	var campus, season string
	var ok bool
	if campus, ok = opts["campus"]; !ok {
		campus = c.repo.Campus(args.UserID)
	}
	if season, ok = opts["season"]; !ok {
		season = c.sCfg.Seasons[0]
	}

	if !c.repo.IsSniping(args.UserID, index, campus, season) {
		rsp := shared.EphemeralEmbedResponse(c.InvalidRemove(index))
		return rsp, nil
	}

	c.repo.RemoveSnipe(args.UserID, index, campus, season)
	c.repo.Remove(index, campus, season)

	course := c.repo.CourseEntry(index, campus, season)
	rsp := shared.EphemeralEmbedResponse(c.SuccessfulRemove(course))
	return rsp, nil
}

func (c *removeCommand) InvalidRemove(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: fmt.Sprintf("You are not currently sniping `%s`.", index),
		Color:       shared.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *removeCommand) SuccessfulRemove(course shared.CourseEntry) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Successfully Removed Request!",
		Description: fmt.Sprintf(
			"`%s` - %s (**Section %s**) was removed from your snipe requests.",
			course.Index,
			course.Title,
			course.Section,
		),
		Color: shared.Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
