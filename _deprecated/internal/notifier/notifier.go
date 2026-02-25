package notifier

import (
	"github.com/bwmarrin/discordgo"
)

type Notifier interface {
	SendResponse(i *discordgo.InteractionCreate, ir *discordgo.InteractionResponse) error
	SendChannelMessage(channelID string, data *discordgo.MessageSend) error
	CreateDMChannel(userID string) (string, error)
}

func MessageSend(c string, e *discordgo.MessageEmbed, mc ...discordgo.MessageComponent) *discordgo.MessageSend {
	messageSend := discordgo.MessageSend{}
	if c != "" {
		messageSend.Content = c
	}
	if e != nil {
		messageSend.Embeds = []*discordgo.MessageEmbed{e}
	}
	if len(mc) != 0 {
		messageSend.Components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: mc,
			},
		}
	}
	return &messageSend
}
