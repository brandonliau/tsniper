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

type searchCommand struct {
	dCfg *config.DiscordConfig
	sCfg *config.ServiceConfig
	repo repository.Repository
	db   database.Database
}

func NewSearchCommand(dCfg *config.DiscordConfig, sCfg *config.ServiceConfig, repo repository.Repository, db database.Database) *searchCommand {
	return &searchCommand{
		dCfg: dCfg,
		sCfg: sCfg,
		repo: repo,
		db:   db,
	}
}

func (c *searchCommand) Command() *discordgo.ApplicationCommand {
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
		Name:        "search",
		Description: "View course information for given index.",
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

func (c *searchCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
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
		rsp := shared.EphemeralEmbedResponse(c.InvalidSearch(index))
		return rsp, nil
	}

	course := c.repo.CourseEntry(index, campus, season)
	rsp := shared.EphemeralEmbedResponse(c.SuccessfulSearch(course, campus, season))
	return rsp, nil
}

func (c *searchCommand) InvalidSearch(index string) *discordgo.MessageEmbed {
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

func (c *searchCommand) SuccessfulSearch(course shared.CourseEntry, campus string, season string) *discordgo.MessageEmbed {
	var lastOpen string
	open := c.repo.LastOpen(course.Index, campus, season)
	if open == -1 {
		lastOpen = "`Unknown`\n"
	} else {
		lastOpen = fmt.Sprintf("<t:%d:R>\n", open)
	}
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (`%s`)", course.Title, course.CourseString),
		Color: shared.Blue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   ":alarm_clock: Section Meeting Times",
				Value:  fmt.Sprintf(">>> %s", course.Meeting),
				Inline: false,
			},
			{
				Name:   "Course Name",
				Value:  fmt.Sprintf("`%s`", course.Title),
				Inline: true,
			},
			{
				Name:   "Index",
				Value:  fmt.Sprintf("`%s`", course.Index),
				Inline: true,
			},
			{
				Name:   "Section",
				Value:  fmt.Sprintf("`%s`", course.Section),
				Inline: true,
			},
			{
				Name:   "Instructors",
				Value:  fmt.Sprintf("```fix\n%s```", course.Instructors),
				Inline: false,
			},
			{
				Name:   "Special Notes",
				Value:  fmt.Sprintf("```fix\n%s```", course.Notes),
				Inline: false,
			},
			{
				Name:   "Insights",
				Value:  fmt.Sprintf("👀`%d`", c.repo.SnipeCount(course.Index, campus, season)),
				Inline: true,
			},
			{
				Name:   "Last open",
				Value:  fmt.Sprintf("%s", lastOpen),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
