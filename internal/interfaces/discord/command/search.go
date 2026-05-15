package command

import (
	"errors"
	"fmt"
	"time"

	"tsniper/internal/application/usecase"
	"tsniper/internal/application/view"
	"tsniper/internal/config"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type searchCommand struct {
	courseService *usecase.CourseService
	customization *config.CustomizationConfig
}

func SearchCommandDefinition(seasons []view.SeasonOption, campuses []view.CampusOption) *discordgo.ApplicationCommand {
	seasonChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(seasons))
	for _, szn := range seasons {
		seasonChoices = append(seasonChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  szn.Name,
			Value: szn.Value,
		})
	}

	campusChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(campuses))
	for _, cmp := range campuses {
		campusChoices = append(campusChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  cmp.Name,
			Value: cmp.Code,
		})
	}

	return &discordgo.ApplicationCommand{
		Name:        "search",
		Description: "View course information for given index.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "index",
				Description: "Index of course to search.",
				Required:    true,
				MinLength:   utils.Ptr(5),
				MaxLength:   5,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "season",
				Description: "Season of course to search. If not provided, the default season will be used.",
				Required:    false,
				Choices:     seasonChoices,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "campus",
				Description: "Campus of course to search. If not provided, the default campus will be used.",
				Required:    false,
				Choices:     campusChoices,
			},
		},
	}
}

func SearchCommandHandler(courseService *usecase.CourseService, customization *config.CustomizationConfig) interaction.HandleFunc {
	c := &searchCommand{
		courseService: courseService,
		customization: customization,
	}
	return c.execute
}

func (c *searchCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	options := interaction.ParseInteractionOptions(i)
	req := usecase.SearchCourseRequest{
		UserID: interaction.GetUserID(i),
		Index:  options["index"].StringValue(),
	}

	if opt, ok := options["season"]; ok {
		req.Season = utils.Ptr(opt.StringValue())
	}

	if opt, ok := options["campus"]; ok {
		req.Campus = utils.Ptr(opt.StringValue())
	}

	var rsp *discordgo.InteractionResponse
	res, err := c.courseService.Search(req)
	switch {
	case err == nil:
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.successfulSearch(res.Course, res.Count)),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrSearchCourseInvalid):
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.invalidSearch(req.Index)),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}

	return rsp, nil
}

func (c *searchCommand) successfulSearch(crs *view.CourseView, count int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (`%s`)", crs.Title, crs.CourseString),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   ":alarm_clock: Section Meeting Times",
				Value:  fmt.Sprintf(">>> %s", crs.Meeting),
				Inline: false,
			},
			{
				Name:   "Course Name",
				Value:  fmt.Sprintf("`%s`", crs.Title),
				Inline: true,
			},
			{
				Name:   "Index",
				Value:  fmt.Sprintf("`%s`", crs.Index),
				Inline: true,
			},
			{
				Name:   "Section",
				Value:  fmt.Sprintf("`%s`", crs.Section),
				Inline: true,
			},
			{
				Name:   "Instructors",
				Value:  fmt.Sprintf("```fix\n%s```", crs.Instructors),
				Inline: false,
			},
			{
				Name:   "Special Notes",
				Value:  fmt.Sprintf("```fix\n%s```", crs.Notes),
				Inline: false,
			},
			{
				Name:   "Insights",
				Value:  fmt.Sprintf("👀`%d`", count),
				Inline: true,
			},
			{
				Name:   "Last Open",
				Value:  presentation.LastOpenDisplayString(crs.LastOpen),
				Inline: true,
			},
		},
		Color: presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *searchCommand) invalidSearch(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: fmt.Sprintf("`%s` does not exist.", index),
		Color:       presentation.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
