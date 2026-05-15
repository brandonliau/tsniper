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

type addCommand struct {
	snipeService  *usecase.SnipeService
	customization *config.CustomizationConfig
}

func AddCommandDefinition(seasons []view.SeasonOption, campuses []view.CampusOption) *discordgo.ApplicationCommand {
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
		Name:        "add",
		Description: "Add a snipe request.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "index",
				Description: "Index of course to add.",
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

func AddCommandHandler(snipeService *usecase.SnipeService, customization *config.CustomizationConfig) interaction.HandleFunc {
	c := &addCommand{
		snipeService:  snipeService,
		customization: customization,
	}
	return c.execute
}

func (c *addCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	options := interaction.ParseInteractionOptions(i)
	req := usecase.AddSnipeRequest{
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
	res, err := c.snipeService.Add(req)
	switch {
	case err == nil:
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.successfulAdd(res.Course)),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrAddSnipeInvalid):
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.invalidAdd(req.Index)),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrAddSnipeDuplicate):
		rsp = interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.duplicateAdd(req.Index)),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}

	return rsp, nil
}

func (c *addCommand) successfulAdd(course *view.CourseView) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Successfully Added Request!",
		Description: fmt.Sprintf(
			"`%s` - %s (**Section %s**) was added to your snipe requests.",
			course.Index,
			course.Title,
			course.Section,
		),
		Color: presentation.Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *addCommand) invalidAdd(index string) *discordgo.MessageEmbed {
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

func (c *addCommand) duplicateAdd(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Duplicate Request!",
		Description: fmt.Sprintf("You are already sniping `%s`.", index),
		Color:       presentation.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
