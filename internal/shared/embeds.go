package shared

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func InvalidCheck() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: "You have no active snipe requests.",
		Color:       Red,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func SuccessfulCheck(text string) *discordgo.MessageEmbed {
	var footer string
	return &discordgo.MessageEmbed{
		Title:       "Active Requests",
		Description: text,
		Color:       Blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
	}
}
