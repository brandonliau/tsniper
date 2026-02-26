package command

import (
	"fmt"
	"strings"
	"time"

	"tsniper/internal/application/usecase"
	"tsniper/internal/domain/course"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type checkCommand struct {
	snipeService *usecase.SnipeService
}

func CheckCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "check",
		Description: "View all active snipe requests.",
	}
}

func CheckCommandHandler(snipeService *usecase.SnipeService) interaction.HandleFunc {
	c := &checkCommand{
		snipeService: snipeService,
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
		interaction.WithEmbeds(c.successfulCheck(res.Courses, res.Counts)),
		interaction.WithEphemeral(),
	)
	return rsp, nil
}

func (c *checkCommand) successfulCheck(courses []*course.Course, counts map[*course.Course]int) *discordgo.MessageEmbed {
	var builder strings.Builder
	for _, crs := range courses {
		fmt.Fprintf(
			&builder,
			"%s `%s` - %s (**Section %s**) | :eyes: `%d` | %s",
			presentation.EmojiMap[crs.Scope.Term.DisplayName()],
			crs.Index,
			crs.Title,
			crs.Section,
			counts[crs],
			presentation.LastOpenDisplayString(crs.LastOpen),
		)
	}
	return &discordgo.MessageEmbed{
		Title:       "Active Requests",
		Description: builder.String(),
		Color:       presentation.Blue,
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
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
