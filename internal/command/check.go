package command

import (
	"fmt"
	"strings"
	"time"

	"Tsniper/internal/repository"
	"Tsniper/internal/shared"
	"Tsniper/pkg/config"
	"Tsniper/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type checkCommand struct {
	dCfg *config.DiscordConfig
	repo repository.Repository
	db   database.Database
}

func NewCheckCommand(dCfg *config.DiscordConfig, repo repository.Repository, db database.Database) *checkCommand {
	return &checkCommand{
		dCfg: dCfg,
		repo: repo,
		db:   db,
	}
}

func (c *checkCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "check",
		Description: "View all active snipe requests.",
	}
}

func (c *checkCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	snipes := c.repo.Snipes(args.UserID)
	if len(snipes) == 0 {
		rsp := shared.EphemeralEmbedResponse(c.InvalidCheck())
		return rsp, nil
	}

	var builder strings.Builder
	for _, snipe := range snipes {
		index, campus, season := snipe[0], snipe[1], snipe[2]
		course := c.repo.CourseEntry(index, campus, season)
		count := c.repo.SnipeCount(index, campus, season)
		lastOpen := c.repo.LastOpen(index, campus, season)
		text := fmt.Sprintf("%s `%s` - %s (**Section %s**) | :eyes: `%d` | ", EmojiMap[season], course.Index, course.Title, course.Section, count)
		builder.WriteString(text)
		if lastOpen == -1 {
			builder.WriteString("`Unknown`\n")
		} else {
			builder.WriteString(fmt.Sprintf("<t:%d:R>\n", lastOpen))
		}
	}

	rsp := shared.EphemeralEmbedResponse(c.SuccessfulCheck(builder.String()))
	return rsp, nil
}

func (c *checkCommand) InvalidCheck() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: "You have no active snipe requests.",
		Color:       shared.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (c *checkCommand) SuccessfulCheck(text string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Active Requests",
		Description: text,
		Color:       shared.Blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
