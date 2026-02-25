package presentation

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func JoinEmbed(user *discordgo.User, memberCount int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Welcome to the TSniper server!",
		Description: fmt.Sprintf("<@%s> has joined the server!\n\nYou are user **#%d**!", user.ID, memberCount),
		Color:       Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
