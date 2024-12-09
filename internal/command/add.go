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

type addCommand struct {
	dCfg *config.DiscordConfig
	sCfg *config.ServiceConfig
	repo repository.Repository
	db   database.Database
}

func NewAddCommand(dCfg *config.DiscordConfig, sCfg *config.ServiceConfig, repo repository.Repository, db database.Database) *addCommand {
	return &addCommand{
		dCfg: dCfg,
		sCfg: sCfg,
		repo: repo,
		db:   db,
	}
}

func (c *addCommand) Command() *discordgo.ApplicationCommand {
	seasonChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(c.sCfg.Seasons))
	for _, season := range c.sCfg.Seasons {
		data := discordgo.ApplicationCommandOptionChoice{
			Name:  season,
			Value: season,
		}
		seasonChoices = append(seasonChoices, &data)
	}

	campusChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(c.sCfg.Campuses))
	for _, campus := range c.sCfg.Campuses {
		data := discordgo.ApplicationCommandOptionChoice{
			Name:  CampusName[campus],
			Value: campus,
		}
		campusChoices = append(campusChoices, &data)
	}

	return &discordgo.ApplicationCommand{
		Name:        "add",
		Description: "Add a snipe request.",
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

func (c *addCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	opts := ParseInteractionOptions(args.Interaction.ApplicationCommandData())
	index := opts["index"]
	var campus, season string
	var ok bool
	if campus, ok = opts["campus"]; !ok {
		campus = c.repo.Campus(args.UserID)
	}
	if season, ok = opts["season"]; !ok {
		season = c.sCfg.DefaultSeason
	}

	if !c.repo.Exists(index, campus, season) {
		rsp := shared.EphemeralEmbedResponse(c.InvalidAdd(index))
		return rsp, nil
	}
	if c.repo.IsSniping(args.UserID, index, campus, season) {
		rsp := shared.EphemeralEmbedResponse(c.DuplicateAdd(index))
		return rsp, nil
	}

	c.repo.AddSnipe(args.UserID, index, campus, season)
	c.repo.Add(index, campus, season)

	course := c.repo.CourseEntry(index, campus, season)
	rsp := shared.EphemeralEmbedResponse(c.SuccessfulAdd(course))
	return rsp, nil
}

func (c *addCommand) InvalidAdd(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: fmt.Sprintf("`%s` does not exist.", index),
		Color:       shared.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *addCommand) DuplicateAdd(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Duplicate Request!",
		Description: fmt.Sprintf("You are already sniping `%s`.", index),
		Color:       shared.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *addCommand) SuccessfulAdd(course shared.CourseEntry) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Successfully Added Request!",
		Description: fmt.Sprintf(
			"`%s` - %s (**Section %s**) was added to your snipe requests.",
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
