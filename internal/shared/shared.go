package shared

import (
	"github.com/bwmarrin/discordgo"
)

const (
	Blue  = 0x5865f2
	Green = 0x2dcc70
	Red   = 0xe74d3b
)

type CmdArgs struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	UserID      string
}

func ContentResponse(c string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c,
		},
	}
}

func EphemeralContentResponse(c string) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: c,
		},
	}
}

func EmbedResponse(e *discordgo.MessageEmbed) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{e},
		},
	}
}

func EphemeralEmbedResponse(e *discordgo.MessageEmbed) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{e},
		},
	}
}

func ModalResponse(ird *discordgo.InteractionResponseData) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: ird,
	}
}

func AddComponent(rsp *discordgo.InteractionResponseData, c ...discordgo.MessageComponent) *discordgo.InteractionResponseData {
	rsp.Components = []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: c,
		},
	}
	return rsp
}
