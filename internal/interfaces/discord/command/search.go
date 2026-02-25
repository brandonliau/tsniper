package command

import (
	"fmt"
	"time"

	"tsniper/internal/application/usecase"
	"tsniper/internal/config"
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type searchCommand struct {
	courseService *usecase.CourseService
	customization *config.CustomizationConfig
}

func SearchCommandDefinition(activeScope ...scope.ActiveScope) *discordgo.ApplicationCommand {
	seasonChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(activeScope[0].Seasons()))
	for _, szn := range activeScope[0].Seasons() {
		choice := &discordgo.ApplicationCommandOptionChoice{
			Name:  szn.DisplayName(),
			Value: szn.DisplayName(),
		}
		seasonChoices = append(seasonChoices, choice)
	}

	campusChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(activeScope[0].Campuses()))
	for _, cmp := range activeScope[0].Campuses() {
		choice := &discordgo.ApplicationCommandOptionChoice{
			Name:  cmp.DisplayName(),
			Value: cmp.DisplayName(),
		}
		campusChoices = append(campusChoices, choice)
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
	req := usecase.SearchCourseRequest{
		Index: i.ApplicationCommandData().Options[0].StringValue(),
	}

	if len(i.ApplicationCommandData().Options) > 1 {
		req.Season = utils.Ptr(i.ApplicationCommandData().Options[1].StringValue())
	}

	var rsp *discordgo.InteractionResponse
	res, err := c.courseService.Search(req)
	switch err {
	case nil:
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.successfulSearch(res.Course)),
			interaction.WithEphemeral(),
		)
	case scope.ErrScopeInvalid, course.ErrCourseNotFound:
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.invalidSearch(req.Index)),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}

	return rsp, nil
}

func (c *searchCommand) successfulSearch(course *course.Course) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (`%s`)", course.Title, course.CourseString),
		Color: presentation.Blue,
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
		},
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
