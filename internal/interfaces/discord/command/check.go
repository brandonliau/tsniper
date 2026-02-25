package command

import (
	"fmt"
	"strings"
	"time"

	"tsniper/internal/application/usecase"
	"tsniper/internal/config"
	"tsniper/internal/domain/course"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type checkCommand struct {
	snipeService  *usecase.SnipeService
	customization *config.CustomizationConfig
}

func CheckCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "check",
		Description: "View all active snipe requests.",
	}
}

func CheckCommandHandler(snipeService *usecase.SnipeService, customization *config.CustomizationConfig) interaction.HandleFunc {
	c := &checkCommand{
		snipeService:  snipeService,
		customization: customization,
	}
	return c.execute
}

func (c *checkCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	req := usecase.CheckSnipeRequest{
		UserID: interaction.GetUserID(i),
	}

	res, err := c.snipeService.Check(req)
	if err != nil {
		return nil, err
	}

	if len(res.Courses) == 0 {
		rsp := interaction.InteractionInitialResponse(
			interaction.WithEmbeds(c.invalidCheck()),
			interaction.WithEphemeral(),
		)
		return rsp, nil
	}

	rsp := interaction.InteractionInitialResponse(
		interaction.WithEmbeds(c.successfulCheck(res.Courses)),
		interaction.WithEphemeral(),
	)
	return rsp, nil
}

func (c *checkCommand) successfulCheck(courses []*course.Course) *discordgo.MessageEmbed {
	var builder strings.Builder
	for _, crs := range courses {
		fmt.Fprintf(&builder, "`%s` - %s (**Section %s**)\n", crs.Index, crs.Title, crs.Section)
	}
	return &discordgo.MessageEmbed{
		Title:       "Active Requests",
		Description: builder.String(),
		Color:       presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *checkCommand) invalidCheck() *discordgo.MessageEmbed {
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
