package command

import (
	"fmt"
	"time"

	"tsniper/internal/application/usecase"
	"tsniper/internal/config"
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"
	"tsniper/internal/domain/snipe"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type removeCommand struct {
	snipeService  *usecase.SnipeService
	customization *config.CustomizationConfig
}

func RemoveCommandDefinition(activeScope ...scope.ActiveScope) *discordgo.ApplicationCommand {
	var seasonChoices []*discordgo.ApplicationCommandOptionChoice
	var campusChoices []*discordgo.ApplicationCommandOptionChoice

	if len(activeScope) > 0 {
		seasonChoices = make([]*discordgo.ApplicationCommandOptionChoice, 0, len(activeScope[0].Seasons()))
		for _, szn := range activeScope[0].Seasons() {
			seasonChoices = append(seasonChoices, &discordgo.ApplicationCommandOptionChoice{
				Name:  szn.DisplayName(),
				Value: szn.DisplayName(),
			})
		}

		campusChoices = make([]*discordgo.ApplicationCommandOptionChoice, 0, len(activeScope[0].Campuses()))
		for _, cmp := range activeScope[0].Campuses() {
			campusChoices = append(campusChoices, &discordgo.ApplicationCommandOptionChoice{
				Name:  cmp.DisplayName(),
				Value: cmp.DisplayName(),
			})
		}
	}

	return &discordgo.ApplicationCommand{
		Name:        "remove",
		Description: "Remove a snipe request.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "index",
				Description: "Index of course to remove.",
				Required:    true,
				MinLength:   utils.Ptr(5),
				MaxLength:   5,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "season",
				Description: "Season of course to add. If not provided, the default season will be used.",
				Required:    false,
				Choices:     seasonChoices,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "campus",
				Description: "Campus of course to add. If not provided, the default campus will be used.",
				Required:    false,
				Choices:     campusChoices,
			},
		},
	}
}

func RemoveCommandHandler(snipeService *usecase.SnipeService, customization *config.CustomizationConfig) interaction.HandleFunc {
	c := &removeCommand{
		snipeService:  snipeService,
		customization: customization,
	}
	return c.execute
}

func (c *removeCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	options := interaction.ParseInteractionOptions(i)
	req := usecase.RemoveSnipeRequest{
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
	res, err := c.snipeService.Remove(req)
	switch err {
	case nil:
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.successfulRemove(res.Course)),
			interaction.WithEphemeral(),
		)
	case scope.ErrScopeInvalid, snipe.ErrSnipeNotFound, course.ErrCourseNotFound:
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.invalidRemove(req.Index)),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}

	return rsp, nil
}

func (c *removeCommand) successfulRemove(course *course.Course) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Successfully Removed Request!",
		Description: fmt.Sprintf(
			"`%s` - %s (**Section %s**) was removed from your snipe requests.", course.Index, course.Title, course.Section),
		Color: presentation.Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *removeCommand) invalidRemove(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: fmt.Sprintf("You are not currently sniping `%s`.", index),
		Color:       presentation.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
