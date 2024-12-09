package command

import (
	"fmt"
	"time"

	"Tsniper/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type uptimeCommand struct {
	start int64
}

func NewUptimeCommand(start int64) *uptimeCommand {
	return &uptimeCommand{
		start: start,
	}
}

func (c *uptimeCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "uptime",
		Description: "Check bot uptime.",
	}
}

func (c *uptimeCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	rsp := shared.EphemeralEmbedResponse(c.uptimeEmbed())
	return rsp, nil
}

func (c *uptimeCommand) uptimeEmbed() *discordgo.MessageEmbed {
	diff := time.Now().Unix() - c.start
	return &discordgo.MessageEmbed{
		Title: "Uptime",
		Description: fmt.Sprintf(
			"Last restart: <t:%d:R>\n"+
				"Uptime: %d days, %d hours, %d min, %d sec",
			c.start,
			(diff / 86400),
			((diff / 3600) % 24),
			((diff / 60) % 60),
			(diff % 60),
		),
		Color: shared.Blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
